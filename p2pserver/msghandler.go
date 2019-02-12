package p2pserver

import (
	"fmt"
	"time"
)

func VersionHandle(data* MsgPayload, p2p *P2PServer, args ...interface{}) {
	fmt.Printf("receive version message")
	version := data.Payload.(*Version)
	remotepeer := p2p.GetPeerFromAddr(data.Addr)
	if remotepeer == nil {
		return
	}

	if version.P.IsConsensus == true {
		remotepeer.UpdateInfo(time.Now(), version.P.Version, version.P.Services, version.P.SyncPort,
			version.P.ConsPort, version.P.Nonce, version.P.Relay, version.P.StartHeight)

		s := remotepeer.GetConsState()
		var msg Message
		if s == INIT {
			remotepeer.SetConsState(HAND_SHAKE)
			msg = NewVersion()
		} else if (s == HAND) {
			remotepeer.SetConsState(HAND_SHAKED)
			msg = NewVerAck(true)
		}

		err := p2p.Send(remotepeer, msg, true)
		if err != nil {
			return
		}
	} else {
		remotepeer.UpdateInfo(time.Now(), version.P.Version, version.P.Services, version.P.SyncPort,
			version.P.ConsPort, version.P.Nonce, version.P.Relay, version.P.StartHeight)

		s := remotepeer.GetSyncState()
		var msg Message
		if s == INIT {
			remotepeer.SetSyncState(HAND_SHAKE)
			msg = NewVersion()
		} else if (s == HAND) {
			remotepeer.SetSyncState(HAND_SHAKED)
			msg = NewVerAck(false)
		}

		err := p2p.Send(remotepeer, msg, false)
		if err != nil {
			return
		}
	}
}

func VersionAck(data* MsgPayload, p2p *P2PServer, args ...interface{}) {
	fmt.Printf("receive verack message")
	verack := data.Payload.(*VerACK)
	remotepeer := p2p.GetPeerFromAddr(data.Addr)
	if remotepeer == nil {
		return
	}

	if verack.IsConsensus == true {
		s := remotepeer.GetConsState()
		if s != HAND_SHAKE && s != HAND_SHAKED {
			return
		}

		remotepeer.SetConsState(ESTABLISH)
		if s == HAND_SHAKE {
			msg := NewVerAck(true)
			p2p.Send(remotepeer, msg, true)
		}
	} else {
		s := remotepeer.GetSyncState()
		if s != HAND_SHAKE && s != HAND_SHAKED {
			return
		}

		remotepeer.SetSyncState(ESTABLISH)
		if s == HAND_SHAKE {
			msg := NewVerAck(false)
			p2p.Send(remotepeer, msg, false)
		}
	}
}

func PingHandle(data* MsgPayload, p2p *P2PServer, args ...interface{}) {

}

func PongHandle(data* MsgPayload, p2p *P2PServer, args ...interface{}) {

}
