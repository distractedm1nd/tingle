package main

import (
	"context"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
)

func ListCmd(ctx context.Context) error {
	client, err := NewCelestiaClient(ctx)
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
	return listRooms(GetMessagesBackwardsAsync(ctx, client, chatNamespace, earliestHeight, currentHeight)), nil
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
	for ID, _ := range doesIDExist {
		IDs[i] = ID
		i++
	}
	return IDs
}
