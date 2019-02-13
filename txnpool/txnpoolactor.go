package txnpool

type TxActor struct {
	server    *TxPoolServer
}

func (ta *TxActor) setServer(s *TxPoolServer) {
	ta.server = s
}

func NewTxActor(s *TxPoolServer) *TxActor {
	a := &TxActor{}
	a.setServer(s)
	return a
}


type TxPoolActor struct {
	server    *TxPoolServer
}

func (tpa *TxPoolActor) setServer(s *TxPoolServer) {
	tpa.server = s
}

func NewTxPoolActor(s *TxPoolServer) *TxPoolActor {
	a := &TxPoolActor{}
	a.setServer(s)
	return a
}

type VerifyRspActor struct {
	server    *TxPoolServer
}

func (vpa *VerifyRspActor) setServer(s *TxPoolServer) {
	vpa.server = s
}

func NewVerifyRspActor(s *TxPoolServer) *VerifyRspActor {
	a := &VerifyRspActor {}
	a.setServer(s)
	return a
}
