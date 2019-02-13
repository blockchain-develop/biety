package p2pserver

import (
	"github.com/biety/common"
	"io"
)

type DataReq struct {
	DataType common.InventoryType
	Hash     common.Uint256
}

//Serialize message payload
func (this DataReq) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteByte(byte(this.DataType))
	sink.WriteHash(this.Hash)

	return nil
}

func (this *DataReq) CmdType() string {
	return GET_DATA_TYPE
}

//Deserialize message payload
func (this *DataReq) Deserialization(source *common.ZeroCopySource) error {
	ty, eof := source.NextByte()
	this.DataType = common.InventoryType(ty)

	this.Hash, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func NewBlkDataReq(hash common.Uint256) Message {
	var dataReq DataReq
	dataReq.DataType = common.BLOCK
	dataReq.Hash = hash

	return &dataReq
}
