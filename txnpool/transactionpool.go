package txnpool

import (
	"github.com/biety/block"
	"github.com/biety/common"
)


type TxAttr struct {
	Height    uint32
	ErrCode   error
}

type TxEntry struct {
	Tx     *block.Transaction
	Attrs  []*TxAttr
}


type TxPool struct {
	txList  map[common.Uint256]*TxEntry
}

func (tp *TxPool) Init() {
	tp.txList = make(map[common.Uint256]*TxEntry)
}