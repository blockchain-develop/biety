package validator

import "github.com/ontio/ontology-eventbus/actor"

type Validator interface {
	Register(poolId  *actor.PID)
	UnRegister(poolId *actor.PID)
}

type validator struct {
	pid *actor.PID
	id string
}

func NewValidator(id string) (Validator, error) {
	validator := &validator{
		id : id,
	}

	props := actor.FromProducer(func() actor.Actor {return validator})

	pid, err := actor.SpawnNamed(props, id)
	validator.pid = pid
	return validator, err
}

func (self *validator) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *CheckTx:
		errCode := VerifyTransaction(msg.Tx)
		response := &CheckResponse {
			ErrCode:errCode,
			Hash:msg.Tx.Hash(),
			Height:0,
		}

		sender := context.Sender()
		sender.Tell(response)
	}
}

func (self *validator) Register(poolId *actor.PID) {
	rv := &RegisterValidator{
		Sender : self.pid,
		Id : self.id,
	}
	poolId.Tell(rv)
}

func (self *validator) UnRegister(poolId *actor.PID) {
	urv := &UnRegisterValidator{
		Id : self.id,
	}
	poolId.Tell(urv)
}