package main

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

func JoinCmd(ctx context.Context, key string) error {
	celestiaClient, err := NewCelestiaClient(ctx)
	if err != nil {
		return err
	}

	m, err := NewModel(ctx, celestiaClient)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}
