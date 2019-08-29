package consensus

import "time"

type StateMgr struct {
	server          *Server
	liveTicker      *time.Timer
	StateEventC     chan *StateEvent

	currentState    ServerState
}

func NewStateMgr(server *Server) *StateMgr {
	sm := &StateMgr{
		server: server,
	}
	return sm
}

func (self *StateMgr) run() {
	self.liveTicker = time.AfterFunc(peerHandshakeTimeout * 5, func() {
		self.StateEventC <- &StateEvent {
			Type: LiveTick,
			blockNum : 0,
		}
		self.liveTicker.Reset(peerHandshakeTimeout * 3)
	})

	for {
		select {
		case evt := <- self.StateEventC:
			switch evt.Type {
			case LiveTick:
				self.onLiveTick(evt)
			case SyncDone:
				self.setSyncedReady()
			}
		}
	}
}

func (self *StateMgr) onLiveTick(evt *StateEvent) error {
	if self.getState() != Synced {
		return nil
	}

	self.server.makeFastForward()

	self.server.reBroadcastCurrentRoundMsgs()

	return nil
}

func (self *StateMgr) getState() ServerState {
	return self.currentState
}
