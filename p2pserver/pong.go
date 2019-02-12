package p2pserver

type Pong struct {
	Height    uint64
}

func (this *Pong) CmdType() string {
	return PONG_TYPE
}

func (this *Pong) Deserialization(source *ZeroCopySource) error {
	return nil
}

func (this *Pong) Serialization(sink *ZeroCopySink) (err error) {
	return nil
}
