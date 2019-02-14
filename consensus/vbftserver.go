package consensus

import (
	"github.com/biety/p2pserver"
	"github.com/ontio/ontology-eventbus/actor"
)

type Server struct {
	pid               *actor.PID
	txpoolactorpid    *actor.PID
	p2pactorpid       *actor.PID

	msgRecvC          map[uint32]chan *p2pMsgPayload
}

func NewVbftServer() (*Server, error) {
	server := &Server {

	}

	props := actor.FromProducer(func() actor.Actor {return server})
	pid, err := actor.SpawnNamed(props, "consensus_vbft")
	if err != nil {
		return nil, err
	}
	server.pid = pid
	return server, nil
}

func (self *Server) Receive(context actor.Context) {
	switch msg:= context.Message().(type) {
	case *StartConsensus:
	case *StopConsensus:
		self.stop()
	case *p2pserver.ConsensusPayload:
		self.NewConsensusPayload(msg)
	}
}


func (self *Server) GetPID() *actor.PID {
	return self.pid
}

func (self *Server) Start() error {
	return self.start()
}

func (self *Server) Init(p2pactorpid *actor.PID, txpoolactorpid *actor.PID) error {
	return self.init(p2pactorpid, txpoolactorpid)
}

func (self *Server) init(p2pactorpid *actor.PID, txpoolactorpid *actor.PID) error {
	self.txpoolactorpid = txpoolactorpid
	self.p2pactorpid = p2pactorpid

	self.msgRecvC = make(map[uint32]chan *p2pMsgPayload)
	self.msgRecvC[0] = make(chan *p2pMsgPayload, 1024)

	return nil
}


func (self *Server) start() error {
	return nil
}

func (self *Server) stop() error {
	return nil
}

func (self *Server) NewConsensusPayload(payload *p2pserver.ConsensusPayload) {
	pl := &p2pMsgPayload{
		payload: payload,
	}

	self.msgRecvC[0] <- pl
}



