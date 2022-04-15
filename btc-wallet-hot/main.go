package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

func main() {
	method := flag.String("m", "", "Method to be executed")
	addrHash := flag.String("addr", "", "Wallet Address Hash")
	txHex := flag.String("tx-hex", "", "Raw Transaction Hex to be broadcast to the BTC Testnet")
	confirmed := flag.Bool("confirmed", true, "If true, only return confirmed TXs")
	txHash := flag.String("tx-hash", "", "Transaction hash to look up")

	flag.Parse()

	switch *method {
	case "get-balance":
		if *addrHash == "" {
			log.Fatalf("Address must be provided to get balance")
		}
		addr, err := getAddr(*addrHash)
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.Println("Balance: " + addr.Balance.String() + " satoshi")
		return

	case "get-txs":
		if *addrHash == "" {
			log.Fatalf("Address must be provided to obtain Transactions")
		}
		addr, err := getAddr(*addrHash)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if *confirmed {
			log.Println(addr.NumTX, " confirmed Transaction(s)")
			for i, tx := range addr.TXRefs {
				out, err := json.MarshalIndent(tx, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(i+1, ")\n", string(out))
			}
		} else {
			log.Println(addr.UnconfirmedNumTX, " unconfirmed Transaction(s)")
			for i, tx := range addr.UnconfirmedTXRefs {
				out, err := json.MarshalIndent(tx, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(i+1, ")\n", string(out))
			}
		}
		return

	case "push-tx":
		if *txHex == "" {
			log.Fatalf("Raw Transaction hex must be provided to broadcast to network")
		}
		txSkel, err := pushTx(*txHex)
		if err != nil {
			log.Fatalf(err.Error())
		}
		out, err := json.MarshalIndent(txSkel, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
		return

	case "get-tx":
		if *txHash == "" {
			log.Fatalf("Transaction hash must be provided to obtain tx details")
		}
		tx, err := getTx(*txHash)
		if err != nil {
			log.Fatalf(err.Error())
		}
		out, err := json.MarshalIndent(tx, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
		return
	default:
		log.Fatalf("Unsupported method")
	}
}
