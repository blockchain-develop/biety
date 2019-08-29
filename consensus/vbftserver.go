package consensus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/biety/p2pserver"
	"github.com/ontio/ontology-eventbus/actor"
	"math"
)

type BftActionType uint8

const (
	MakeProposal BftActionType = iota
	EndorseBlock
	CommitBlock
	SealBlock
	FastForward // for syncer catch up
	ReBroadcast
)

type BftAction struct {
	Type   BftActionType
	BlockNum      uint32
	Proposal      *blockProposalMsg
}

type SendMsgEvent struct {
	ToPeer uint32 // peer index
	Msg    ConsensusMsg
}

type Server struct {
	Index             uint32

	pid               *actor.PID
	txpoolactorpid    *actor.PID
	p2pactorpid       *actor.PID

	msgRecvC          map[uint32]chan *p2pMsgPayload

	syncer            *Syncer
	stateMgr          *StateMgr

	bftActionC        chan *BftAction
	msgSendC          chan *SendMsgEvent

	msgPool           *MsgPool
	blockPool         *BlockPool
}

func NewVbftServer() (*Server, error) {
	server := &Server {

	}

	props := actor.FromProducer(func() actor.Actor {return server})
	pid, err := actor.SpawnNamed(props, "consensus_vbft")
	if err != nil {
		return nil, err
	}
	server.pid = pid
	return server, nil
}

func (self *Server) Receive(context actor.Context) {
	switch msg:= context.Message().(type) {
	case *StartConsensus:
	case *StopConsensus:
		self.stop()
	case *p2pserver.ConsensusPayload:
		self.NewConsensusPayload(msg)
	}
}


func (self *Server) GetPID() *actor.PID {
	return self.pid
}

func (self *Server) Start() error {
	return self.start()
}

func (self *Server) Init(p2pactorpid *actor.PID, txpoolactorpid *actor.PID) error {
	return self.init(p2pactorpid, txpoolactorpid)
}

func (self *Server) init(p2pactorpid *actor.PID, txpoolactorpid *actor.PID) error {
	self.txpoolactorpid = txpoolactorpid
	self.p2pactorpid = p2pactorpid

	self.msgRecvC = make(map[uint32]chan *p2pMsgPayload)
	self.msgRecvC[0] = make(chan *p2pMsgPayload, 1024)

	for i,_ := range self.msgRecvC {
		self.run(i)
	}

	go self.syncer.run()
	go self.stateMgr.run()
	go self.msgSendLoop()
	go self.actionLoop()

	return nil
}


func (self *Server) start() error {
	return nil
}

func (self *Server) stop() error {
	return nil
}

func (self *Server) NewConsensusPayload(payload *p2pserver.ConsensusPayload) {
	pl := &p2pMsgPayload{
		payload: payload,
	}

	self.msgRecvC[0] <- pl
}

func (self *Server) run(i uint32) {
	go func() {
		for {
			_, msgData, err := self.receiveFromPeer(i)
			if err != nil {
				return
			}

			msg, err := self.DeserializeVbftMsg(msgData)
			if err != nil {
				return
			}

			self.onConsensusMsg(msg)
		}
	}()
}

func (self *Server) receiveFromPeer(id uint32) (uint32, []byte, error) {
	c, present := self.msgRecvC[id]
	if present {
		select {
		case payload := <-c:
			if payload != nil {
				return payload.fromPeer, payload.payload.Data, nil
			}
		}
	}

	return 0, nil, fmt.Errorf("nil consensus payload")
}

