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
	dir      string
	Addr     string
	Token    string
	Username string
}

func NewConfig(addr, token string) *Config {
	return &Config{
		Addr:  addr,
		Token: token,
	}
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
	c.dir = path
	return c, err
}

func (c *Config) Save() error {
	f, err := os.Create(c.dir)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	err = encoder.Encode(c)
	return err
}

func ConnectCmd(_ context.Context, addr, token, username string) error {
	c := NewConfig(addr, token)
	err := c.Save()
	if err != nil {
		return err
	}
	return nil
}
