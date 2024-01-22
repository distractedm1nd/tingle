# Tingle - A Community Board 💬

Tingle is a community board that offers two types of chat rooms: private and public.

## Private Rooms

Private rooms are created with a unique encryption key. To join a private room, you will need the encryption key provided by the room creator.

## Public Rooms

Public rooms are open to all users. You can join a public room using its public key.

To join either a private or public room, use the `join` command followed by the appropriate key.

## Installation Guide

Before you can use Tingle, you need to install it. Follow these steps:

### Prerequisites

- You need to have Go installed on your machine. If it's not installed, you can download it from [here](https://golang.org/dl/).

### Installation Steps

1. Clone the Tingle repository:

```bash
git clone https://github.com/distractedm1nd/tingle.git
```

2. Navigate to the cloned directory:

```bash
cd tingle
```

3. Install tingle:

```bash
go install
```

After following these steps, Tingle should be installed on your machine and ready to use.

## Usage Guide

If you want to test out Tingle local network, you can use the `localnet.sh` script to set up a local network. If you want to use mainnet, you can skip this step. To run the script install [celestia-app](https://docs.celestia.org/nodes/celestia-app#install-celestia-app).

### localnet.sh

This script is used to set up a local network for testing purposes. It initializes a new blockchain, creates a new account, and configures the network settings.

To run the script, use the following command:

```bash
./scripts/localnet.sh
```

### fund.sh

This script is used to fund an account with tokens. It takes two arguments: the address to fund and the amount of tokens to send. If no amount is specified, it defaults to 10000000utia.

To run the script, use the following command:

```bash
./scripts/fund.sh <address> <amount>
```
