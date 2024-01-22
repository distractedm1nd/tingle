package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/share"
)

// @HliBussy
func GetMessagesBackwardsAsync(ctx context.Context, client *client.Client, namespace share.Namespace, startHeight, endHeight uint64) <-chan Message {
	msgCh := make(chan Message, 8)
	go func() {
		defer close(msgCh)
		for height := endHeight; height >= startHeight; height-- {
			blobs, err := client.Blob.GetAll(ctx, height, []share.Namespace{namespace})
			if err != nil {
				slog.Error("can't get blobs for height and namespace, skipping...", height, namespace)
				continue
			}

			for _, b := range blobs {
				var msg Message
				err := json.Unmarshal(b.Data, &msg)
				if err != nil {
					slog.Error("can't unmarshal msg for height and namespace, skipping...", height, namespace)
					continue
				}
				select {
				case msgCh <- msg:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return msgCh
}
