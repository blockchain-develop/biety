package p2pserver

import "github.com/biety/common"

type Ping struct {
	Height    uint64
}

func (this *Ping) CmdType() string {
	return PING_TYPE
}

func (this *Ping) Deserialization(source *common.ZeroCopySource) error {
	return nil
}

func (this *Ping) Serialization(sink *common.ZeroCopySink) (err error) {
	return nil
}
