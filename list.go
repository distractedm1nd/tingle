package main

import (
	"context"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/share"
)

func ListCmd(ctx context.Context) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	client, err := NewCelestiaClient(ctx, cfg)
	if err != nil {
		return err
	}
	rooms, err := List(ctx, client)
	if err != nil {
		return err
	}
	for _, room := range rooms {
		println(room)
	}
	return nil
}

func List(ctx context.Context, client *client.Client) ([]string, error) {
	header, err := client.Header.NetworkHead(ctx)
	if err != nil {
		return nil, err
	}
	currentHeight := header.Height()
	earliestHeight := currentHeight - syncPeriod
	namespace, err := share.NewBlobNamespaceV0([]byte(chatNamespaceStr))
	if err != nil {
		return nil, err
	}

	return listRooms(GetMessagesBackwardsAsync(ctx, client, namespace, earliestHeight, currentHeight)), nil
}

// Author: @manav
// Takes in a list of messages, filters them by `isPublic` = true,
// and returns a list of namespace IDs strings.
func listRooms(messages <-chan Message) []string {
	doesIDExist := make(map[string]bool)
	for message := range messages {
		if message.Public {
			doesIDExist[message.ID] = true
		}
	}
	IDs := make([]string, len(doesIDExist))
	i := 0
	for ID := range doesIDExist {
		IDs[i] = ID
		i++
	}
	return IDs
}
