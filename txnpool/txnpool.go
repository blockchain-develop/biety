package txnpool

func StartTxnPoolServer() (*TxPoolServer, error) {
	s := NewTxPoolServer()
	return s, nil
}
