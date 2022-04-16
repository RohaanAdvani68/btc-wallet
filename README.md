# BTC Wallet
Lightweight cold storage Wallet solution that allows for receiving and sending BTC on Testnet. This project consists of btc-wallet-cold, which provides a CLI to perform private key operations, and btc-wallet-hot, which provides a CLI for networked operations, including obtaining account balance, obtaining transactions details, and broadcasting raw transactions to Testnet.

## Installation Guide
Once this repo is cloned, navigate to the *btc-wallet-cold* directory. Run the following command:  
`go build -o lib/main .`  
  
Next, navigate to the *btc-wallet-hot* directory and run the same command:  
 `go build -o lib/main .`

 We will be using the Blockcypher API to interact with the BTC Testnet. You may obtain a free Blockcypher API Key from [here](https://accounts.blockcypher.com). In the *btc-wallet-hot* directory, create a .env file with the following environment variable:
 * BLOCK_CYPHER_API_KEY={your_api_key}

You are now all set to start interacting with the applications via the command line.

## Wallet Interactions

### Create wallet 
`./btc-wallet-cold/lib/main -m=create-wallet -key={secret_key}`  
This command creates a new cold wallet, if wallet does not already exist in storage. key is the password used to generate the private key of the wallet, and so must be kept secure. 

### Get wallet address
`./btc-wallet-cold/lib/main -m=get-wallet -key={secret_key}`  
This command obtains wallet adrress (compressed, uncompressed) from storage.  
You can now try receiving funds to the newly created uncompressed address using any Testnet BTC faucet, such as [this one](https://testnet-faucet.mempool.co). 

### Destroy wallet 
`./btc-wallet-cold/lib/main -m=destroy-wallet -key={secret_key}`  
This command destroys an existing wallet, if wallet exists in storage.

### Get wallet balance
`./btc-wallet-hot/lib/main -m=get-balance -addr={addr_hash}`

### Send BTC to another wallet
1. Obtain Transactions for account  
`./btc-wallet-hot/lib/main -m=get-txs -addr={addr_hash} -confirmed=true`
2. Obtain detailed transaction information for selected transaction  
`./btc-wallet-hot/lib/main -m=get-tx -tx-hash={tx_hash}`  
From here, we can obtain detailed information for UTXO.
3. Create raw transaction hex  
`./btc-wallet-cold/lib/main -m=post-tx -key={secret_key} -dest={destination_address} -amount={amount_satoshi} -txindex={tx_index} -pkscript={pk_script}`  
Here we will create the raw signed transaction hex, which is to be broadcast to the nextwork in the next step. The unsigned transaction hex is also printed to the screen which can be decoded using any online tool to debug transaction data.  
Note, if specified amount is less than UTXO unspent balance, the difference goes to the miner as we have not implemented change via multiple outputs.  
*tx_index* and *pk_script* can be obtained from the detailed tx details in step 2. 

4. Broadcast transaction to BTC Testnet  
`./btc-wallet-hot/lib/main -m=push-tx -tx-hex={tx_hex}`  
This broadcasts the raw hex encoded transaction to the BTC Testnet. *tx_hex* is simply the output from step 3 above. 

## References 
* https://github.com/LuisAcerv/btchdwallet
* https://github.com/nraboy/open-ledger-micro