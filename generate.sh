#!/bin/bash


for (( c=0; c<100000; c++ ))
do 
   address=$(docker run -it   -v tatum:/root/.tatumrc -e TATUM_KMS_PASSWORD="e1e351df-26ca-40e6-8804-4e2f5cd7ba22" tatumio/tatum-kms --testnet getaddress 85fc7f7c-175d-4b3d-9a5c-e1c0b5f5ed44 $c | jq -r ".address")
   echo "$address, $c" >> addresses.txt
done

