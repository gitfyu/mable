// Package config holds the server's configuration values. The values are automatically parsed from command-line
// flags when this package is initialized.
package config

import (
	"flag"
	"github.com/gitfyu/mable/internal/server"
)

var Srv server.Config

func initServerFlags() {
	flag.StringVar(&Srv.Addr, "srv-bind", ":25565", "address to bind to, such as :25565 or 123.123.123.123:123")
	flag.IntVar(&Srv.MaxPacketSize, "srv-max-packet-size", 1<<16, "Maximum size of a single packet, in bytes")
	flag.IntVar(&Srv.Timeout, "srv-timeout", 20, "Time in seconds after which idle clients are kicked")
	flag.StringVar(&Srv.LogLevel, "srv-log-level", "info", "The server will print debug logs")
}

func init() {
	initServerFlags()
	flag.Parse()
}
