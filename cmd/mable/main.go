package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	var dataDir string
	flag.StringVar(&dataDir, "data-dir", "", "Directory containing the application data")
	flag.Parse()

	if len(dataDir) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		log.Info().Str("path", dataDir).Msg("Data directory does not exist, creating it now")

		if err := os.MkdirAll(dataDir, os.ModeDir|0600); err != nil {
			log.Fatal().Err(err).Msg("Could not create data directory")
		}
	}
}
