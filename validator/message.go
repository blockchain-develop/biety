package validator

import (
	"github.com/biety/block"
	"github.com/biety/common"
	"github.com/ontio/ontology-eventbus/actor"
)


type RegisterValidator struct {
	Sender    *actor.PID
	Id        string
}

type UnRegisterValidator struct {
	Id        string
}

type CheckTx struct {
	Tx       *block.Transaction
}

type CheckResponse struct {
	Hash     common.Uint256
	Height   uint32
	ErrCode  error
}
