package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	method := flag.String("m", "", "Method to be executed")
	key := flag.String("key", "", "Wallet key")
	dest := flag.String("dest", "", "Destination Address")
	amount := flag.Int64("amount", 0, "Transaction amount (Satoshi)")
	txHash := flag.String("txid", "", "Previous transaction hash on blockchain")
	txIndex := flag.Uint("txindex", 0, "Index of UTXO in the specified transaction to be used as input")
	pkScript := flag.String("pkscript", "", "PubKey Script")

	flag.Parse()

	switch *method {
	case "create-wallet":
		if *key == "" {
			log.Fatalf("Key must be provided to create new wallet")
		}
		var wallet Wallet
		err := wallet.DecryptFile(*key)
		if err == nil {
			log.Fatalf("Wallet already exists")
		}
		err = wallet.Generate(*key)
		if err != nil {
			log.Fatalln(err, "Error!")
		}
		return

	case "get-wallet":
		if *key == "" {
			log.Fatalf("Key must be provided to obtain wallet")
		}
		var wallet Wallet
		err := wallet.GetAddresses(*key)
		if err != nil {
			log.Fatalln(err, "Error!")
		}
		fmt.Println("Compressed Address:", wallet.CompressedAddress)
		fmt.Println("Uncompressed Address:", wallet.UncompressedAddress)
		return

	case "destroy-wallet":
		if *key == "" {
			log.Fatalf("Key must be provided to destroy wallet")
		}
		var wallet Wallet
		err := wallet.Destroy(*key)
		if err != nil {
			log.Fatalln(err, "Error!")
		}
		return

	case "post-transaction":
		if *key == "" {
			log.Fatalf("Key must be provided to post transaction")
		}
		if *dest == "" {
			log.Fatalf("Destination Address must be provided to post transaction")
		}
		if *amount == 0 {
			log.Fatalf("Non-zero amount must be provided to post transaction")
		}
		if *txHash == "" {
			log.Fatalf("Previous transaction hash must be provided to post transaction")
		}
		if *pkScript == "" {
			log.Fatalf("PK Script must be provided to post transaction")
		}
		var wallet Wallet
		err := wallet.DecryptFile(*key)
		if err != nil {
			log.Fatalln(err, "Error!")
		}
		transaction, err := CreateTx(wallet.WIF, *dest, *amount, *txHash, *pkScript, *txIndex)
		if err != nil {
			log.Fatalln(err, "Error!")
		}
		fmt.Println("Signed Transaction:", transaction.SignedTx)
		fmt.Println("Unsigned Transaction:", transaction.UnsignedTx)
		fmt.Println("TX Hash:", transaction.TxId)
		return

	default:
		log.Fatalf("Unsupported method")
	}
}
