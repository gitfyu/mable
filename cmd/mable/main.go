package main

import (
	"errors"
	"github.com/gitfyu/mable/internal/config"
	"github.com/gitfyu/mable/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if !config.DebugLogs {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Warn().Msg("Debug logs are disabled")
	}

	srv, err := server.NewServer(config.Srv)

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
