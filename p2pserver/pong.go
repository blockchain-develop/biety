package p2pserver

type Pong struct {
	Height    uint64
}

func (this *Pong) CmdType() string {
	return PONG_TYPE
}

func (this *Pong) Deserialization(data []byte) error {
	return nil
}
