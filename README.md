
# zeno-wallet-manager
Will locally use tatum KMS for wallet management


## Creating wallets
To install on local development environment, you will need to generate your own custodial wallet for local testing, the original wallet.dat file will not be shared for security reasons.
```
docker pull tatumio/tatum-kms

docker volume create  --driver local --opt type=none --opt device=$HOME/tatum --opt o=bind tatum

docker run -it -v tatum:/root/.tatumrc -e TATUM_KMS_PASSWORD="Your own local password here" tatumio/tatum-kms --testnet generatemanagedwallet MATIC

docker run -it -v tatum:/root/.tatumrc -e TATUM_KMS_PASSWORD="Your own local password here" tatumio/tatum-kms --testnet getprivatekey {{wallet_id}} //returned in previous step

docker run -it -v tatum:/root/.tatumrc -e TATUM_KMS_PASSWORD="Your own local password here" tatumio/tatum-kms --testnet storemanagedprivatekey MATIC
```
copy the wallet.dat file created in ```$HOME/tatum``` it will be required for creating the internal api service that generates wallet addresses for users/businesses in the below steps.

More details on tatum kms [here](https://docs.tatum.io/tutorials/how-to-securely-store-private-keys#5.-store-the-private-key-to-your-wallet)

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

