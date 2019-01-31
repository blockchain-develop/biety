package p2pserver

type Ping struct {
	Height    uint64
}

func (this *Ping) CmdType() string {
	return PING_TYPE
}

func (this *Ping) Deserialization(data []byte) error {
	return nil
}
