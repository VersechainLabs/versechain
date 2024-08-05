package MEVless

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func (m *MEVless) SubscribeOrderCommitment(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("SubscribeOrderCommitment: websocket upgrade failed: %s", err)
		return
	}
	defer c.Close()
	//_, msg, err := c.ReadMessage()
	//if err != nil {
	//	logrus.Error("SubscribeOrderCommitment ReadMessage failed: ", err)
	//	return
	//}
	//txnHash := common.BytesToHash(msg)
	//txOrderByt, _, err := m.commitmentsDB.Get(txnHash.Bytes())
	//if err != nil {
	//	logrus.Errorf("SubscribeOrderCommitment Get txOrder(%s) failed: %s", txnHash.String(), err)
	//	return
	//}

	for {
		select {
		case oc := <-m.notifyCh:
			byt, err := json.Marshal(oc)
			if err != nil {
				logrus.Error("SubscribeOrderCommitment json.Marshal failed: ", err)
				continue
			}
			err = c.WriteMessage(websocket.TextMessage, byt)
			if err != nil {
				logrus.Error("SubscribeOrderCommitment WriteMessage failed: ", err)
			}
		}
	}

}

func (m *MEVless) HandleSubscribe() {
	http.HandleFunc("/mev_less", m.SubscribeOrderCommitment)
	http.ListenAndServe(m.cfg.Addr, nil)
}
