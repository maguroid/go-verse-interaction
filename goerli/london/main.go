package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maguroid/go-verse-interaction/lib/counter"
)

var chainId = big.NewInt(5)

func main() {
	ctx := context.Background()

	// initialize client
	cli, err := ethclient.DialContext(ctx, os.Getenv("GOERLI_RPC_URL"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	key, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := sendTransaction(ctx, cli, key); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := deployContract(ctx, cli, key); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func sendTransaction(ctx context.Context, cli *ethclient.Client, key *ecdsa.PrivateKey) error {
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
	// value := big.NewInt(1000) // in wei (0.000000000000001 eth)
	value := big.NewInt(0) // in wei (0.000000000000001 eth)
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

	return nil
}

func deployContract(ctx context.Context, cli *ethclient.Client, key *ecdsa.PrivateKey) error {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	if err != nil {
		return err
	}

	addr, tx, _, err := counter.DeployCounter(auth, cli)
	if err != nil {
		return err
	}

	log.Println("deploying counter contract...")
	log.Printf("contract address: %s\n", addr.Hex())
	log.Printf("tx hash: %s\n", tx.Hash().Hex())
	log.Println("waiting for tx to be mined...")

	receipt, err := bind.WaitMined(ctx, cli, tx)
	if err != nil {
		return err
	}

	log.Printf("contract deployed at block %d\n", receipt.BlockNumber)

	return nil
}
