package validator

import (
	"errors"
	"fmt"
	"github.com/biety/base"
	"github.com/biety/block"
	"github.com/biety/common"
	"github.com/biety/signature"
)

func VerifyTransaction(tx *block.Transaction) error {
	err := checkTransactionSignatures(tx)
	if err != nil {
		return err
	}

	err = checkTransactionPayload(tx)
	if err != nil {
		return err
	}

	return nil
}

func checkTransactionSignatures(tx *block.Transaction) error {
	hash := tx.Hash()
	lensig := len(tx.Sigs)

	if lensig > base.TX_MAX_SIG_SIZE {
		return fmt.Errorf("transaction signature number")
	}

	address := make(map[common.Address]bool, len(tx.Sigs))
	for _, sigdata := range tx.Sigs {
		sig, err := sigdata.GetSig()
		if err != nil {
			return err
		}

		m := int(sig.M)
		kn := len(sig.PubKeys)
		sn := len(sig.SigData)

		if kn > base.MULTI_SIG_MAX_PUBKEY_SIZE || sn < m {
			return errors.New("wrong tx sig param length")
		}

		if kn == 1 {
			err := signature.Verify(sig.PubKeys[0], hash[:], sig.SigData[0])
			if err != nil {
				return errors.New("signature verification failed")
			}
			address[block.AddressFromPubKey(sig.PubKeys[0])] = true
		} else {
			err := signature.VerifyMultiSignature(hash[:], sig.PubKeys, m, sig.SigData)
			if err != nil {
				return err
			}

			addr, err := block.AddressFromMultiPubKeys(sig.PubKeys, m)
			if err != nil {
				return err
			}
			address[addr] = true
		}
	}

	addrList := make([]common.Address, 0, len(address))
	for addr := range address {
		addrList = append(addrList, addr)
	}

	tx.SignedAddr = addrList
	return nil
}

func checkTransactionPayload(tx *block.Transaction) error {
	switch tx.Payload.(type) {
	case *block.DeployCode:
		return nil
	case *block.InvokeCode:
		return nil
	default:
		return errors.New(fmt.Sprint("unimplemented transaction payload type."))
	}
}
