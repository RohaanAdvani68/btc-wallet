package main

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Transaction struct {
	TxId               string
	SourceAddress      string
	DestinationAddress string
	Amount             int64 
	UnsignedTx         string
	SignedTx           string 
}

func CreateTransaction(secret string, destination string, amount int64, txHash string) (Transaction, error) {
	var transaction Transaction
	wif, err := btcutil.DecodeWIF(secret)
	if err != nil {
		return Transaction{}, err
	}
	addresspubkey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return Transaction{}, err
	}
	// since we are creating raw tx offline, we will not have access to all UTXO info
	// so we create our own UTXO
	sourceTx := wire.NewMsgTx(wire.TxVersion)
	sourceUtxoHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return Transaction{}, err
	}
	sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
	sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
	destinationAddress, err := btcutil.DecodeAddress(destination, &chaincfg.TestNet3Params)
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.TestNet3Params)
	if err != nil {
		return Transaction{}, err
	}
	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return Transaction{}, err
	}
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return Transaction{}, err
	}
	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)
	sourceTxHash := sourceTx.TxHash()
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	destinationUtxo := wire.NewOutPoint(&sourceTxHash, 0)
	redeemTxIn := wire.NewTxIn(destinationUtxo, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)
	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)
	// sign using sender's pvt key
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return Transaction{}, err
	}
	redeemTx.TxIn[0].SignatureScript = sigScript
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return Transaction{}, err
	}
	if err := vm.Execute(); err != nil {
		return Transaction{}, err
	}
	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)
	transaction.TxId = sourceTxHash.String()
	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
	transaction.Amount = amount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()
	return transaction, nil
}