package consensus

import "github.com/biety/p2pserver"

type StartConsensus struct {}
type StopConsensus struct {}

type p2pMsgPayload struct {
	fromPeer      uint32
	payload       *p2pserver.ConsensusPayload
}
