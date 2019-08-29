package consensus

type MsgPool struct {

}

func (pool *MsgPool) onBlockSealed(num  uint32) {

}

func (pool *MsgPool) GetProposalMsgs(num uint32) []ConsensusMsg {
	return nil
}

func (pool *MsgPool) GetCommitMsgs(num uint32) []ConsensusMsg {
	return nil
}
