package p2pserver

import (
	"errors"
	"github.com/biety/common"
	"io"
)

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

func (this *Version) Deserialization(source *common.ZeroCopySource) error {
	var irregular, eof bool
	this.P.Version, eof = source.NextUint32()
	this.P.Services, eof = source.NextUint64()
	this.P.TimeStamp, eof = source.NextInt64()
	this.P.SyncPort, eof = source.NextUint16()
	this.P.HttpInfoPort, eof = source.NextUint16()
	this.P.ConsPort, eof = source.NextUint16()
	var buf []byte
	buf, eof = source.NextBytes(uint64(len(this.P.Cap[:])))
	copy(this.P.Cap[:], buf)

	this.P.Nonce, eof = source.NextUint64()
	this.P.StartHeight, eof = source.NextUint64()
	this.P.Relay, eof = source.NextUint8()
	this.P.IsConsensus, irregular, eof = source.NextBool()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return errors.New("irregular data")
	}

	return nil
}

func (this *Version) Serialization(sink *common.ZeroCopySink) (err error) {
	sink.WriteUint32(this.P.Version)
	sink.WriteUint64(this.P.Services)
	sink.WriteInt64(this.P.TimeStamp)
	sink.WriteUint16(this.P.SyncPort)
	sink.WriteUint16(this.P.HttpInfoPort)
	sink.WriteUint16(this.P.ConsPort)
	sink.WriteBytes(this.P.Cap[:])
	sink.WriteUint64(this.P.Nonce)
	sink.WriteUint64(this.P.StartHeight)
	sink.WriteUint8(this.P.Relay)
	sink.WriteBool(this.P.IsConsensus)
	return nil
}

func NewVersion() Message {
	var version Version
	version.P.Version = 8
	version.P.IsConsensus = false

	return &version
}