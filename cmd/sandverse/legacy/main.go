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

var chainId = big.NewInt(20197)
var rpcUrl = os.Getenv("SAND_VERSE_RPC_URL")

func main() {
	ctx := context.Background()
	cli, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	key, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("chain id: %s\n", chainId.String())

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
	bn, _ := cli.BlockNumber(ctx)
	log.Printf("block number: %d\n", bn)

	from := crypto.PubkeyToAddress(key.PublicKey)

	nonce, err := cli.PendingNonceAt(ctx, from)
	if err != nil {
		return err
	}
	log.Printf("nonce: %d\n", nonce)

	// gasPrice, err := cli.SuggestGasPrice(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// => 0
	gasPrice := big.NewInt(0)
	log.Printf("gas price: %s\n", gasPrice.String())

	// value := big.NewInt(1000) // in wei (0.000000000000001 eth)
	value := big.NewInt(0) // in wei (0.000000000000001 eth)
	log.Printf("sending %s wei to my own address\n", value.String())

	tx := types.NewTransaction(nonce, from, value, 21000, gasPrice, []byte{})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), key)
	if err != nil {
		return err
	}

	if err := cli.SendTransaction(ctx, signedTx); err != nil {
		return err
	}

	log.Println("waiting for tx to be mined...")

	receipt, err := bind.WaitMined(ctx, cli, signedTx)
	if err != nil {
		return err
	}

	log.Printf("tx hash: %s\n", receipt.TxHash.Hex())

	return nil
}

func deployContract(ctx context.Context, cli *ethclient.Client, key *ecdsa.PrivateKey) error {
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	if err != nil {
		return err
	}

	auth.GasPrice = big.NewInt(0)

	addr, tx, _, err := counter.DeployCounter(auth, cli)
	if err != nil {
		return err
	}

	// addr, tx, _, err := ft.DeployFt(auth, cli)
	// if err != nil {
	// 	return err
	// }

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
