package p2pserver

import (
	"fmt"
	"github.com/biety/block"
	"github.com/biety/common"
	"github.com/ontio/ontology-eventbus/actor"
	"time"
)

func VersionHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {
	fmt.Printf("receive version message")
	version := data.Payload.(*Version)
	remotepeer := p2p.GetPeerFromAddr(data.Addr)
	if remotepeer == nil {
		return
	}

	if version.P.IsConsensus == true {
		remotepeer.UpdateInfo(time.Now(), version.P.Version, version.P.Services, version.P.SyncPort,
			version.P.ConsPort, version.P.Nonce, version.P.Relay, version.P.StartHeight)

		s := remotepeer.GetConsState()
		var msg Message
		if s == INIT {
			remotepeer.SetConsState(HAND_SHAKE)
			msg = NewVersion()
		} else if (s == HAND) {
			remotepeer.SetConsState(HAND_SHAKED)
			msg = NewVerAck(true)
		}

		err := p2p.Send(remotepeer, msg, true)
		if err != nil {
			return
		}
	} else {
		remotepeer.UpdateInfo(time.Now(), version.P.Version, version.P.Services, version.P.SyncPort,
			version.P.ConsPort, version.P.Nonce, version.P.Relay, version.P.StartHeight)

		s := remotepeer.GetSyncState()
		var msg Message
		if s == INIT {
			remotepeer.SetSyncState(HAND_SHAKE)
			msg = NewVersion()
		} else if (s == HAND) {
			remotepeer.SetSyncState(HAND_SHAKED)
			msg = NewVerAck(false)
		}

		err := p2p.Send(remotepeer, msg, false)
		if err != nil {
			return
		}
	}
}

func VersionAck(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {
	fmt.Printf("receive verack message")
	verack := data.Payload.(*VerACK)
	remotepeer := p2p.GetPeerFromAddr(data.Addr)
	if remotepeer == nil {
		return
	}

	if verack.IsConsensus == true {
		s := remotepeer.GetConsState()
		if s != HAND_SHAKE && s != HAND_SHAKED {
			return
		}

		remotepeer.SetConsState(ESTABLISH)
		if s == HAND_SHAKE {
			msg := NewVerAck(true)
			p2p.Send(remotepeer, msg, true)
		}
	} else {
		s := remotepeer.GetSyncState()
		if s != HAND_SHAKE && s != HAND_SHAKED {
			return
		}

		remotepeer.SetSyncState(ESTABLISH)
		if s == HAND_SHAKE {
			msg := NewVerAck(false)
			p2p.Send(remotepeer, msg, false)
		}
	}
}

func PingHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {

}

func PongHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {

}

func HeadersReqHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {
	fmt.Printf("receive headers request message\n")

	headersReq := data.Payload.(*HeadersReq)
	startHash := headersReq.HashStart
	stopHash := headersReq.HashEnd

	headers, err := GetHeadersFromHash(startHash, stopHash)
	if err != nil {
		return
	}

	remotepeer := p2p.GetPeerFromAddr(data.Addr)
	if remotepeer == nil {
		return
	}

	msg := NewBlkHeaders(headers)
	err = p2p.Send(remotepeer, msg, false)
	if err != nil {
		return
	}

	return
}

func BlkHeaderHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {
	fmt.Printf("receive block headers message\n")
	blkHeader := data.Payload.(*BlkHeader)
	input := &AppendHeaders{
		FromID: 1,
		Headers: blkHeader.BlkHdr,
	}

	pid.Tell(input)
}

func DataReqHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {
	fmt.Printf("receive data request message\n")

	dataReq := data.Payload.(*DataReq)
	remotepeer := p2p.GetPeerFromAddr(data.Addr)
	if remotepeer == nil {
		return
	}

	reqType := common.InventoryType(dataReq.DataType)
	hash := dataReq.Hash
	switch reqType {
	case common.BLOCK:
		block, err := GetBlockByHash(hash)
		if err != nil || block == nil || block.Header == nil {
			msg := NewNotFound(hash)
			err := p2p.Send(remotepeer, msg, false)
			if err != nil {
				return
			}
			return
		}
		msg := NewBlock(block)
		err = p2p.Send(remotepeer, msg, false)
		if err != nil {
			return
		}
	}
}

func BlockHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {
	fmt.Printf("receive block message\n")
	block := data.Payload.(*Block)
	input := &AppendBlock {
		FromID: data.Id,
		BlockSize:data.PayloadSize,
		Block:block.Blk,
	}
	pid.Tell(input)
}

func TransactionHandle(data* MsgPayload, p2p *P2PServer, pid *actor.PID, args ...interface{}) {
	fmt.Printf("receive transaction message\n")
	trn := data.Payload.(*Trn)
	AddTransaction(trn.Txn)
}

















func GetHeadersFromHash(startHash common.Uint256, stopHash common.Uint256) ([]*block.Header, error) {
	return nil, nil
}

func GetBlockByHash(uint256 common.Uint256) (*block.Block, error) {
	return nil, nil
}
