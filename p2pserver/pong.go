package p2pserver

import "github.com/biety/common"

type Pong struct {
	Height    uint64
}

func (this *Pong) CmdType() string {
	return PONG_TYPE
}

func (this *Pong) Deserialization(source *common.ZeroCopySource) error {
	return nil
}

func (this *Pong) Serialization(sink *common.ZeroCopySink) (err error) {
	return nil
}
