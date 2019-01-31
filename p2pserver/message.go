package p2pserver

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	VERSION_TYPE     = "version"    //peer`s information
	VERACK_TYPE      = "verack"     //ack msg after version recv
	GetADDR_TYPE     = "getaddr"    //req nbr address from peer
	ADDR_TYPE        = "addr"       //nbr address
	PING_TYPE        = "ping"       //ping  sync height
	PONG_TYPE        = "pong"       //pong  recv nbr height
	GET_HEADERS_TYPE = "getheaders" //req blk hdr
	HEADERS_TYPE     = "headers"    //blk hdr
	INV_TYPE         = "inv"        //inv payload
	GET_DATA_TYPE    = "getdata"    //req data from peer
	BLOCK_TYPE       = "block"      //blk payload
	TX_TYPE          = "tx"         //transaction
	CONSENSUS_TYPE   = "consensus"  //consensus payload
	GET_BLOCKS_TYPE  = "getblocks"  //req blks from peer
	NOT_FOUND_TYPE   = "notfound"   //peer can`t find blk according to the hash
	DISCONNECT_TYPE  = "disconnect" //peer disconnect info raise by link
)

type Message interface {
	//Serialization(sink *comm.ZeroCopySink) error
	Deserialization(data []byte) error
	CmdType() string
}

//MsgPayload in link channel
type MsgPayload struct {
	Id          uint64  //peer ID
	Addr        string  //link address
	PayloadSize uint32  //payload size
	Payload     Message //msg payload
}

type messageHeader struct {
	Magic        uint32
	CMD          [12]byte
	Length       uint32
	Checksum     [4]byte
}


func ReadMessage(reader io.Reader) (Message, uint32, error) {
	hdr, err := readMessageHeader(reader)
	if err != nil {
		return nil, 0, err
	}

	buf := make([]byte, hdr.Length)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return nil, 0, err
	}

	cmdType := string(bytes.TrimRight(hdr.CMD[:], string(0)))
	msg, err := MakeEmptyMessage(cmdType)
	if err != nil {
		return nil, 0, err
	}

	err = msg.Deserialization(buf)
	if err != nil {
		return nil, 0, err
	}
	return msg, hdr.Length, nil
}

func readMessageHeader(reader io.Reader) (messageHeader, error) {
	msgh := messageHeader{}
	err := binary.Read(reader, binary.LittleEndian, &msgh)
	return msgh, err
}

















func MakeEmptyMessage(cmdType string) (Message, error) {
	switch cmdType {
	case VERSION_TYPE:
		return &Version{},nil
	case VERACK_TYPE:
		return &VerACK{},nil
	case PING_TYPE:
		return &Ping{},nil
	case PONG_TYPE:
		return &Pong{},nil
	default:
		return nil, nil
	}
}
