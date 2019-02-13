package p2pserver

import (
	"fmt"
	"github.com/biety/block"
	"github.com/biety/common"
	"io"
)

type BlkHeader struct {
	BlkHdr   []*block.Header
}


//Serialize message payload
func (this BlkHeader) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteUint32(uint32(len(this.BlkHdr)))

	for _, header := range this.BlkHdr {
		err := header.Serialization(sink)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *BlkHeader) CmdType() string {
	return HEADERS_TYPE
}

//Deserialize message payload
func (this *BlkHeader) Deserialization(source *common.ZeroCopySource) error {
	var count uint32
	count, eof := source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}

	for i := 0; i < int(count); i++ {
		var headers block.Header
		err := headers.Deserialization(source)
		if err != nil {
			return fmt.Errorf("deserialze BlkHeader error: %v", err)
		}
		this.BlkHdr = append(this.BlkHdr, &headers)
	}
	return nil
}

func NewBlkHeaders(headers []*block.Header) Message {
	var blkHdr  BlkHeader
	blkHdr.BlkHdr = headers

	return &blkHdr
}


type AppendHeaders struct {
	FromID    uint64
	Headers   []*block.Header
}