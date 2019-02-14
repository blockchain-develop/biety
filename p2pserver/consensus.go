package p2pserver

import (
	"bytes"
	"fmt"
	"github.com/biety/common"
	"github.com/biety/signature"
	"github.com/ontio/ontology-crypto/keypair"
	"io"
)

type ConsensusPayload struct {
	Version         uint32
	PrevHash        common.Uint256
	Height          uint32
	BookkeeperIndex uint16
	Timestamp       uint32
	Data            []byte
	Owner           keypair.PublicKey
	Signature       []byte
	PeerId          uint64
	hash            common.Uint256
}

//get the consensus payload hash
func (this *ConsensusPayload) Hash() common.Uint256 {
	return common.Uint256{}
}

//Check whether header is correct
func (this *ConsensusPayload) Verify() error {
	buf := new(bytes.Buffer)
	err := this.SerializeUnsigned(buf)
	if err != nil {
		return err
	}
	err = signature.Verify(this.Owner, buf.Bytes(), this.Signature)
	if err != nil {
		return fmt.Errorf("signature verify error.")
	}
	return nil
}

//serialize the consensus payload
func (this *ConsensusPayload) ToArray() []byte {
	b := new(bytes.Buffer)
	err := this.Serialize(b)
	if err != nil {
		return nil
	}
	return b.Bytes()
}

//return inventory type
func (this *ConsensusPayload) InventoryType() common.InventoryType {
	return common.CONSENSUS
}

func (this *ConsensusPayload) GetMessage() []byte {
	//TODO: GetMessage
	//return sig.GetHashData(cp)
	return []byte{}
}

func (this *ConsensusPayload) Type() common.InventoryType {

	//TODO:Temporary add for Interface signature.SignableData use.
	return common.CONSENSUS
}

func (this *ConsensusPayload) Serialization(sink *common.ZeroCopySink) error {
	this.serializationUnsigned(sink)
	buf := keypair.SerializePublicKey(this.Owner)
	sink.WriteVarBytes(buf)
	sink.WriteVarBytes(this.Signature)

	return nil
}

//Serialize message payload
func (this *ConsensusPayload) Serialize(w io.Writer) error {
	err := this.SerializeUnsigned(w)
	if err != nil {
		return err
	}
	buf := keypair.SerializePublicKey(this.Owner)
	err = common.WriteVarBytes(w, buf)
	if err != nil {
		return fmt.Errorf("write publickey error.")
	}

	err = common.WriteVarBytes(w, this.Signature)
	if err != nil {
		return fmt.Errorf("write Signature error.")
	}

	return nil
}

//Deserialize message payload
func (this *ConsensusPayload) Deserialization(source *common.ZeroCopySource) error {
	err := this.deserializationUnsigned(source)
	if err != nil {
		return err
	}
	buf, _, irregular, eof := source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	this.Owner, err = keypair.DeserializePublicKey(buf)
	if err != nil {
		return fmt.Errorf("deserialize publickey error.")
	}

	this.Signature, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}

//Deserialize message payload
func (this *ConsensusPayload) Deserialize(r io.Reader) error {
	err := this.DeserializeUnsigned(r)
	if err != nil {
		return err
	}
	buf, err := common.ReadVarBytes(r)
	if err != nil {
		return fmt.Errorf("read buf error.")
	}
	this.Owner, err = keypair.DeserializePublicKey(buf)
	if err != nil {
		return fmt.Errorf("deserialize publickey error.")
	}

	this.Signature, err = common.ReadVarBytes(r)
	if err != nil {
		return fmt.Errorf("read Signature error.")
	}

	return err
}

func (this *ConsensusPayload) serializationUnsigned(sink *common.ZeroCopySink) {
	sink.WriteUint32(this.Version)
	sink.WriteHash(this.PrevHash)
	sink.WriteUint32(this.Height)
	sink.WriteUint16(this.BookkeeperIndex)
	sink.WriteUint32(this.Timestamp)
	sink.WriteVarBytes(this.Data)
}

//Serialize message payload
func (this *ConsensusPayload) SerializeUnsigned(w io.Writer) error {
	err := common.WriteUint32(w, this.Version)
	if err != nil {
		return fmt.Errorf("write error.")
	}
	err = this.PrevHash.Serialize(w)
	if err != nil {
		return fmt.Errorf("serialize error.")
	}
	err = common.WriteUint32(w, this.Height)
	if err != nil {
		return fmt.Errorf("write error.")
	}
	err = common.WriteUint16(w, this.BookkeeperIndex)
	if err != nil {
		return fmt.Errorf("write error.")
	}
	err = common.WriteUint32(w, this.Timestamp)
	if err != nil {
		return fmt.Errorf("write error.")
	}
	err = common.WriteVarBytes(w, this.Data)
	if err != nil {
		return fmt.Errorf("write error.")
	}
	return nil
}

func (this *ConsensusPayload) deserializationUnsigned(source *common.ZeroCopySource) error {
	var irregular, eof bool
	this.Version, eof = source.NextUint32()
	this.PrevHash, eof = source.NextHash()
	this.Height, eof = source.NextUint32()
	this.BookkeeperIndex, eof = source.NextUint16()
	this.Timestamp, eof = source.NextUint32()
	this.Data, _, irregular, eof = source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	return nil
}

//Deserialize message payload
func (this *ConsensusPayload) DeserializeUnsigned(r io.Reader) error {
	var err error
	this.Version, err = common.ReadUint32(r)
	if err != nil {
		return fmt.Errorf("read version error.")
	}

	preBlock := new(common.Uint256)
	err = preBlock.Deserialize(r)
	if err != nil {
		return fmt.Errorf("read preBlock error.")
	}
	this.PrevHash = *preBlock

	this.Height, err = common.ReadUint32(r)
	if err != nil {
		return fmt.Errorf("read Height error.")
	}

	this.BookkeeperIndex, err = common.ReadUint16(r)
	if err != nil {
		return fmt.Errorf("read BookkeeperIndex error.")
	}

	this.Timestamp, err = common.ReadUint32(r)
	if err != nil {
		return fmt.Errorf("read Timestamp error.")
	}

	this.Data, err = common.ReadVarBytes(r)
	if err != nil {
		return fmt.Errorf("read  Data error.")
	}

	return nil
}

type Consensus struct {
	Cons  ConsensusPayload
}

//Serialize message payload
func (this *Consensus) Serialization(sink *common.ZeroCopySink) error {
	return this.Cons.Serialization(sink)
}

func (this *Consensus) CmdType() string {
	return CONSENSUS_TYPE
}

//Deserialize message payload
func (this *Consensus) Deserialization(source *common.ZeroCopySource) error {
	return this.Cons.Deserialization(source)
}

