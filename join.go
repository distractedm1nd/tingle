package main

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

func JoinCmd(ctx context.Context, key string, public bool) error {
	celestiaClient, err := NewCelestiaClient(ctx)
	if err != nil {
		return err
	}

	m, err := NewModel(ctx, celestiaClient, key, public)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}
