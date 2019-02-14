package txnpool

import (
	"github.com/biety/base"
	"github.com/biety/block"
	"github.com/biety/p2pserver"
	"github.com/biety/validator"
	"github.com/ontio/ontology-eventbus/actor"
)

type TxActor struct {
	server    *TxPoolServer
}

func (ta *TxActor) setServer(s *TxPoolServer) {
	ta.server = s
}

func (ta *TxActor) handleTransaction(self *actor.PID, txn *block.Transaction) {
	if len(txn.ToArray()) > base.MAX_TX_SIZE {
		return
	}

	// todo
	// if transaction has exist

	// todo
	// if transaction pool is full

	ta.server.assignTxToWorker(txn)
}

func (ta *TxActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *p2pserver.TxReq:
		ta.handleTransaction(context.Self(), msg.Tx)
	}
}

func NewTxActor(s *TxPoolServer) *TxActor {
	a := &TxActor{}
	a.setServer(s)
	return a
}


type TxPoolActor struct {
	server    *TxPoolServer
}

func (tpa *TxPoolActor) setServer(s *TxPoolServer) {
	tpa.server = s
}

func NewTxPoolActor(s *TxPoolServer) *TxPoolActor {
	a := &TxPoolActor{}
	a.setServer(s)
	return a
}

func (tpa *TxPoolActor) Receive(context actor.Context) {

}

type VerifyRspActor struct {
	server    *TxPoolServer
}

func (vpa *VerifyRspActor) setServer(s *TxPoolServer) {
	vpa.server = s
}

func NewVerifyRspActor(s *TxPoolServer) *VerifyRspActor {
	a := &VerifyRspActor {}
	a.setServer(s)
	return a
}

func (vpa *VerifyRspActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *validator.RegisterValidator:
		vpa.server.registerValidator(msg)
	case *validator.UnRegisterValidator:
		vpa.server.unRegisterValidator(msg.Id)
	case *validator.CheckResponse:
		vpa.server.assignRspToWorker(msg)
	}
}
