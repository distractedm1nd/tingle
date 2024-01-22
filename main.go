package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
)

var chatNamespaceStr = "tingle"

const (
	// TODO: make this configurable
	syncPeriod uint64 = 100 // last 100 blocks

	writeKeyLength = 64
	readKeyLength  = 32
	idLength       = 10
)

const usageStr = `
Usage:
	chat create
	chat join <pub>/<priv>/<id>
	chat list 
`

type Message struct {
	ID       string
	Username string
	Content  string
	Public   bool
}

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if len(os.Args) == 1 {
		fmt.Println(usageStr)
		return nil
	}

	switch os.Args[1] {
	case "connect":
		if len(os.Args) != 5 {
			return errors.New("please provide the nodes address and token and a username")
		}
		return ConnectCmd(ctx, os.Args[2], os.Args[3], os.Args[4])
	case "create":
		return CreateCmd()
	case "join":
		if len(os.Args) != 3 {
			return errors.New("please provide a write key, read key, or room id")
		}
		return JoinCmd(ctx, os.Args[2])
	case "list":
		return ListCmd(ctx)
	default:
		return fmt.Errorf("invalid command")
	}
}

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
