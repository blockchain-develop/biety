package txnpool

import (
	"github.com/biety/base"
	"github.com/biety/block"
	"github.com/biety/common"
	"github.com/biety/validator"
	"github.com/ontio/ontology/errors"
	"time"
)


type pendingTx struct {
	tx    *block.Transaction
	req   *validator.CheckTx
	ret   []*TxAttr
}

type TxPoolWorker struct {
	timer          *time.Timer
	rcvTxCh        chan *block.Transaction
	rspCh          chan *validator.CheckResponse

	pendingTxList  map[common.Uint256]*pendingTx

	server         *TxPoolServer
}

func (worker  *TxPoolWorker) init(s *TxPoolServer) {
	worker.rcvTxCh = make(chan *block.Transaction, base.MAX_PENDING_TXN)
	worker.rspCh = make(chan *validator.CheckResponse, base.MAX_PENDING_TXN)
	worker.pendingTxList = make(map[common.Uint256]*pendingTx)
	worker.server = s
}

func (worker *TxPoolWorker) start() {
	worker.timer = time.NewTimer(time.Second * base.EXPIRE_INTERVAL)
	for {
		select {
		case rcvTx, ok := <- worker.rcvTxCh:
			if ok {
				worker.verifyTx(rcvTx)
			}
		case rsp, ok := <- worker.rspCh:
			if ok {
				worker.handleRsp(rsp)
			}
		}
	}
}

func (worker *TxPoolWorker) verifyTx(tx *block.Transaction) {
	req := &validator.CheckTx {
		Tx: tx,
	}

	worker.sendReq2Validator(req)

	pt := &pendingTx {
		tx: tx,
		req: req,
	}
	worker.pendingTxList[tx.Hash()] = pt
}

func (worker *TxPoolWorker) sendReq2Validator(req *validator.CheckTx) bool {
	rspactorpid := worker.server.GetVerifyRspActorPID()
	if rspactorpid == nil {
		return false
	}

	pids := worker.server.getNextValidatorPIDs()
	if pids == nil {
		return false
	}

	for _, pid := range pids {
		pid.Request(req, rspactorpid)
	}

	return true
}

func (worker *TxPoolWorker) handleRsp(rsp *validator.CheckResponse) {
	if rsp.ErrCode != errors.ErrNoError {
		delete(worker.pendingTxList, rsp.Hash)
		return
	}

	pt,ok := worker.pendingTxList[rsp.Hash]
	if !ok {
		return
	}

	retAttr := &TxAttr{
		Height:rsp.Height,
		ErrCode:rsp.ErrCode,
	}
	pt.ret = append(pt.ret, retAttr)
	worker.putTxPool(pt)
	delete(worker.pendingTxList, rsp.Hash)
}

func (worker *TxPoolWorker) putTxPool(pt *pendingTx) bool {
	txentry := &TxEntry{
		Tx : pt.tx,
		Attrs: pt.ret,
	}

	worker.server.addTxList(txentry)
	return true
}