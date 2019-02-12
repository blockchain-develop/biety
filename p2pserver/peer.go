package p2pserver

import (
	"errors"
	"time"
)

type PeerCom struct {
	id             uint64
	version        uint32
	services       uint64
	relay          bool
	httpInfoPort   uint16
	syncPort       uint16
	consPort       uint16
	height         uint64
}

// SetID sets a peer's id
func (this *PeerCom) SetID(id uint64) {
	this.id = id
}

// GetID returns a peer's id
func (this *PeerCom) GetID() uint64 {
	return this.id
}

// SetVersion sets a peer's version
func (this *PeerCom) SetVersion(version uint32) {
	this.version = version
}

// GetVersion returns a peer's version
func (this *PeerCom) GetVersion() uint32 {
	return this.version
}

// SetServices sets a peer's services
func (this *PeerCom) SetServices(services uint64) {
	this.services = services
}

// GetServices returns a peer's services
func (this *PeerCom) GetServices() uint64 {
	return this.services
}

// SerRelay sets a peer's relay
func (this *PeerCom) SetRelay(relay bool) {
	this.relay = relay
}

// GetRelay returns a peer's relay
func (this *PeerCom) GetRelay() bool {
	return this.relay
}

// SetSyncPort sets a peer's sync port
func (this *PeerCom) SetSyncPort(port uint16) {
	this.syncPort = port
}

// GetSyncPort returns a peer's sync port
func (this *PeerCom) GetSyncPort() uint16 {
	return this.syncPort
}

// SetConsPort sets a peer's consensus port
func (this *PeerCom) SetConsPort(port uint16) {
	this.consPort = port
}

// GetConsPort returns a peer's consensus port
func (this *PeerCom) GetConsPort() uint16 {
	return this.consPort
}

// SetHeight sets a peer's height
func (this *PeerCom) SetHeight(height uint64) {
	this.height = height
}

// GetHeight returns a peer's height
func (this *PeerCom) GetHeight() uint64 {
	return this.height
}

type Peer struct {
	base            PeerCom

	SyncLink        *Link
	ConsLink        *Link

	syncState        uint32
	consState        uint32
}


func NewPeer()  *Peer {
	p := &Peer {
		syncState: INIT,
		consState: INIT,
	}

	p.SyncLink = NewLink()
	p.ConsLink = NewLink()

	return p
}

func (this *Peer) Send(msg Message, isConsensus bool) error {
	if isConsensus {
		return this.SendToCons(msg)
	} else {
		return this.SendToSync(msg)
	}
}

func (this *Peer) SendToSync(msg Message) error {
	if this.SyncLink != nil {
		return this.SyncLink.Tx(msg)
	}

	return errors.New("sync link invalid")
}

func (this *Peer) SendToCons(msg Message) error {
	if this.ConsLink != nil {
		return this.ConsLink.Tx(msg)
	}
	return errors.New("cons link invalid")
}

func (this *Peer) SetSyncState(state uint32) {
	this.syncState = state
}

func (this *Peer) GetSyncState() uint32 {
	return this.syncState
}

func (this *Peer) SetConsState(state uint32) {
	this.consState = state
}

func (this *Peer) GetConsState() uint32 {
	return this.consState
}

func (this *Peer) UpdateInfo(t time.Time, version uint32, services uint64,
	syncPort uint16, consPort uint16, nonce uint64, relay uint8, height uint64) {
	this.base.SetID(nonce)
	this.base.SetVersion(version)
	this.base.SetServices(services)
	this.base.SetSyncPort(syncPort)
	this.base.SetConsPort(consPort)
	this.SyncLink.SetPort(syncPort)
	this.ConsLink.SetPort(consPort)
	if relay == 0 {
		this.base.SetRelay(false)
	} else {
		this.base.SetRelay(true)
	}
	this.base.SetHeight(uint64(height))
}
