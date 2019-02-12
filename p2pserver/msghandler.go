package p2pserver

import "fmt"

func VersionHandle(data* MsgPayload, args ...interface{}) {
	fmt.Printf("receive version message")

}

func VersionAck(data* MsgPayload, args ...interface{}) {

}

func PingHandle(data* MsgPayload, args ...interface{}) {

}

func PongHandle(data* MsgPayload, args ...interface{}) {

}
