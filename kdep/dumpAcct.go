package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

func main() {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks := keystore.NewKeyStore(os.Getenv("HOME")+"/bc/eth", scryptN, scryptP)
	var unlockAcct common.Address
	if len(os.Args) > 1 {
		unlockAcct = common.HexToAddress(os.Args[1])
	}
	var acct accounts.Account
	for _, act := range ks.Accounts() {
		fmt.Println("Addr: ", act.Address)
		if act.Address == unlockAcct {
			acct = act
		}
	}
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	pwd := utils.GetPassPhrase("unlock acct "+unlockAcct.String(), false)
	if err := ks.Unlock(acct, pwd); err == nil {
		fmt.Println("Unlock successfully")
	} else {
		fmt.Println("Unlock failed:", err)
	}
	// func (ks *KeyStore) SignTx(a accounts.Account, tx *types.Transaction,
	//						chainID *big.Int) (*types.Transaction, error)
	fmt.Println("vim-go")
}
