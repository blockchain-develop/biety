package p2pserver

type Ping struct {
	Height    uint64
}

func (this *Ping) CmdType() string {
	return PING_TYPE
}

func (this *Ping) Deserialization(source *ZeroCopySource) error {
	return nil
}

func (this *Ping) Serialization(sink *ZeroCopySink) (err error) {
	return nil
}
