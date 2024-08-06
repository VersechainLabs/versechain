package MEVless

import (
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
	"slices"
	"sort"
	"strings"
	"time"
)

const notifyBufferLen = 10

type MEVless struct {
	*tripod.Tripod
	cfg *Config

	commitmentsDB *pebble.DB

	notifyCh chan *OrderCommitment
}

const Prefix = "MEVless_"

func NewMEVless(cfg *Config) (*MEVless, error) {
	db, err := pebble.Open(cfg.DbPath, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	tri := &MEVless{
		Tripod:        tripod.NewTripod(),
		cfg:           cfg,
		commitmentsDB: db,
		notifyCh:      make(chan *OrderCommitment, notifyBufferLen),
	}

	tri.SetWritings(tri.OrderTx)

	go tri.HandleSubscribe()
	return tri, nil
}

func (m *MEVless) OrderTx(ctx *context.WriteContext) error {
	fmt.Printf("[OrderTx] %s\n]", ctx.Txn.TxnHash.Hex())
	return nil
}

func (m *MEVless) Pack(blockNum common.BlockNum, numLimit uint64) ([]*types.SignedTxn, error) {
	err := m.OrderCommitment(blockNum)
	if err != nil {
		return nil, err
	}
	return m.Pool.PackFor(numLimit, func(txn *types.SignedTxn) bool {
		return txn.ParamsIsJson()
	})
}

func (m *MEVless) PackFor(blockNum common.BlockNum, numLimit uint64, filter func(*types.SignedTxn) bool) ([]*types.SignedTxn, error) {
	err := m.OrderCommitment(blockNum)
	if err != nil {
		return nil, err
	}
	return m.Pool.PackFor(numLimit, func(txn *types.SignedTxn) bool {
		if txn.ParamsIsJson() {
			return filter(txn)
		}
		return false
	})
}

// wrCall.params = "MEVless_(TxnHash)"
func (m *MEVless) OrderCommitment(blockNum common.BlockNum) error {
	hashTxns, err := m.Pool.PackFor(m.cfg.PackNumber, func(txn *types.SignedTxn) bool {
		paramStr := txn.GetParams()
		if !strings.HasPrefix(paramStr, Prefix) {
			return false
		}
		hashStr := strings.TrimPrefix(paramStr, Prefix)
		hashByt := common.HexToHash(hashStr)
		return len(hashByt) == common.HashLen
	})
	if err != nil {
		return err
	}
	if len(hashTxns) == 0 {
		return nil
	}

	sequence := m.makeOrder(hashTxns)

	orderCommitment := &OrderCommitment{
		BlockNumber: blockNum,
		Sequences:   sequence,
	}

	m.notifyClient(orderCommitment)

	err = m.storeOrderCommitment(orderCommitment)
	if err != nil {
		return err
	}

	// TODO: sync the OrderCommitment to other P2P nodes

	// sleep for a while so that clients can send their tx-content onchain.
	time.Sleep(800 * time.Millisecond)

	m.Pool.Reset(hashTxns)

	m.Pool.SortTxns(func(txs []*types.SignedTxn) []*types.SignedTxn {
		sorted := make([]*types.SignedTxn, len(sequence))
		for num, hash := range sequence {
			txn, err := m.Pool.GetTxn(hash)
			if err != nil {
				logrus.Error("MEVless get txn from txpool failed: ", err)
				continue
			}
			if txn != nil {
				sorted[num] = txn
				txs = slices.DeleteFunc(txs, func(txn *types.SignedTxn) bool {
					return txn.TxnHash == hash
				})
			}
		}

		sorted = append(sorted, txs...)

		return sorted
	})

	return nil
}

func (m *MEVless) makeOrder(hashTxns []*types.SignedTxn) map[int]common.Hash {
	order := make(map[int]common.Hash)
	sort.Slice(hashTxns, func(i, j int) bool {
		return hashTxns[i].GetTips() > hashTxns[j].GetTips()
	})
	for i, txn := range hashTxns {
		spew.Dump(txn)
		hashStr := strings.TrimPrefix(txn.GetParams(), Prefix)
		order[i] = common.BytesToHash([]byte(hashStr))
	}
	return order
}

func (m *MEVless) VerifyBlock(block *types.Block) error {
	// TODO: verify tx order commitment from other miner node
	// TODO: for double-check, fetch txs from DA layers if the block does not have commitment txs.
	return nil
}

func (m *MEVless) Charge() uint64 {
	return m.cfg.Charge
}

type OrderCommitment struct {
	BlockNumber common.BlockNum     `json:"block_number"`
	Sequences   map[int]common.Hash `json:"sequences"`
}

type TxOrder struct {
	BlockNumber common.BlockNum `json:"block_number"`
	Sequence    int             `json:"sequence"`
}

func (m *MEVless) storeOrderCommitment(oc *OrderCommitment) error {
	batch := m.commitmentsDB.NewBatch()
	for seq, txnHash := range oc.Sequences {
		txOrder := &TxOrder{
			BlockNumber: oc.BlockNumber,
			Sequence:    seq,
		}
		byt, err := json.Marshal(txOrder)
		if err != nil {
			return err
		}
		err = batch.Set(txnHash.Bytes(), byt, pebble.NoSync)
		if err != nil {
			return err
		}
	}
	return batch.Commit(pebble.Sync)
}

func (m *MEVless) notifyClient(oc *OrderCommitment) {
	fmt.Printf("[NotifyClient] %#v\n", oc)
	if len(m.notifyCh) == notifyBufferLen {
		_ = <-m.notifyCh
	}
	m.notifyCh <- oc
}
