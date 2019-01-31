package p2pserver

type VersionPayload struct {
	Version      uint32
	Services     uint64
	TimeStamp    int64
	SyncPort     uint16
	HttpInfoPort uint16
	ConsPort     uint16
	Cap          [32]byte
	Nonce        uint64
	StartHeight  uint64
	Relay        uint8
	IsConsensus  bool
}

type Version struct {
	P VersionPayload
}

func (this *Version) CmdType() string {
	return VERSION_TYPE
}

func (this *Version) Deserialization(data []byte) error {
	return nil
}