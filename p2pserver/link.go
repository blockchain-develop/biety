package p2pserver

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/biety/common"
	"net"
)

type Link struct {
	addr           string
	conn           net.Conn
	port           uint16

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

		fmt.Printf("recv msg : %v\n", msg.CmdType())

		this.recvChan <- &MsgPayload {
			Addr :  this.addr,
			PayloadSize: payloadsize,
			Payload: msg,
		}
	}
}

func (this *Link) Tx(msg Message) error {
	conn := this.conn
	if conn == nil {
		return errors.New("tx link invalid")
	}

	fmt.Printf("send msg : %v\n", msg.CmdType())

	sink := common.NewZeroCopySink(nil)
	err := WriteMessage(sink, msg)
	if err != nil {
		return err
	}

	payload := sink.Bytes()
	nByteCnt := len(payload)
	fmt.Printf("buf length: %d\n", nByteCnt)

	_, err = conn.Write(payload)
	if err != nil {
		return err
	}

	return nil
}

func (this *Link) SetPort(p uint16) {
	this.port = p
}

func (this *Link) GetPort() uint16 {
	return this.port
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