func (self *Server) DeserializeVbftMsg(msgPayload []byte) (ConsensusMsg, error) {
	m := &ConsensusMsgPayload{}
	err := json.Unmarshal(msgPayload, m)
	if err != nil {
		return nil, err
	}

	if m.Len < uint32(len(m.Payload)) {
		return nil, err
	}

	switch m.Type {
	case BlockProposalMessage:
		t := &blockProposalMsg{}
		err := t.UnmarshalJSON(m.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal msg")
		}
		return t, nil
	case BlockEndorseMessage:
		t := &blockEndorseMsg{}
		if err := json.Unmarshal(m.Payload, t); err != nil {
			return nil, fmt.Errorf("failed to unmarshal msg (type: %d): %s", m.Type, err)
		}
		return t, nil
	case BlockCommitMessage:
		t := &blockCommitMsg{}
		err := json.Unmarshal(m.Payload, t)
		if err != nil {
			return nil, err
		}
		return t, nil
	case PeerHandshakeMessage:
		t := &peerHandshakeMsg{}
		err := json.Unmarshal(m.Payload, t)
		if err != nil {
			return nil, err
		}
		return t, nil
	case PeerHeartbeatMessage:
		t := &peerHeartbeatMsg{}
		err := json.Unmarshal(m.Payload, t)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, nil
}

func (self *Server) onConsensusMsg(msg ConsensusMsg) {
	switch msg.Type() {
	case BlockProposalMessage:
	case BlockEndorseMessage:
	case BlockCommitMessage:
	case PeerHeartbeatMessage:
	}
}


func (self *Server) actionLoop() {
	for {
		select {
		case action := <-self.bftActionC:
			switch action.Type {
			case MakeProposal:
			case EndorseBlock:
			case CommitBlock:
			case SealBlock:
			case FastForward:
				var blkNum uint32 = 0
				pmsgs := self.msgPool.GetProposalMsgs(blkNum)
				for _, msg := range pmsgs {
					p := msg.(*blockProposalMsg)
					if p != nil {
						self.blockPool.NewBlockProposal(p)
					}
				}

				cmsgs := self.msgPool.GetCommitMsgs(blkNum)
				commitMsgs := make([]*blockCommitMsg, 0)
				for _, msg := range cmsgs {
					c := msg.(*blockCommitMsg)
					if c != nil {
						self.blockPool.NewBlockCommitment(c)
						commitMsgs = append(commitMsgs, c)
					}
				}

				if len(pmsgs) == 0 && len(cmsgs) == 0 {
					self.startNewRound()
					break
				}

				proposer, forEmpty := self.getCommitConsensus(commitMsgs)
				if proposer == math.MaxUint32 {
					break
				}

				var proposal *blockProposalMsg
				for _, m := range pmsgs {
					p, ok := m.(*blockProposalMsg)
					if !ok {
						continue
					}

					if p.Block.getProposer() == proposer {
						proposal = p
						break
					}
				}

				if proposal == nil {
					break
				}

				self.sealBlock(proposal.Block, forEmpty)
			case ReBroadcast:
			}
		}
	}
}

func (self *Server) msgSendLoop() {
	for {
		select {
		case evt := <-self.msgSendC:
			if self.nonConsensusNode() {
				continue
			}
			payload, err := self.SerializeVbftMsg(evt.Msg)
			if err != nil {
				continue
			}

			if evt.ToPeer == math.MaxUint32 {
				self.broadcastToAll(payload)
			} else {
				self.sendToPeer(evt.ToPeer, payload)
			}
		}
	}
}

func (self *Server) nonConsensusNode() bool {
	return self.Index == math.MaxUint32
}


func (self *Server) SerializeVbftMsg(msg ConsensusMsg) ([]byte, error) {
	payload, err := msg.Serialize()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&ConsensusMsgPayload{
		Type:    msg.Type(),
		Len:     uint32(len(payload)),
		Payload: payload,
	})
}

func (self *Server) broadcastToAll(data []byte) error {
	msg := &p2pserver.ConsensusPayload{
		Data: data,
	}

	buf := new(bytes.Buffer)
	err := msg.SerializeUnsigned(buf)
	if err != nil {
		return err
	}

	// sign

	self.p2pactorpid.Tell(msg)
	return nil
}

func (self *Server) sendToPeer(index uint32, data []byte) error {
	return nil
}

func (self *Server) makeFastForward() error {
	ba := &BftAction{
		Type: FastForward,
		BlockNum:0,
	}
	self.bftActionC <- ba
	return nil
}

func (self *Server) sealBlock(block *Block, empty bool) error {
	var sealedBlkNum uint32 = 0

	self.blockPool.setBlockSealed(block, empty)
	self.msgPool.onBlockSealed(sealedBlkNum)
	self.blockPool.onBlockSealed(sealedBlkNum)

	return nil
}

func (self *Server) startNewRound() error {
	return nil
}

func (self *Server) getCommitConsensus(commitMsgs []*blockCommitMsg) (uint32, bool) {
	return 0, true
}



