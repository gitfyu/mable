package main

import (
	"flag"
	"github.com/gitfyu/mable/internal/mable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

func notExist(file string) bool {
	_, err := os.Stat(file)
	return os.IsNotExist(err)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	var dataDir string
	flag.StringVar(&dataDir, "data-dir", "", "Directory containing the application data")
	flag.Parse()

	if len(dataDir) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	if notExist(dataDir) {
		log.Info().Str("path", dataDir).Msg("Data directory does not exist, creating it now")

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatal().Err(err).Msg("Could not create data directory")
		}
	}

	cfgFile := path.Join(dataDir, "config.toml")
	cfg := mable.DefaultConfig()

	if notExist(cfgFile) {
		log.Warn().Msg("No config file found, using default config")
	} else if err := mable.LoadConfig(cfgFile, cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

}
