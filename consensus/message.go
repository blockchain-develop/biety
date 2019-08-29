package consensus

import (
	"bytes"
	"fmt"
	"github.com/biety/block"
	"github.com/biety/common"
	"github.com/biety/p2pserver"
	"io"
)

type StartConsensus struct {}
type StopConsensus struct {}

type p2pMsgPayload struct {
	fromPeer      uint32
	payload       *p2pserver.ConsensusPayload
}

type ConsensusMsg interface {
	Type()    MsgType
	Verify() error
	Serialize()  ([]byte, error)
}

type MsgType  uint8

const (
	BlockProposalMessage MsgType = iota
	BlockEndorseMessage
	BlockCommitMessage

	PeerHandshakeMessage
	PeerHeartbeatMessage

	BlockInfoFetchMessage
	BlockInfoFetchRespMessage
	ProposalFetchMessage
	BlockFetchMessage
	BlockFetchRespMessage
)

type ConsensusMsgPayload struct {
	Type  MsgType `json:"type"`
	Len   uint32  `json:"len"`
	Payload []byte `json:"payload"`
}


type Block struct {
	Block         *block.Block
	EmptyBlock    *block.Block
}


func (blk *Block) getProposer() uint32 {
	return 0
	//return blk.Info.Proposer
}

func (blk *Block) getBlockNum() uint32 {
	return blk.Block.Header.Height
}

func (blk *Block) getPrevBlockHash() common.Uint256 {
	return blk.Block.Header.PrevBlockHash
}

/*
func (blk *Block) getLastConfigBlockNum() uint32 {
	return blk.Info.LastConfigBlockNum
}

func (blk *Block) getNewChainConfig() *vconfig.ChainConfig {
	return blk.Info.NewChainConfig
}
*/

//
// getVrfValue() is a helper function for participant selection.
//
/*
func (blk *Block) getVrfValue() []byte {
	return blk.Info.VrfValue
}

func (blk *Block) getVrfProof() []byte {
	return blk.Info.VrfProof
}
*/

func (blk *Block) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := blk.Block.Serialize(buf); err != nil {
		return nil, fmt.Errorf("serialize block: %s", err)
	}

	payload := bytes.NewBuffer([]byte{})
	if err := common.WriteVarBytes(payload, buf.Bytes()); err != nil {
		return nil, fmt.Errorf("serialize block buf: %s", err)
	}

	if blk.EmptyBlock != nil {
		buf2 := bytes.NewBuffer([]byte{})
		if err := blk.EmptyBlock.Serialize(buf2); err != nil {
			return nil, fmt.Errorf("serialize empty block: %s", err)
		}
		if err := common.WriteVarBytes(payload, buf2.Bytes()); err != nil {
			return nil, fmt.Errorf("serialize empty block buf: %s", err)
		}
	}

	return payload.Bytes(), nil
}

func (blk *Block) Deserialize(data []byte) error {
	source := common.NewZeroCopySource(data)
	//buf := bytes.NewBuffer(data)
	buf1, _, irregular, eof := source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}

	block1, err := block.BlockFromRawBytes(buf1)
	if err != nil {
		return fmt.Errorf("deserialize block: %s", err)
	}

	/*
	info := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(block.Header.ConsensusPayload, info); err != nil {
		return fmt.Errorf("unmarshal vbft info: %s", err)
	}
	*/

	var emptyBlock *block.Block
	if source.Len() > 0 {
		buf2, _, irregular, eof := source.NextVarBytes()
		if irregular == false && eof == false {
			block2, err := block.BlockFromRawBytes(buf2)
			if err == nil {
				emptyBlock = block2
			}
		}
	}

	blk.Block = block1
	blk.EmptyBlock = emptyBlock
	//blk.Info = info

	return nil
}

type blockProposalMsg struct {
	Block  *Block `json:"block"`
}

func (msg *blockProposalMsg) UnmarshalJSON(data []byte) error {
	blk := &Block{}
	if err := blk.Deserialize(data); err != nil {
		return err
	}

	msg.Block = blk
	return nil
}

func (msg *blockProposalMsg) MarshalJSON() ([]byte, error) {
	return msg.Block.Serialize()
}

type blockEndorseMsg struct {

}

type blockCommitMsg struct {

}

type peerHandshakeMsg struct {

}

type peerHeartbeatMsg struct {

}