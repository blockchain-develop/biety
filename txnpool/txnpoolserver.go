package txnpool

import "github.com/ontio/ontology-eventbus/actor"

type TxPoolServer struct {
	actors    map[string]*actor.PID
}

func NewTxPoolServer() *TxPoolServer {
	s := &TxPoolServer {}
	return s
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
