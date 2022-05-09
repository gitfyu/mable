// Package main is responsible for starting the server.
package main

import (
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gitfyu/mable/internal/server"
	"github.com/gitfyu/mable/log"
)

var (
	srvConf server.Config
	logger  = log.Logger{
		Name: "MAIN",
	}
)

func init() {
	flag.StringVar(&srvConf.Addr, "srv-bind", ":25565", "address to bind to, such as :25565 or 123.123.123.123:123")
	flag.IntVar(&srvConf.MaxPacketSize, "srv-max-packet-size", 1<<16, "Maximum size of a single packet, in bytes")
	flag.IntVar(&srvConf.Timeout, "srv-timeout", 20, "Time in seconds after which idle clients are kicked")
	flag.StringVar(&srvConf.LogLevel, "srv-log-level", "info", "The server will print debug logs")
	flag.Parse()
}

func main() {
	srv, err := server.NewServer(srvConf)

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
