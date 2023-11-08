package eth

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kjx98/golib/to"
)

type callMsg = ethereum.CallMsg

// erc20@kraken   0xf499de5d77d511c8b7d3102978c5ca2cba40e0d5
// ETH@kraken	  0xeb8f5d4f02e15441282408c822d8931f5f2d9670

var client *ethclient.Client

// New client connection
func NewClient(url string) error {
	if c, err := ethclient.Dial(url); err == nil {
		client = c
	} else {
		return err
	}
	return nil
}

func NewInfura() error {
	url := "https://mainnet.infura.io/v3/" + os.Getenv("INFURA_API")
	return NewClient(url)
}

// EstGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func EstGas(fromAddr string) uint64 {
	toAddr := common.HexToAddress("0xf499de5d77d511c8b7d3102978c5ca2cba40e0d5")
	txCall := callMsg{
		From:  common.HexToAddress(fromAddr),
		To:    &toAddr,
		Value: to.FromEther(3.0),
	}
	var ret uint64
	if r, err := client.EstimateGas(context.Background(), txCall); err == nil {
		ret = r
	} else {
		log.Fatal("client.EstimateGas failed:", err)
	}
	return ret
}

func BlockNumber() (num uint64) {
	if ret, err := client.BlockNumber(context.Background()); err != nil {
		log.Fatal("client.BlockNumber failed:", err)
	} else {
		num = ret
	}
	return
}

// GasTipCap retrieves the currently suggested gas tip cap, in GWei
func GasTipCap() float64 {
	var res float64
	if ret, err := client.SuggestGasTipCap(context.Background()); err != nil {
		log.Fatal("client.SuggestGasTipCap:", err)
	} else {
		res = to.ToGWei(ret.Uint64())
	}
	return res
}

// GasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func GasPrice() float64 {
	var res float64
	if ret, err := client.SuggestGasPrice(context.Background()); err != nil {
		log.Fatal("client.SuggestGasPrice:", err)
	} else {
		res = to.ToGWei(ret.Uint64())
	}
	return res
}

func FeeHistory() (base float64, reward []float64) {
	if ret, err := client.FeeHistory(context.Background(), 1, nil,
		[]float64{10, 50, 90}); err != nil {
		log.Fatal("client.FeeHistory failed:", err)
	} else {
		base = to.ToGWei(ret.BaseFee[0].Uint64())
		reward = make([]float64, len(ret.Reward[0]))
		for i, r := range ret.Reward[0] {
			reward[i] = to.ToGWei(r.Uint64())
		}
	}
	return
}

// NewTx  create Tx for transfer and signed
func NewTx(fromAddr, toAddr common.Address, valueEth float64, feeMax float64) *types.Transaction {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Fatal("client.PendingNonceAt failed:", err)
	}

	value := to.FromEther(valueEth) // in wei
	gasLimit := uint64(52680)       // in units, erc20@kraken 52653
	tip := big.NewInt(2000000000)   // maxPriorityFeePerGas = 2 Gwei
	feeCap := new(big.Int)
	feeCap.SetUint64(to.FromGWei(feeMax)) // maxFeePerGas Gwei

	chainID, err := client.ChainID(context.Background())
	//chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: feeCap,
		GasTipCap: tip,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     value,
	})
	acct, err := Find(fromAddr)
	if err != nil {
		log.Fatal(fromAddr.String(), "NOT FOUND:", err)
	}
	signedTx, err := SignTx(acct, tx, chainID)
	if err != nil {
		log.Fatal("SignTx failed:", err)
	}
	return signedTx
}

func SendTx(tx *types.Transaction) string {
	if err := client.SendTransaction(context.Background(), tx); err != nil {
		log.Fatal("client.SendTransaction failed:", err)
	}

	return tx.Hash().Hex()
}
