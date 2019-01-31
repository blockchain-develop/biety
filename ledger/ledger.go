package ledger

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
