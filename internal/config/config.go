package config

import (
	"flag"
	"github.com/gitfyu/mable/internal/server"
)

var Srv server.Config

var DebugLogs bool

func init() {
	flag.StringVar(&Srv.Addr, "bind", ":25565", "address to bind to, such as :25565 or 123.123.123.123:123")
	flag.IntVar(&Srv.MaxPacketSize, "max-packet-size", 1<<16, "Maximum size of a single packet, in bytes")
	flag.IntVar(&Srv.Timeout, "timeout", 20, "Time in seconds after which idle clients are kicked")
	flag.BoolVar(&DebugLogs, "debug-logs", true, "Whether to print debug logs")

	flag.Parse()
}
