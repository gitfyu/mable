package game

import (
	"time"
)

// TODO make these configurable

const jobQueueSize = 1000
const entityTickRate = time.Second

type World struct {
	chunks   map[ChunkPos]*Chunk
	entities map[ID]Entity
	jobs     chan func()
	done     chan struct{}
}

// NewWorld constructs a new World
func NewWorld(chunks map[ChunkPos]*Chunk) *World {
	w := &World{
		chunks:   chunks,
		entities: make(map[ID]Entity),
		jobs:     make(chan func(), jobQueueSize),
		done:     make(chan struct{}),
	}
	go w.handle()
	return w
}

func (w *World) AddEntity(e Entity) {
	w.Schedule(func() {
		w.entities[e.GetEntityID()] = e
	})
}

func (w *World) RemoveEntity(id ID) {
	w.Schedule(func() {
		delete(w.entities, id)
	})
}

// Schedule schedules a job to be executed by this world
func (w *World) Schedule(job func()) {
	w.jobs <- job
}

// GetChunk gets the Chunk at the specified position, or nil if it does not exist
func (w *World) GetChunk(pos ChunkPos) *Chunk {
	return w.chunks[pos]
}

func (w *World) handle() {
	entityTicks := time.NewTicker(entityTickRate)
	defer entityTicks.Stop()

	for {
		select {
		case <-entityTicks.C:
			w.tickEntities()
		case job := <-w.jobs:
			job()
		case <-w.done:
			return
		}
	}
}

func (w *World) tickEntities() {
	for _, e := range w.entities {
		e.tick()
	}
}

// Close releases the resources for this World. The Handle function will return after the World is closed. This function
// may only be called once and always returns nil.
func (w *World) Close() error {
	close(w.done)
	return nil
}
