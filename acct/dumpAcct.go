package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kjx98/go-eth"
	"os"
)

func main() {
	for _, act := range eth.Accounts() {
		fmt.Println("Addr: ", act.Address)
	}
	if len(os.Args) > 1 {
		unlockAc := common.HexToAddress(os.Args[1])
		if acct, err := eth.Find(unlockAc); err == nil {
			pwd := utils.GetPassPhrase("unlock acct "+unlockAc.String(), false)
			if err := eth.Unlock(acct, pwd); err == nil {
				fmt.Println("Unlock successfully")
				if key, err := eth.GetKey(unlockAc, pwd); err != nil {
					fmt.Println("Decrypt error:", err)
				} else {
					privKey := hex.EncodeToString(key.PrivateKey.D.Bytes())
					fmt.Println("Secret Key: ", privKey)
				}
			} else {
				fmt.Println("Unlock failed:", err)
			}
		} else {
			fmt.Printf("No such account: %s\n", unlockAc.String())
		}
	}
	// func (ks *KeyStore) SignTx(a accounts.Account, tx *types.Transaction,
	//						chainID *big.Int) (*types.Transaction, error)
}
