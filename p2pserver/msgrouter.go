package p2pserver

import (
	"fmt"
)

type MessageHandler func(data* MsgPayload, args ...interface{})

type MessageRouter struct {
	msgHandlers map[string]MessageHandler

	RecvSyncChan    chan *MsgPayload
	RecyConsChan    chan *MsgPayload

	stopSyncCh      chan bool
	stopConsCh      chan bool
}

func NewMsgRouter() *MessageRouter {
	msgRouter := &MessageRouter{}
	return msgRouter
}

func (this *MessageRouter) init(syncchan chan *MsgPayload, conschan chan *MsgPayload) {
	this.msgHandlers = make(map[string]MessageHandler)
	this.RecvSyncChan = syncchan
	this.RecyConsChan = conschan
	this.stopSyncCh = make(chan bool)
	this.stopConsCh = make(chan bool)

	this.msgHandlers[VERSION_TYPE] = VersionHandle
	this.msgHandlers[VERACK_TYPE] = VersionAck
	this.msgHandlers[PING_TYPE] = PingHandle
	this.msgHandlers[PONG_TYPE] = PongHandle
}

func (this *MessageRouter) start() error {
	go this.hookChan(this.RecvSyncChan, this.stopSyncCh)
	go this.hookChan(this.RecyConsChan, this.stopConsCh)
	return nil
}

func (this *MessageRouter) hookChan(channel chan *MsgPayload, stopch chan bool) {
	for {
		select {
		case data, ok := <-channel:
			if ok {
				msgType := data.Payload.CmdType()
				handler, ok := this.msgHandlers[msgType]
				if ok {
					go handler(data)
				} else {
					fmt.Printf("unknow message!")
				}
			}
		case <- stopch:
			return
		}
	}
}
