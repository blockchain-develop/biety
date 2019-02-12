package p2pserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

//msg cmd const
const (
	MSG_CMD_LEN      = 12               //msg type length in byte
	CMD_OFFSET       = 4                //cmd type offet in msg hdr
	CHECKSUM_LEN     = 4                //checksum length in byte
	MSG_HDR_LEN      = 24               //msg hdr length in byte
	MAX_BLK_HDR_CNT  = 500              //hdr count once when sync header
	MAX_INV_HDR_CNT  = 500              //inventory count once when req inv
	MAX_REQ_BLK_ONCE = 16               //req blk count once from one peer when sync blk
	MAX_MSG_LEN      = 30 * 1024 * 1024 //the maximum message length
	MAX_PAYLOAD_LEN  = MAX_MSG_LEN - MSG_HDR_LEN
)

type Message interface {
	Serialization(sink *ZeroCopySink) (err error)
	Deserialization(source *ZeroCopySource) error
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
	CMD          [MSG_CMD_LEN]byte
	Length       uint32
	Checksum     [CHECKSUM_LEN]byte
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

	checksum := Checksum(buf)
	if checksum != hdr.Checksum {
		return nil, 0, fmt.Errorf("message checksum mismatch\n")
	}

	cmdType := string(bytes.TrimRight(hdr.CMD[:], string(0)))
	msg, err := MakeEmptyMessage(cmdType)
	if err != nil {
		return nil, 0, err
	}


	source := NewZeroCopySource(buf)
	err = msg.Deserialization(source)
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


func WriteMessage(sink* ZeroCopySink, msg Message) (err error) {
	pstart := sink.Size()
	sink.NextBytes(MSG_HDR_LEN)
	err = msg.Serialization(sink)
	if err != nil {
		return err
	}

	pend := sink.Size()
	total := pend - pstart
	payLen := total - MSG_HDR_LEN

	sink.BackUp(total)
	buf := sink.NextBytes(total)
	checksum := Checksum(buf[MSG_HDR_LEN:])
	hdr := newMessageHeader(msg.CmdType(), uint32(payLen), checksum)

	sink.BackUp(total)
	writeMessageHeaderInto(sink, hdr)
	sink.NextBytes(payLen)

	return nil
}

func writeMessageHeaderInto(sink *ZeroCopySink, msgh messageHeader) {
	sink.WriteUint32(msgh.Magic)
	sink.WriteBytes(msgh.CMD[:])
	sink.WriteUint32(msgh.Length)
	sink.WriteBytes(msgh.Checksum[:])
}

func newMessageHeader(cmd string, length uint32, checksum [CHECKSUM_LEN]byte) messageHeader {
	msgh := messageHeader{}
	msgh.Magic = 1
	copy(msgh.CMD[:], cmd)
	msgh.Checksum = checksum
	msgh.Length = length
	return msgh
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
