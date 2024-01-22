#!/bin/bash

# Get the address to fund from the first argument
ADDRESS_TO_FUND=$1

# Check if there is a second argument, if so set it as the amount of tokens to be sent, else set a default of 1000000
if [ -z "$2" ]
then
    AMOUNT_TO_SEND="10000000utia"
else
    AMOUNT_TO_SEND=$2
fi


# Get the validator's account
VALIDATOR_ACCOUNT=$(celestia-appd keys show validator -a --home temp/celestia)

# Fund the address
celestia-appd tx bank send $VALIDATOR_ACCOUNT $ADDRESS_TO_FUND $AMOUNT_TO_SEND --chain-id private --yes --home temp/celestia --keyring-backend test --fees 30000utia