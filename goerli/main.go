package main

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx := context.Background()

	// initialize client
	cli, err := ethclient.DialContext(ctx, os.Getenv("GOERLI_RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}

	key, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	chainId := big.NewInt(5)
	log.Printf("chain id: %s\n", chainId.String())

	from := crypto.PubkeyToAddress(key.PublicKey)
	log.Printf("sender: %s\n", from.Hex())

	bn, _ := cli.BlockNumber(ctx)
	log.Printf("block number: %d\n", bn)

	// next nonce of the account
	nonce, err := cli.PendingNonceAt(ctx, from)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("nonce: %d\n", nonce)

	// priority fee price cap
	tipCap, _ := cli.SuggestGasTipCap(ctx)
	log.Printf("tip cap: %s\n", tipCap.String())

	// base fee price cap
	feeCap, _ := cli.SuggestGasPrice(ctx)
	log.Printf("fee cap: %s\n", feeCap.String())

	// value to send
	value := big.NewInt(1000) // in wei (0.000000000000001 eth)
	log.Printf("sending %s wei to my own address\n", value.String())

	// gas limit
	gasLimit := uint64(21000)
	log.Printf("gas limit: %d\n", gasLimit)

	// create tx
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainId,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        &from,
		Value:     value,
		Data:      []byte{},
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), key)
	if err != nil {
		log.Fatal(err)
	}

	// send tx
	if err := cli.SendTransaction(ctx, signedTx); err != nil {
		log.Fatal(err)
	}

	log.Println("waiting for tx to be mined...")

	receipt, err := bind.WaitMined(ctx, cli, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tx hash: %s\n", receipt.TxHash.Hex())
}
