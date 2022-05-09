package game

import (
	"time"
)

// Config is used to configure a Game instance.
type Config struct {
	// MaxJobs specifies how many jobs can be queued at the same time.
	MaxJobs int
	// TickInterval specifies how often the game state should be updated.
	TickInterval time.Duration
}

// Game manages the state for game related things, such as
// worlds and entities.
type Game struct {
	cfg    Config
	worlds []*World
	closed chan struct{}
	jobs   chan func()
}

// NewGame constructs a new Game. Panics if len(worlds)==0.
func NewGame(worlds []*World, cfg Config) *Game {
	if len(worlds) == 0 {
		panic("no worlds specified")
	}
	return &Game{
		cfg:    cfg,
		worlds: worlds,
		closed: make(chan struct{}),
		jobs:   make(chan func(), cfg.MaxJobs),
	}
}

// Schedule schedules a job to be executed in the same goroutine
// that called Run.
func (g *Game) Schedule(job func()) {
	g.jobs <- job
}

// Run will process game updates until Close is called.
func (g *Game) Run() error {
	ticker := time.NewTicker(g.cfg.TickInterval)
	defer ticker.Stop()

	for {
		select {
		case job := <-g.jobs:
			job()
		case <-ticker.C:
			for _, w := range g.worlds {
				w.tick()
			}
		case <-g.closed:
			return nil
		}
	}
}

// DefaultWorld returns the default World.
func (g *Game) DefaultWorld() *World {
	return g.worlds[0]
}

// Close releases resources associated with the Game.
// Any ongoing Run calls will exit.
// This function may only be called once and always returns nil.
func (g *Game) Close() error {
	close(g.closed)
	return nil
}
