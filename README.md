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

If you want to test out Tingle on a local network, you can use the Docker image provided by the `local-celestia-devnet` repository. If you want to use mainnet, you can skip this step.

### Local Celestia Devnet

This Docker image allows developers to spin up a local Celestia devnet composed of a celestia-app validator node and a celestia-node bridge node.

To run the Docker image, use the following command:

```bash
docker run -t -i \
-p 26650:26650 -p 26657:26657 -p 26658:26658 -p 26659:26659 -p 9090:9090 \
ghcr.io/rollkit/local-celestia-devnet:latest
```

For more information, please refer to the [local-celestia-devnet repository](https://github.com/rollkit/local-celestia-devnet).

## Connect Tingle to your node

```bash
./tingle connect http://localhost:26658 <auth key> <username>
```

## Connect Tingle to remote node (running on scaleway)

```bash
./tingle ws://51.159.130.129:26658 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.mo6fGPd2qEfeftaagxPSZNwbQ0wix2PagDPF9FrA-bI <username>
```

Let`s talk !

![Lets talk](./screenshots/Screenshot.png)