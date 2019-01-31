package p2pserver

type VerACK struct {
	IsConsensus  bool
}

func (this *VerACK) CmdType() string {
	return VERACK_TYPE
}

//Deserialize message payload
func (this *VerACK) Deserialization(data []byte) error {
	return nil
}
