package p2pserver

import (
	"github.com/biety/block"
	"github.com/biety/common"
)

// Transaction message
type Trn struct {
	Txn *block.Transaction
}

//Serialize message payload
func (this Trn) Serialization(sink *common.ZeroCopySink) error {
	return this.Txn.Serialization(sink)
}

func (this *Trn) CmdType() string {
	return TX_TYPE
}

//Deserialize message payload
func (this *Trn) Deserialization(source *common.ZeroCopySource) error {
	tx := &block.Transaction{}
	err := tx.Deserialization(source)
	if err != nil {
		return err
	}

	this.Txn = tx
	return nil
}

type TxReq struct {
	Tx    *block.Transaction
}

