package p2pserver

import (
	"github.com/biety/common"
	"io"
)


type NotFound struct {
	Hash    common.Uint256
}

//Serialize message payload
func (this NotFound) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteHash(this.Hash)
	return nil
}

func (this NotFound) CmdType() string {
	return NOT_FOUND_TYPE
}

//Deserialize message payload
func (this *NotFound) Deserialization(source *common.ZeroCopySource) error {
	var eof bool
	this.Hash, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func NewNotFound(hash common.Uint256) Message {
	var notfound NotFound
	notfound.Hash = hash

	return &notfound
}
