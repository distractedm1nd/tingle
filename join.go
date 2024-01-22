package main

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

func JoinCmd(ctx context.Context, key string, public bool) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	celestiaClient, err := NewCelestiaClient(ctx, cfg)
	if err != nil {
		return err
	}

	m, err := NewModel(ctx, celestiaClient, key, public, cfg.Username)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}
