package block

import (
	"errors"
	"github.com/biety/base"
	"github.com/biety/common"
	"github.com/ontio/ontology-crypto/keypair"
)

func AddressFromPubKey(pubkey keypair.PublicKey) common.Address {
	prog := ProgramFromPubKey(pubkey)

	return common.AddressFromVmCode(prog)
}

func AddressFromMultiPubKeys(pubkeys []keypair.PublicKey, m int) (common.Address, error) {
	var addr common.Address
	n := len(pubkeys)
	if !(1 <= m && m <= n && n > 1 && n <= base.MULTI_SIG_MAX_PUBKEY_SIZE) {
		return addr, errors.New("wrong multi-sig param")
	}

	prog, err := ProgramFromMultiPubKey(pubkeys, m)
	if err != nil {
		return addr, err
	}

	return common.AddressFromVmCode(prog), nil
}

func AddressFromBookkeepers(bookkeepers []keypair.PublicKey) (common.Address, error) {
	if len(bookkeepers) == 1 {
		return AddressFromPubKey(bookkeepers[0]), nil
	}
	return AddressFromMultiPubKeys(bookkeepers, len(bookkeepers)-(len(bookkeepers)-1)/3)
}

