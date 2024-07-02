package pkg

import (
	"math/rand"
	"time"
)

const (
	maxTransfer = 100
)

type CaseEthWallet struct {
	*EthWallet
	EthCount int64 `json:"ethCount"`
}

func (c *CaseEthWallet) Copy() *CaseEthWallet {
	return &CaseEthWallet{
		EthWallet: c.EthWallet.Copy(),
		EthCount:  c.EthCount,
	}
}

type TransferManager struct {
}

func NewTransferManager() *TransferManager {
	return &TransferManager{}
}

func GenerateCaseWallets(initialEthCount int64, wallets []*EthWallet) []*CaseEthWallet {
	c := make([]*CaseEthWallet, 0)
	for _, w := range wallets {
		c = append(c, &CaseEthWallet{
			EthWallet: w,
			EthCount:  initialEthCount,
		})
	}
	return c
}

func (m *TransferManager) GenerateTransferSteps(stepCount int, wallets []*CaseEthWallet) *TransferCase {
	t := &TransferCase{
		Original: getCopy(wallets),
		Expect:   getCopy(wallets),
	}
	steps := make([]*Step, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < stepCount; i++ {
		steps = append(steps, generateStep(r, wallets, maxTransfer))
	}
	t.Steps = steps
	calculateExpect(t)
	return t
}

func calculateExpect(tc *TransferCase) {
	for _, step := range tc.Steps {
		calculate(step, tc.Expect)
	}
}

func calculate(step *Step, expect map[string]*CaseEthWallet) {
	fromWallet := expect[step.From.Address]
	toWallet := expect[step.To.Address]
	fromWallet.EthCount = fromWallet.EthCount - step.Count
	toWallet.EthCount = toWallet.EthCount + step.Count
	expect[step.From.Address] = fromWallet
	expect[step.To.Address] = toWallet
}

func generateStep(r *rand.Rand, wallets []*CaseEthWallet, maxTransfer int) *Step {
	from := r.Intn(len(wallets))
	to := r.Intn(len(wallets))
	for from == to {
		to = r.Intn(len(wallets))
	}
	transferCount := r.Intn(maxTransfer) + 1
	return &Step{
		From:  wallets[from].EthWallet,
		To:    wallets[to].EthWallet,
		Count: int64(transferCount),
	}
}

func getCopy(wallets []*CaseEthWallet) map[string]*CaseEthWallet {
	m := make(map[string]*CaseEthWallet)
	for _, w := range wallets {
		m[w.Address] = w.Copy()
	}
	return m
}

type TransferCase struct {
	Steps []*Step `json:"steps"`
	// address to wallet
	Original map[string]*CaseEthWallet `json:"original"`
	Expect   map[string]*CaseEthWallet `json:"expect"`
}

type Step struct {
	From  *EthWallet `json:"from"`
	To    *EthWallet `json:"to"`
	Count int64      `json:"count"`
}
