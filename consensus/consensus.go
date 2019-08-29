package consensus

import "github.com/ontio/ontology-eventbus/actor"

type ConsensusService interface {
	Start() error
	GetPID()   *actor.PID
	Init(p2pactorpid *actor.PID, txpoolactorpid *actor.PID) error
}

func NewConsensueService() (ConsensusService, error) {
	server, err := NewVbftServer()
	return server, err
}
