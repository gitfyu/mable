package main

import (
	"errors"
	"flag"
	"github.com/gitfyu/mable/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
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
	cfg := server.DefaultConfig()

	if notExist(cfgFile) {
		log.Warn().Msg("No config file found, using default config")
	} else if err := server.LoadConfig(cfgFile, cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	if !cfg.DebugLogs {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Warn().Msg("Debug logs are disabled")
	}

	srv, err := server.NewServer(cfg)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	log.Info().Msgf("Listening on %s", srv.Addr().String())

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		<-ch
		log.Info().Msg("Shutting down")
		srv.Close()
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, net.ErrClosed) {
		log.Fatal().Err(err).Msg("Server execution failed")
	}

	log.Info().Msg("Goodbye")
}
