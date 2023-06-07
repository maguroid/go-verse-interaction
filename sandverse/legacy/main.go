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
	cli, err := ethclient.DialContext(ctx, os.Getenv("SAND_VERSE_RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}

	key, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	chainId := big.NewInt(20197)
	log.Printf("chain id: %s\n", chainId.String())

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	if err != nil {
		log.Fatal(err)
	}

	bn, _ := cli.BlockNumber(ctx)
	log.Printf("block number: %d\n", bn)

	nonce, err := cli.PendingNonceAt(ctx, auth.From)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("nonce: %d\n", nonce)

	// gasPrice, err := cli.SuggestGasPrice(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// => 0
	gasPrice := big.NewInt(0)
	log.Printf("gas price: %s\n", gasPrice.String())

	value := big.NewInt(1000) // in wei (0.000000000000001 eth)
	log.Printf("sending %s wei to my own address\n", value.String())

	tx := types.NewTransaction(nonce, auth.From, value, 21000, gasPrice, []byte{})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), key)
	if err != nil {
		log.Fatal(err)
	}

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
