package p2pserver

import (
	"github.com/biety/block"
	"github.com/ontio/ontology-eventbus/actor"
)

var txnPoolPid *actor.PID

func SetTxnPoolPid(txnPid *actor.PID) {
	txnPoolPid = txnPid
}

func AddTransaction(transaction *block.Transaction) {
	if txnPoolPid == nil {
		return
	}

	txReq := &TxReq {
		Tx : transaction,
	}

	txnPoolPid.Tell(txReq)
}
