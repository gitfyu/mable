// Package main is responsible for starting the server.
package main

import (
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gitfyu/mable/block"
	"github.com/gitfyu/mable/game"
	"github.com/gitfyu/mable/internal/server"
	"github.com/gitfyu/mable/log"
)

var (
	srvConf  server.Config
	gameConf game.Config

	defaultWorld *game.World

	logger = log.Logger{
		Name: "MAIN",
	}
)

func createDefaultWorld() *game.World {
	chunks := make(map[game.ChunkPos]*game.Chunk)
	for x := int32(-2); x <= 2; x++ {
		for z := int32(-2); z <= 2; z++ {
			c := game.NewChunk()
			for dx := uint8(0); dx < 16; dx++ {
				for dz := uint8(0); dz < 16; dz++ {
					for dy := uint8(1); dy < 100; dy += 5 {
						c.SetBlock(dx, dy, dz, block.Stone.ToData())
					}
				}
			}

			chunks[game.ChunkPos{X: x, Z: z}] = c
		}
	}
	return game.NewWorld(chunks)
}

func init() {
	// Server config
	flag.StringVar(&srvConf.Addr, "srv-bind", ":25565", "address to bind to, such as :25565 or 123.123.123.123:123")
	flag.IntVar(&srvConf.MaxPacketSize, "srv-max-packet-size", 1<<16, "Maximum size of a single packet, in bytes")
	flag.IntVar(&srvConf.Timeout, "srv-timeout", 20, "Time in seconds after which idle clients are kicked")
	flag.StringVar(&srvConf.LogLevel, "srv-log-level", "debug", "The minimum level that will be logged")

	// Game config
	flag.IntVar(&gameConf.MaxJobs, "game-max-jobs", 100, "Maximum number of pending jobs")
	tickIntervalStr := flag.String("game-tick-interval", "10ms", "How often a tick should occur")

	var err error
	gameConf.TickInterval, err = time.ParseDuration(*tickIntervalStr)
	if err != nil {
		panic("invalid tick interval value")
	}

	flag.Parse()

	defaultWorld = createDefaultWorld()
}

func main() {
	game := game.NewGame([]*game.World{defaultWorld}, gameConf)
	defer game.Close()

	srv, err := server.NewServer(srvConf, game)
	if err != nil {
		logger.Error("Failed to start").Err(err).Log()
		os.Exit(-1)
	}

	logger.Info("Server started").Log()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := game.Run(); err != nil {
			logger.Error("Game execution failed").Err(err).Log()
			ch <- syscall.SIGINT
		}
	}()

	go func() {
		<-ch
		logger.Info("Shutting down").Log()
		srv.Close()
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, net.ErrClosed) {
		logger.Error("Server execution failed").Err(err).Log()
	}
}
