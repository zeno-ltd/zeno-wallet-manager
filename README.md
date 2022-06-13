
# zeno-wallet-manager
Will locally use tatum KMS for wallet management
An internal api to generate an address from a protected pre configured wallet,
each address will be associated with a user.

## Creating wallets
To install on local development environment, you will need to generate your own custodial wallet for local testing, the original wallet.dat file will not be shared for security reasons.
```
docker pull tatumio/tatum-kms
docker volume create  --driver local --opt type=none --opt device=$HOME/tatum --opt o=bind tatum
docker run -it -v tatum:/root/.tatumrc -e TATUM_KMS_PASSWORD="Your own local password here" tatumio/tatum-kms --testnet generatemanagedwallet MATIC
docker run -it -v tatum:/root/.tatumrc -e TATUM_KMS_PASSWORD="Your own local password here" tatumio/tatum-kms --testnet getprivatekey {{wallet_id}} //returned in previous step
docker run -it -v tatum:/root/.tatumrc -e TATUM_KMS_PASSWORD="Your own local password here" tatumio/tatum-kms --testnet storemanagedprivatekey MATIC
```
copy the wallet.dat file created in ```$HOME/tatum``` it will required for creating the internal api service for generating wallet addresses for users in the below steps.

More details on tatum kms [here](%28https://docs.tatum.io/tutorials/how-to-securely-store-private-keys#5.-store-the-private-key-to-your-wallet%29)

## Setup (internal api service)
```
git clone git@github.com:sheldondz/zeno-wallet-manager.git
cd zeno-wallet-manager
cp $HOME/tatum/wallet.dat .
sh build.sh
docker run --rm -d  -p 4000:4000 -e TATUM_KMS="/usr/src/app/dist/index.js" -e NODE_EXEC="/usr/local/bin/node" -e HTTP_KMS_PORT=0.0.0.0:4000 -e TATUM_KMS_PASSWORD="same password used above"  zeno-wallet-manager
```
This will run an api wrapped around the tatum-kms, which will only be used to generate addresses for businesses and users who signup. 
>**Note:**  in staging and production environments this service will only be exposed to the zeno-api in the same VPC and not exposed to the internet.

