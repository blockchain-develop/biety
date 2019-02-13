package p2pserver

import (
	"errors"
	"fmt"
	"github.com/biety/base"
	"github.com/biety/config"
	"github.com/ontio/ontology-eventbus/actor"
	"net"
	"strconv"
	"time"
)

type P2PServer struct {
	pid       *actor.PID

	msgRouter     *MessageRouter

	synclistener  net.Listener
	conslistener  net.Listener

	SyncChan    chan *MsgPayload
	ConsChan    chan *MsgPayload

	PeerSyncAddress map[string]*Peer
	PeerConsAddress map[string]*Peer

	Peerindex       int32
	PeerList        map[int32]*Peer
}

func NewServer() *P2PServer {
	p := &P2PServer {
		SyncChan:make(chan *MsgPayload, 10000),
		ConsChan:make(chan *MsgPayload, 10000),
	}
	p.PeerConsAddress = make(map[string]*Peer)
	p.PeerSyncAddress = make(map[string]*Peer)

	p.Peerindex = 0
	p.PeerList = make(map[int32]*Peer)

	p.msgRouter = NewMsgRouter()
	return p
}

func (this *P2PServer) Start() error {
	//
	err := this.startNet()
	if err != nil {
		return err
	}

	//
	this.msgRouter.init(this.SyncChan, this.ConsChan, this)
	err = this.msgRouter.start()
	if err != nil {
		return err
	}

	//
	go this.connectSeedService()

	return nil
}

func (this *P2PServer) startNet() error {
	syncPort := config.Sync_port
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(syncPort))
	if err != nil {
		return err
	}
	this.synclistener = listener
	go this.StartSyncAccept(this.synclistener)

	consPort := config.Cons_port
	listener, err = net.Listen("tcp", ":"+strconv.Itoa(consPort))
	if err != nil {
		return err
	}
	this.conslistener = listener
	go this.StartConsAccept(this.conslistener)

	return nil
}

func (this *P2PServer) StartSyncAccept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("StartAccept err : %s", err)
			return
		}

		addr := conn.RemoteAddr().String()
		remotepeer := NewPeer()
		remotepeer.SyncLink.SetAddr(addr)
		remotepeer.SyncLink.SetConn(conn)
		remotepeer.SyncLink.SetChan(this.SyncChan)
		this.AddPeerSyncAddress(addr, remotepeer)
		fmt.Printf("%s connect to me\n", addr)

		go remotepeer.SyncLink.Rx()
	}
}

func (this *P2PServer) StartConsAccept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("StartAccept err : %s", err)
			return
		}

		addr := conn.RemoteAddr().String()
		remotepeer := NewPeer()
		remotepeer.ConsLink.SetAddr(addr)
		remotepeer.ConsLink.SetConn(conn)
		remotepeer.ConsLink.SetChan(this.ConsChan)
		this.AddPeerConsAddress(addr, remotepeer)
		fmt.Printf("%s connect to me\n", addr)

		go remotepeer.ConsLink.Rx()
	}
}


func (this *P2PServer) connectSeedService() {
	t := time.NewTimer(time.Second * base.CONN_MONITOR)
	for {
		select {
		case <-t.C:
			this.connectSeeds()
			t.Stop()
		}
	}
}

func (this *P2PServer) connectSeeds() {
	seedNodes := make([]string, 0)
	//pList := make([]*Peer, 0)
	for _, n := range config.DefaultConfig.Genesis.SeedList {
		ip, err := base.ParseIPAddr(n)
		if err != nil {
			fmt.Printf("seed peer %s address format is wrong", n)
			continue
		}

		ns,err := net.LookupHost(ip)
		if err != nil {
			fmt.Printf("resolve err: %s", err)
			continue
		}

		port, err := base.ParseIPPort(n)
		if err != nil {
			fmt.Printf("seed peer %s address format is wrong", n)
			continue
		}

		seedNodes = append(seedNodes, ns[0] + port)
	}

	for _,nodeaddr := range seedNodes {
		go this.Connect(nodeaddr, false)
	}
}

func (this *P2PServer) Connect(addr string, isConsensus bool) error {
	fmt.Printf("try to connect %s......\n", addr)
	conn, err := net.DialTimeout("tcp", addr, time.Second*base.DIAL_TIMEOUT)
	if err != nil {
		fmt.Printf("connect %s failed:%s\n", addr, err)
		return err
	}
	addr = conn.RemoteAddr().String()
	fmt.Printf("connect %s successful.\n", addr)

	//
	var remotepeer *Peer
	if !isConsensus {
		remotepeer = NewPeer()
		remotepeer.SyncLink.SetAddr(addr)
		remotepeer.SyncLink.SetConn(conn)
		remotepeer.SyncLink.SetChan(this.SyncChan)
		this.AddPeerSyncAddress(addr, remotepeer)
		go remotepeer.SyncLink.Rx()
		remotepeer.SetSyncState(HAND)
	} else {
		remotepeer = NewPeer()
		remotepeer.SyncLink.SetAddr(addr)
		remotepeer.SyncLink.SetConn(conn)
		remotepeer.SyncLink.SetChan(this.SyncChan)
		this.AddPeerConsAddress(addr, remotepeer)
		go remotepeer.SyncLink.Rx()
		remotepeer.SetConsState(HAND)
	}

	fmt.Printf("try version......\n")
	version := NewVersion()
	err = remotepeer.Send(version, isConsensus)
	if err != nil {
		fmt.Printf("send version error: %s\n", err)
		return err
	}

	return nil
}

func (this *P2PServer) Send(p *Peer, msg Message, isConsensus bool) error {
	if p != nil {
		return p.Send(msg, isConsensus)
	}

	return errors.New("send to a invalid peer")
}

func (this *P2PServer) AddPeerSyncAddress(addr string, p* Peer) {
	this.PeerSyncAddress[addr] = p

	this.PeerList[this.Peerindex] = p
	this.Peerindex ++
}

func (this *P2PServer) AddPeerConsAddress(addr string, p *Peer) {
	this.PeerConsAddress[addr] = p

	this.PeerList[this.Peerindex] = p
	this.Peerindex ++
}

func (this *P2PServer) GetPeerFromAddr(addr string) *Peer {
	p, ok := this.PeerSyncAddress[addr]
	if ok {
		return p
	}

	p, ok = this.PeerConsAddress[addr]
	if ok {
		return p
	}

	return nil
}

func (this *P2PServer) SetPID(pid *actor.PID) {
	this.pid = pid
	this.msgRouter.SetPID(pid)
}

func (this *P2PServer) GetNode(id int32) *Peer {
	n, ok := this.PeerList[id]
	if !ok {
		return nil
	}

	return n
}