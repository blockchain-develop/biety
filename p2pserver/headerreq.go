package p2pserver

import (
	"github.com/biety/common"
	"io"
)

type HeadersReq struct {
	Len        uint8
	HashStart  common.Uint256
	HashEnd    common.Uint256
}

//Serialize message payload
func (this *HeadersReq) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteUint8(this.Len)
	sink.WriteHash(this.HashStart)
	sink.WriteHash(this.HashEnd)
	return nil
}

func (this *HeadersReq) CmdType() string {
	return GET_HEADERS_TYPE
}

//Deserialize message payload
func (this *HeadersReq) Deserialization(source *common.ZeroCopySource) error {
	var eof bool
	this.Len, eof = source.NextUint8()
	this.HashStart, eof = source.NextHash()
	this.HashEnd, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func NewHeadersReq(curHdrHash common.Uint256) Message {
	var h HeadersReq
	h.Len = 1
	h.HashEnd = curHdrHash

	return &h
}
