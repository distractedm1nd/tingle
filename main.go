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
)

const usageStr = `
Usage:
	chat connect <node address> <token> <username>
	chat join <key>
	chat join public <room id>
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
	case "join":
		if len(os.Args) == 4 && os.Args[2] == "public" {
			return JoinCmd(ctx, os.Args[3], true)
		} else if len(os.Args) == 3 {
			return JoinCmd(ctx, os.Args[2], false)
		} else {
			return errors.New("please provide a room id or encryption key")
		}
	case "list":
		return ListCmd(ctx)
	default:
		return fmt.Errorf("unknown command")
	}
}

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
