package p2pserver

import "github.com/biety/common"

type VerACK struct {
	IsConsensus  bool
}

func (this *VerACK) CmdType() string {
	return VERACK_TYPE
}

//Deserialize message payload
func (this *VerACK) Deserialization(source *common.ZeroCopySource) error {
	return nil
}

func (this *VerACK) Serialization(sink *common.ZeroCopySink) (err error) {
	return nil
}

func NewVerAck(isConsensus bool) Message {
	var verack VerACK
	verack.IsConsensus = isConsensus

	return &verack
}
