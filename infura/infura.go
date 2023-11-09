package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kjx98/go-eth"
)

// erc20@kraken   0xf499de5d77d511c8b7d3102978c5ca2cba40e0d5
// ETH@kraken	  0xeb8f5d4f02e15441282408c822d8931f5f2d9670

func main() {
	if err := eth.NewInfura(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Lastest blockNumber: %d\n", eth.BlockNumber())
	gasPrice, _ := eth.GasPrice()
	tipCap, _ := eth.GasTipCap()
	fmt.Printf("gasPrice: %8.3f,  TipCap: %8.3f\n", gasPrice, tipCap)

	base, reward := eth.FeeHistory()
	fmt.Printf("Base Gas price: %8.3f\nRewards: ", base)
	//for _, r := range reward {
	//	fmt.Printf("%8.3f, ", r)
	//}
	//fmt.Println("ok")
	fmt.Printf("%8.3f(fast), %8.3f(normal), %8.3f(slow)\n", reward[2],
		reward[1], reward[0])
	if len(os.Args) > 1 {
		acct := common.HexToAddress(os.Args[1])
		bal, err := eth.Balance(acct)
		if err != nil {
			log.Fatal("Get Balance: ", err)
		} else {
			fmt.Printf("%v latest balance: %.8f\n", acct, bal)
		}
	}
}
