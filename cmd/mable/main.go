// Package main is responsible for starting the server.
package main

import (
	"errors"
	"github.com/gitfyu/mable/internal/config"
	"github.com/gitfyu/mable/internal/server"
	"github.com/gitfyu/mable/log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var logger = log.Logger{
	Name: "MAIN",
}

func main() {
	srv, err := server.NewServer(config.Srv)

	if err != nil {
		logger.Error("Failed to start").Err(err).Log()
		os.Exit(-1)
	}

	logger.Info("Server started").Log()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		<-ch
		logger.Info("Shutting down").Log()
		srv.Close()
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, net.ErrClosed) {
		logger.Error("Server execution failed").Err(err).Log()
	}
}
