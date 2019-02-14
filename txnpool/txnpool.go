package txnpool

import (
	"fmt"
	"github.com/ontio/ontology-eventbus/actor"
)

func startActor(obj interface{}, id string) (*actor.PID, error) {
	props := actor.FromProducer(func() actor.Actor {return obj.(actor.Actor)})
	pid, _ := actor.SpawnNamed(props, id)
	if pid == nil {
		return nil, fmt.Errorf("fail to start actor")
	}
	return pid, nil
}


func StartTxnPoolServer() (*TxPoolServer, error) {
	s := NewTxPoolServer()

	//
	rspActor := NewVerifyRspActor(s)
	rspactorpid, err := startActor(rspActor, "txVerifyRsp")
	if rspactorpid == nil {
		return nil, err
	}
	s.RegisterActor("txVerifyRsp", rspactorpid)

	//
	txpoolactor := NewTxPoolActor(s)
	txpoolactorpid, err := startActor(txpoolactor, "txPool")
	if txpoolactorpid == nil {
		return nil, err
	}
	s.RegisterActor("txPool", rspactorpid)

	//
	txactor := NewTxActor(s)
	txactorpid, err := startActor(txactor, "tx")
	if txactorpid == nil {
		return nil, err
	}
	s.RegisterActor("tx", rspactorpid)

	return s, nil
}

func (tp *TxPool) AddTxList(txentry *TxEntry) bool {
	txHash := txentry.Tx.Hash()
	if _,ok := tp.txList[txHash]; ok {
		return false
	}

	tp.txList[txHash] = txentry
	return true
}
