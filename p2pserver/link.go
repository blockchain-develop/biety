package p2pserver

import (
	"bufio"
	"fmt"
	"net"
)

type Link struct {
	addr           string
	conn           net.Conn

	recvChan       chan *MsgPayload
}

func NewLink()  *Link {
	link := &Link {

	}

	return link
}

func (this *Link) Rx() {
	conn := this.conn
	if conn == nil {
		return
	}

	reader := bufio.NewReaderSize(conn, 1024 * 256)

	for {
		msg, payloadsize, err := ReadMessage(reader)
		if err != nil {
			fmt.Printf("read err : %s", err)
			break
		}

		this.recvChan <- &MsgPayload {
			Addr :  this.addr,
			PayloadSize: payloadsize,
			Payload: msg,
		}
	}
}












func (this *Link) SetAddr(addr string) {
	this.addr = addr
}

func (this *Link) GetAddr() string {
	return this.addr
}

func (this *Link) SetConn(conn net.Conn) {
	this.conn = conn
}

func (this *Link) GetConn() net.Conn {
	return this.conn
}

func (this *Link) SetChan(recvchan chan *MsgPayload) {
	this.recvChan = recvchan
}

func (this *Link) GetChan() chan *MsgPayload {
	return this.recvChan
}
