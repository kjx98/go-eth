package eth

import (
	"errors"
	"fmt"
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

func GetKey(addr common.Address, auth string) (*keystore.Key, error) {
	if !inited {
		InitKeyStore("")
	}
	a, err := Find(addr)
	if err != nil {
		return nil, err
	}

	// Load the key from the keystore and decrypt its contents
	keyjson, err := os.ReadFile(a.URL.Path)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	return key, nil
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
