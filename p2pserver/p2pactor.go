package p2pserver

import (
	"github.com/biety/block"
	"github.com/ontio/ontology-eventbus/actor"
)

type P2PActor struct {
	props           *actor.Props
}

func NewP2PActor() *P2PActor {
	return &P2PActor{
	}
}

func (this *P2PActor) Start() (*actor.PID, error) {
	this.props = actor.FromProducer(func() actor.Actor {return this})
	p2pPID, err := actor.SpawnNamed(this.props, "net_server")
	return p2pPID, err
}

func (this *P2PActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AppendHeaders:
		this.OnHeaderReceive(msg.FromID, msg.Headers)
	case *AppendBlock:
		this.OnBlockReceive(msg.FromID, msg.BlockSize, msg.Block)
	}
}

func (this *P2PActor) OnHeaderReceive(fromID uint64, headers []*block.Header) {

}

func (this *P2PActor) OnBlockReceive(fromID uint64, blockSize uint32, block *block.Block) {

}
