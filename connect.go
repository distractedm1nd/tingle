package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
)

const pathName = ".chat"

func NewCelestiaClient(ctx context.Context) (*client.Client, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return client.NewClient(ctx, cfg.Addr, cfg.Token)
}

type Config struct {
	path      string
	Addr     string
	Token    string
	Username string
}

func NewConfig(addr, token string) (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, pathName, "config.toml")
	return &Config{
		Addr:  addr,
		Token: token,
		path:  path,
	}, nil
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, pathName, "config.toml")
	c := &Config{}
	_, err = toml.DecodeFile(path, c)
	if err != nil {
		return nil, err
	}
	c.path = path
	return c, err
}

func (c *Config) Save() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	err = encoder.Encode(c)
	return err
}

func ConnectCmd(_ context.Context, addr, token, username string) error {
	c, err := NewConfig(addr, token)
	if err != nil {
		return err
	}
	err = c.Save()
	if err != nil {
		return err
	}
	return nil
}
