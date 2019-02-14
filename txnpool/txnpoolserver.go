package txnpool

import (
	"github.com/biety/block"
	"github.com/biety/validator"
	"github.com/ontio/ontology-eventbus/actor"
)

type TxPoolServer struct {
	actors    map[string]*actor.PID
	workers    []TxPoolWorker
	txPool     *TxPool
	entries     map[string][]*validator.RegisterValidator
}

func NewTxPoolServer() *TxPoolServer {
	s := &TxPoolServer {}
	s.init()
	return s
}

func (s *TxPoolServer) init() {
	s.txPool = &TxPool{}
	s.txPool.Init()

	s.actors = make(map[string]*actor.PID)

	s.entries = make(map[string][]*validator.RegisterValidator)

	s.workers = make([]TxPoolWorker, 8)
	var i uint8
	for i = 0;i < 8;i ++ {
		s.workers[i].init(s)
		go s.workers[i].start()
	}
}

func (s *TxPoolServer) registerValidator(v *validator.RegisterValidator) {
	_,ok := s.entries[v.Id]
	if !ok {
		s.entries[v.Id] = make([]*validator.RegisterValidator, 0, 1)
	}
	s.entries[v.Id] = append(s.entries[v.Id], v)
}

func (s *TxPoolServer) unRegisterValidator(id string) {

}

func (s *TxPoolServer) RegisterActor(id string, pid *actor.PID) {
	s.actors[id] = pid
}

func (s *TxPoolServer) GetPID(id string) *actor.PID {
	return s.actors[id]
}

func (s *TxPoolServer) GetVerifyRspActorPID() *actor.PID {
	return s.actors["txVerifyRsp"]
}

func (s *TxPoolServer) GetTxPoolActorPID() *actor.PID {
	return s.actors["txPool"]
}

func (s *TxPoolServer) GetTxActorPID() *actor.PID {
	return s.actors["TxActor"]
}

func (s *TxPoolServer) assignTxToWorker(tx *block.Transaction) bool {
	if tx == nil {
		return false
	}

	s.workers[0].rcvTxCh <- tx
	return true
}

func (s *TxPoolServer) assignRspToWorker(rsp *validator.CheckResponse) bool {
	if rsp == nil {
		return false
	}

	s.workers[0].rspCh <- rsp
	return true
}

func (s *TxPoolServer) getNextValidatorPIDs() []*actor.PID {
	if len(s.entries) == 0 {
		return nil
	}

	ret := make([]*actor.PID, 0, len(s.entries))
	for _,v := range s.entries {
		ret = append(ret, v[0].Sender)
	}

	return ret
}

func (s *TxPoolServer) getNextValidatorPID(id string) *actor.PID {
	length := len(s.entries[id])
	if length == 0 {
		return nil
	}

	return s.entries[id][0].Sender
}


func (s *TxPoolServer) addTxList(txentry *TxEntry) bool {
	ret := s.txPool.AddTxList(txentry)
	return ret
}