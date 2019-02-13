package blocksync

import (
	"fmt"
	"github.com/biety/base"
	"github.com/biety/block"
	"github.com/biety/common"
	"github.com/biety/ledger"
	"github.com/biety/p2pserver"
	"time"
)


type BlockSyncMgr struct {
	server    *p2pserver.P2PServer
	ledger    *ledger.Ledger

	blocksCache map[uint32]*BlockInfo
}

type BlockInfo struct {
	nodeID   uint64
	block    *block.Block
}

func NewBlockSyncMgr(server *p2pserver.P2PServer, ledger *ledger.Ledger) *BlockSyncMgr {
	return &BlockSyncMgr{
		server : server,
		ledger: ledger,
		blocksCache : make(map[uint32]*BlockInfo,0),
	}
}

func (this *BlockSyncMgr) Start() {
	go this.sync()
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <- ticker.C:
			go this.checkTimeout()
			go this.sync()
			go this.saveBlock()
		}
	}
}

func (this *BlockSyncMgr) checkTimeout() {

}

func (this *BlockSyncMgr) sync() {
	this.syncHeader()
	this.syncBlock()
}

func (this *BlockSyncMgr) syncHeader() {
	curBlockHeight := this.ledger.GetCurrentBlockHeight()
	curHeaderHeight := this.ledger.GetCurrentHeaderHeight()
	if curHeaderHeight-curBlockHeight >= base.SYNC_MAX_HEADER_FORWARD_SIZE {
		return
	}

	NextHeaderId := curHeaderHeight + 1
	reqNode := this.getNextNode(NextHeaderId)
	if reqNode == nil {
		return
	}

	headerHash := this.ledger.GetCurrentHeaderHash()
	msg := p2pserver.NewHeadersReq(headerHash)
	err := this.server.Send(reqNode, msg, false)

	if err != nil {
		fmt.Printf("syncHeader failed to send a new headersReq")
	}
}

func (this *BlockSyncMgr) syncBlock() {
	curBlockHeight := this.ledger.GetCurrentBlockHeight()
	curHeaderHeight := this.ledger.GetCurrentHeaderHeight()
	count := curHeaderHeight-curBlockHeight
	if count <= 0 {
		return
	}

	counter := uint32(1)
	i := uint32(0)
	reqTimes := 1
	for {
		if counter > count {
			break
		}
		i ++
		nextBlockHeight := curBlockHeight + i
		nextBlockHash := this.ledger.GetBlockHash(nextBlockHeight)
		if nextBlockHash == common.UINT256_EMPTY {
			return
		}

		for t := 0;t < reqTimes;t ++ {
			reqNode := this.getNextNode(nextBlockHeight)
			if reqNode == nil {
				return
			}

			msg := p2pserver.NewBlkDataReq(nextBlockHash)
			err := this.server.Send(reqNode, msg, false)
			if err != nil {
				fmt.Printf("syncBlock Height error\n")
				return
			}
		}
		counter ++
		reqTimes = 1
	}
}

func (this *BlockSyncMgr) saveBlock() {
	curBlockHeight := this.ledger.GetCurrentBlockHeight()
	nextBlockHeight := curBlockHeight + 1
	for height := range this.blocksCache {
		if height <= curBlockHeight {
			delete(this.blocksCache, height)
		}
	}

	for {
		_, nextBlock := this.getBlockCache(nextBlockHeight)
		if nextBlock == nil {
			return
		}

		this.ledger.AddBlock(nextBlock)
		this.delBlockCache(nextBlockHeight)

		nextBlockHeight ++
	}
}

func (this *BlockSyncMgr) getBlockCache(blockHeight uint32) (uint64, *block.Block) {
	blockInfo, ok := this.blocksCache[blockHeight]
	if !ok {
		return 0, nil
	}

	return blockInfo.nodeID, blockInfo.block
}

func (this *BlockSyncMgr) delBlockCache(blockHeight uint32) {
	delete (this.blocksCache, blockHeight)
}

func (this *BlockSyncMgr) getNextNode(nextBlockHeight uint32) *p2pserver.Peer {
	n := this.server.GetNode(0)
	if n == nil {
		return nil
	}

	return n
}

func (this *BlockSyncMgr) OnHeaderReceive(fromID uint64, headers []*block.Header) {
	if len(headers) == 0 {
		return
	}

	height := headers[0].Height
	curHeadersHeight := this.ledger.GetCurrentHeaderHeight()

	if height < curHeadersHeight {
		return
	}

	err := this.ledger.AddHeader(headers)
	if err != nil {
		return
	}

	this.syncHeader()
}

func (this *BlockSyncMgr) OnBlockReceive(fromID uint64, blockSize uint32, block *block.Block) {
	height := block.Header.Height
	//blockHash := block.Hash()
	curHeaderHeight := this.ledger.GetCurrentHeaderHeight()
	nextHeader := curHeaderHeight + 1
	if height > nextHeader {
		return
	}

	curBlockHeight := this.ledger.GetCurrentBlockHeight()
	if height <= curBlockHeight {
		return
	}

	this.addBlockCache(fromID, block)
	go this.saveBlock()
	this.syncBlock()
}

func (this *BlockSyncMgr) addBlockCache(nodeID uint64, block *block.Block) bool {
	blockInfo := &BlockInfo{
		nodeID : nodeID,
		block: block,
	}

	this.blocksCache[block.Header.Height] = blockInfo
	return true
}
