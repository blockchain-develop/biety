package p2pserver

import (
	"fmt"
	"net"
	"strconv"
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

	return nil
}

func (this *P2PServer) startNet() error {
	syncPort := 6666
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(syncPort))
	if err != nil {
		return err
	}
	this.synclistener = listener
	go this.StartSyncAccept(this.synclistener)

	consPort := 7777
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

		go remotepeer.ConsLink.Rx()
	}
}
