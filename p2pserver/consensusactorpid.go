package p2pserver

import "github.com/ontio/ontology-eventbus/actor"

var ConsensusPid *actor.PID

func SetConsensusPid(pid *actor.PID) {
	ConsensusPid = pid
}
