package p2pserver

import (
	"fmt"
	"github.com/biety/block"
	"github.com/biety/common"
)


type Block struct {
	Blk   *block.Block
}


//Serialize message payload
func (this *Block) Serialization(sink *common.ZeroCopySink) error {
	err := this.Blk.Serialization(sink)
	if err != nil {
		return fmt.Errorf("Block Serialization failed")
	}

	return nil
}

func (this *Block) CmdType() string {
	return BLOCK_TYPE
}

//Deserialize message payload
func (this *Block) Deserialization(source *common.ZeroCopySource) error {
	this.Blk = new(block.Block)
	err := this.Blk.Deserialization(source)
	if err != nil {
		return fmt.Errorf("Block Deserialization failed")
	}

	return nil
}

func NewBlock(bk *block.Block) Message {
	var blk  Block
	blk.Blk = bk

	return &blk
}


type AppendBlock struct {
	FromID      uint64
	BlockSize   uint32
	Block       *block.Block
}