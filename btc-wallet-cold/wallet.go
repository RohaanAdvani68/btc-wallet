package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type Wallet struct {
	WIF                 string
	UncompressedAddress string
	CompressedAddress   string
}

// generate new wallet and encrypt, save to disk
func (wallet Wallet) Generate(key string) error {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return err
	}
	wif, err := btcutil.NewWIF(secret, &chaincfg.TestNet3Params, false)
	if err != nil {
		return err
	}
	uncompressedAddress, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return err
	}
	compressedAddress, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return err
	}

	wallet.WIF = wif.String()
	wallet.UncompressedAddress = uncompressedAddress.EncodeAddress()
	wallet.CompressedAddress = compressedAddress.EncodeAddress()

	return wallet.EncryptFile(key)
}

// destroy existing wallet
func (wallet Wallet) Destroy(key string) error {
	if _, err := os.Stat("wallet.dat"); errors.Is(err, os.ErrNotExist) {
		return errors.New("No wallet exists to be destroyed")
	}
	if wallet.Authenticate(key) {
		return os.Remove("wallet.dat")
	}
	return errors.New("Unauthenticated to destroy existing wallet")
}

func (wallet *Wallet) GetAddresses(passphrase string) error {
	if _, err := os.Stat("wallet.dat"); errors.Is(err, os.ErrNotExist) {
		return errors.New("No wallet exists")
	}
	err := wallet.DecryptFile(passphrase)
	if err != nil {
		return errors.New(`Incorrect password`)
	}
	// scrub private key data from response
	wallet.WIF = ""
	return nil
}

// convert key to hex string of appropriate length
func (wallet Wallet) CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// check if passphrase is valid to access wallet
func (wallet *Wallet) Authenticate(passphrase string) bool {
	err := wallet.DecryptFile(passphrase)
	if err != nil {
		return false
	}
	return true
}

func (wallet Wallet) Encrypt(passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(wallet.CreateHash(passphrase)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	data, err := json.Marshal(wallet)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (wallet *Wallet) Decrypt(data []byte, passphrase string) error {
	block, err := aes.NewCipher([]byte(wallet.CreateHash(passphrase)))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	json.Unmarshal(plaintext, wallet)
	return nil
}

// encrypt wallet and save to disk
func (wallet Wallet) EncryptFile(passphrase string) error {
	file, err := os.Create("wallet.dat")
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := wallet.Encrypt(passphrase)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

// read encrypted wallet from disk and decrypt
func (wallet *Wallet) DecryptFile(passphrase string) error {
	ciphertext, err := ioutil.ReadFile("wallet.dat")
	if err != nil {
		return err
	}
	err = wallet.Decrypt(ciphertext, passphrase)
	if err != nil {
		return err
	}
	return nil
}
