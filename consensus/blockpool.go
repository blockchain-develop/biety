package consensus

type BlockPool struct {

}

func (pool *BlockPool) onBlockSealed(num uint32) {

}

func (pool *BlockPool) setBlockSealed(block *Block, forEmpty bool) error {
	return nil
}

func (pool *BlockPool) NewBlockProposal(msg *blockProposalMsg) error {
	return nil
}

func (pool *BlockPool) NewBlockCommitment(msg *blockCommitMsg) error {
	return nil
}