package p2pserver

import "errors"

type Peer struct {
	SyncLink        *Link
	ConsLink        *Link
}


func NewPeer()  *Peer {
	p := &Peer {

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
