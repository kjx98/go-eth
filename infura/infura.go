package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kjx98/go-eth"
)

// erc20@kraken   0xf499de5d77d511c8b7d3102978c5ca2cba40e0d5
// ETH@kraken	  0xeb8f5d4f02e15441282408c822d8931f5f2d9670

func main() {
	if err := eth.NewInfura(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Lastest blockNumber: %d\n", eth.BlockNumber())
	fmt.Printf("gasPrice: %8.3f,  TipCap: %8.3f\n", eth.GasPrice(),
		eth.GasTipCap())

	if len(os.Args) > 1 {
		fmt.Printf("EstGas for transfer: %d\n", eth.EstGas(os.Args[1]))
	}

	base, reward := eth.FeeHistory()
	fmt.Printf("Base Gas price: %8.3f\nRewards: ", base)
	//for _, r := range reward {
	//	fmt.Printf("%8.3f, ", r)
	//}
	//fmt.Println("ok")
	fmt.Printf("%8.3f(fast), %8.3f(normal), %8.3f(slow)\n", reward[2],
		reward[1], reward[0])
}
