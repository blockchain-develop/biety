package p2pserver

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
