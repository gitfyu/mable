package mable

import (
	"github.com/BurntSushi/toml"
	"os"
)

// Config represents the values stored in the server's configuration file
type Config struct {
	Address string `toml:"address"`
}

// DefaultConfig returns a Config containing the default configuration values
func DefaultConfig() *Config {
	return &Config{
		Address: ":25565",
	}
}

// LoadConfig decodes a TOML file into a Config
func LoadConfig(file string, cfg *Config) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = toml.NewDecoder(f).Decode(cfg)
	return err
}
