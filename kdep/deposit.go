package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kjx98/go-eth"
	"github.com/kjx98/golib/to"
)

// 0xf499... out, expired
// erc20@kraken   0xf499de5d77d511c8b7d3102978c5ca2cba40e0d5
// ETH@kraken	  0xeb8f5d4f02e15441282408c822d8931f5f2d9670
// ETH@hlp		0xfEb8F1aA128FC340A91cDDF5FA23dEcc7C329E23
var (
	useNative     bool
	waitTx        bool
	useHLP        bool
	gasPriceLimit float64
	tipLimit      float64
	nonce         uint64
)

func main() {
	flag.BoolVar(&useNative, "eth", true, "deposit to ETH address")
	flag.BoolVar(&useHLP, "hlp", false, "deposit to HLP ETH address")
	flag.BoolVar(&waitTx, "wait", false, "wait for TransactionReceipt")
	flag.Uint64Var(&nonce, "nonce", 0, "special nonce tx ro replace")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: kdep [options] <acct/hash> [value] [gasPriceLimit]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if err := eth.NewInfura(); err != nil {
		log.Fatal(err)
	}

	if len(flag.Args()) == 0 {
		log.Fatal("Account/Hash missing")
	}
	if waitTx {
		// call waitTx
		txHash := common.HexToHash(flag.Arg(0))
		waitConfirm(txHash)
		os.Exit(0)
	}
	if gasPrice, err := eth.GasPrice(); err == nil {
		gasPriceLimit = gasPrice // normal, fast set to 110%
	} else {
		log.Fatal("Get gasPrice:", err)
	}
	if tipCap, err := eth.GasTipCap(); err == nil {
		//tipLimit = tipCap * 1.1 // faster tip
		tipLimit = tipCap * 0.5 // slower tip
	} else {
		log.Fatal("Get gaTipCap:", err)
	}
	if len(flag.Args()) > 2 {
		gasPriceLimit = to.Double(flag.Arg(2))
		tipLimit = 3.0
	}
	fromAddr := common.HexToAddress(flag.Arg(0))
	toAddr := common.HexToAddress("0xf499de5d77d511c8b7d3102978c5ca2cba40e0d5")
	if useHLP {
		toAddr = common.HexToAddress("0xfEb8F1aA128FC340A91cDDF5FA23dEcc7C329E23")
		tipLimit = 0
	} else if useNative {
		toAddr = common.HexToAddress("eb8f5d4f02e15441282408c822d8931f5f2d9670")
	}
	if nonce == 0 {
		if n, err := eth.PendingNonce(fromAddr); err == nil {
			nonce = n
		}
	}
	if acct, err := eth.Find(fromAddr); err == nil {
		pwd := utils.GetPassPhrase("unlock acct "+fromAddr.String(), false)
		if err := eth.Unlock(acct, pwd); err != nil {
			log.Fatal("Unlock failed:", err)
		}
	} else {
		log.Fatal("No such account: ", fromAddr)
	}

	var vETH float64
	var gasLimit uint64
	if ret, err := eth.EstGas(fromAddr, toAddr); err == nil {
		gasLimit = ret
	} else {
		gasLimit = 65000
	}
	feeETH := float64(gasLimit) * gasPriceLimit * 0.000000001
	if len(flag.Args()) < 2 {
		if ret, err := eth.PendingBalance(fromAddr); err == nil {
			vETH = ret - feeETH
		} else {
			log.Fatal("get PendingBalance:", err)
		}
	} else {
		vETH = to.Double(flag.Arg(1))
	}

	fmt.Printf("Deposit %.4f(ETH) from %v \n", vETH, fromAddr)
	fmt.Printf("Use gasPrice: %.4f,  TipCap: %.4f\n", gasPriceLimit, tipLimit)
	fmt.Printf("Tx fee %.8f ETH\n", feeETH)

	if tx, err := eth.NewTx(fromAddr, toAddr, vETH, nonce, gasLimit,
		gasPriceLimit, tipLimit); err != nil {
		log.Fatal("NewTX: ", err)
	} else {
		txHash := eth.SendTx(tx)
		fmt.Printf("Deposit tx: %s\n", txHash.Hex())
	}
}

func waitConfirm(txHash common.Hash) {
	tEnd := time.Now().Unix() + 300 // 5 minutes
	for time.Now().Unix() < tEnd {
		if res, err := eth.TransactionReceipt(txHash); err == nil {
			// dump confirm
			fmt.Printf("%s mined @%d gasUsed %d\n", txHash.Hex(),
				res.BlockNumber.Uint64(), res.GasUsed)
			if res.EffectiveGasPrice != nil {
				fmt.Printf("GasPrice: %.8f\n", to.ToGWei(res.EffectiveGasPrice.Uint64()))
			}
			break
		} else {
			fmt.Println("Tx confirmation", err)
		}
		// sleep 5 seconds
		time.Sleep(5 * time.Second)
	}
}
