package eth

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"os"
)

var ks *keystore.KeyStore
var inited = false

func InitKeyStore(ksdir string) {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if ksdir == "" {
		ksdir = os.Getenv("HOME") + "/bc/eth"
	}
	ks = keystore.NewKeyStore(ksdir, scryptN, scryptP)
	inited = true
}

func Accounts() []accounts.Account {
	if !inited {
		InitKeyStore("")
	}
	return ks.Accounts()
}

func Find(addr common.Address) (accounts.Account, error) {
	if !inited {
		InitKeyStore("")
	}
	for _, act := range ks.Accounts() {
		if act.Address == addr {
			return act, nil
		}
	}
	return accounts.Account{}, errors.New("Address NOT FOUND")
}

func Unlock(acct accounts.Account, pwd string) error {
	if !inited {
		InitKeyStore("")
	}
	return ks.Unlock(acct, pwd)
}

func SignTx(a accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if !inited {
		InitKeyStore("")
	}
	return ks.SignTx(a, tx, chainID)
}

func SignTxWithPassphrase(a accounts.Account, passphrase string,
	tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if !inited {
		InitKeyStore("")
	}
	return ks.SignTxWithPassphrase(a, passphrase, tx, chainID)
}
