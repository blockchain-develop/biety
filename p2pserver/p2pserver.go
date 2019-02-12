package p2pserver

import (
	"fmt"
	"github.com/biety/base"
	"github.com/biety/config"
	"net"
	"strconv"
	"time"
)

type P2PServer struct {
	msgRouter     *MessageRouter

	synclistener  net.Listener
	conslistener  net.Listener

	SyncChan    chan *MsgPayload
	ConsChan    chan *MsgPayload
}

func NewServer() *P2PServer {
	p := &P2PServer {
		SyncChan:make(chan *MsgPayload, 10000),
		ConsChan:make(chan *MsgPayload, 10000),
	}
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
	this.msgRouter.init(this.SyncChan, this.ConsChan)
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
	fmt.Printf("connect %s successful.\n", addr)

	//
	remotepeer := NewPeer()
	remotepeer.SyncLink.SetAddr(addr)
	remotepeer.SyncLink.SetConn(conn)
	remotepeer.SyncLink.SetChan(this.SyncChan)
	go remotepeer.SyncLink.Rx()

	fmt.Printf("try version......\n")
	version := NewVersion()
	err = remotepeer.Send(version, isConsensus)
	if err != nil {
		fmt.Printf("send version error: %s\n", err)
		return err
	}

	return nil
}
