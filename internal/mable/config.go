package mable

import (
	"github.com/BurntSushi/toml"
	"os"
)

// Config represents the values stored in the server's configuration file
type Config struct {
	// Address is the address the server will listen on
	Address string `toml:"address"`
	// MaxPacketSize is the maximum size in bytes a packet is allowed to have without being rejected
	MaxPacketSize int `toml:"max_packet_size"`
}

// DefaultConfig returns a Config containing the default configuration values
func DefaultConfig() *Config {
	return &Config{
		Address: ":25565",
		// TODO currently this is just an arbitrarily chosen limit
		MaxPacketSize: 1 << 16,
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
