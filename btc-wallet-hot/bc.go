package main

import (
	"os"

	"github.com/blockcypher/gobcy/v2"
)

func getAddr(addrHash string) (gobcy.Addr, error) {
	bc := gobcy.API{os.Getenv("BLOCK_CYPHER_API_KEY"), "btc", "test3"}
	return bc.GetAddr(addrHash, nil)
}

func pushTx(hex string) (gobcy.TXSkel, error) {
	bc := gobcy.API{os.Getenv("BLOCK_CYPHER_API_KEY"), "btc", "test3"}
	return bc.PushTX(hex)
}
