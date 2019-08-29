package consensus

import (
	"time"
)

var (
	makeProposalTimeout    = 300 * time.Millisecond
	make2ndProposalTimeout = 300 * time.Millisecond
	endorseBlockTimeout    = 100 * time.Millisecond
	commitBlockTimeout     = 200 * time.Millisecond
	peerHandshakeTimeout   = 10 * time.Second
	txPooltimeout          = 1 * time.Second
	zeroTxBlockTimeout     = 10 * time.Second
)

type ServerState int

const (
	Init ServerState = iota
	LocalConfigured
	Configured       // config loaded from chain
	Syncing          // syncing block from neighbours
	WaitNetworkReady // sync reached, and keep synced, try connecting with more peers
	SyncReady        // start processing consensus msg, but not broadcasting proposal/endorse/commit
	Synced           // start bft
	SyncingCheck     // potentially lost syncing
)

type StateEventType int

const (
	ConfigLoaded     StateEventType = iota
	UpdatePeerConfig                // notify statemgmt on peer heartbeat
	UpdatePeerState                 // notify statemgmt on peer heartbeat
	SyncReadyTimeout
	SyncDone
	LiveTick
)

type PeerState struct {
	peerIdx           uint32
	chainConfigView   uint32
	committedBlockNum uint32
	connected         bool
}

type StateEvent struct {
	Type         StateEventType
	peerState    *PeerState
	blockNum     uint32
}
