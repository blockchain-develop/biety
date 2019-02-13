package ledger

import (
	"github.com/biety/block"
	"github.com/biety/common"
)

type Ledger struct {
	dir     string
}

func NewLedger(dir string) (*Ledger, error) {
	return &Ledger {
		dir: dir,
	}, nil
}

func (self *Ledger) Close() error {
	return nil
}

func (self *Ledger) AddHeader(headers []*block.Header) error {
	return nil
}

func (self *Ledger) AddBlock(block *block.Block) error {
	return nil
}

func (self *Ledger) GetBLockByHeight(height uint32) (*block.Block, error) {
	return nil, nil
}

func (self *Ledger) GetBlockByHash(blockhash common.Uint256) (*block.Block, error) {
	return nil, nil
}

func (self *Ledger) GetHeaderByHeight(height uint32) (*block.Header, error) {
	return nil, nil
}

func (self *Ledger) GetHeaderByHash(blockhash common.Uint256) (*block.Header, error) {
	return nil,nil
}

func (self *Ledger) GetBlockHash(height uint32) common.Uint256 {
	return common.UINT256_EMPTY
}

func (self *Ledger) GetCurrentBlockHeight() uint32{
	return 0
}

func (self *Ledger) GetCurrentBlockHash() common.Uint256 {
	return common.UINT256_EMPTY
}

func (self *Ledger) GetCurrentHeaderHeight() uint32 {
	return 0
}

func (self *Ledger) GetCurrentHeaderHash() common.Uint256 {
	return common.UINT256_EMPTY
}
