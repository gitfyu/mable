package world

import (
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Default is a global default world for testing, will be removed in the future
var Default = NewWorld()

// BroadcastDefault broadcasts a random value to the Default world every second, only used for testing and will be
// removed in the future
func BroadcastDefault() {
	for range time.Tick(time.Second) {
		Default.Broadcast(rand.Int())
	}
}

type World struct {
	listeners     map[uint32]chan<- interface{}
	listenersLock sync.RWMutex
}

func NewWorld() *World {
	return &World{
		listeners: make(map[uint32]chan<- interface{}),
	}
}

var subIdCounter uint32

// Subscribe subscribes the provided channel to receive world updates. This function returns a unique ID for this
// subscription, which can be passed to Unsubscribe
func (w *World) Subscribe(ch chan<- interface{}) uint32 {
	w.listenersLock.Lock()
	defer w.listenersLock.Unlock()

	id := atomic.AddUint32(&subIdCounter, 1)
	w.listeners[id] = ch

	return id
}

// Unsubscribe cancels a previous subscription. You should call this function using the ID returned from Subscribe.
func (w *World) Unsubscribe(id uint32) {
	w.listenersLock.Lock()
	defer w.listenersLock.Unlock()

	delete(w.listeners, id)
}

// Broadcast broadcasts a message to all subscribers of this world
func (w *World) Broadcast(msg interface{}) {
	w.listenersLock.RLock()
	defer w.listenersLock.RUnlock()

	for _, ch := range w.listeners {
		ch <- msg
	}

	log.Debug().
		Interface("msg", msg).
		Int("listeners", len(w.listeners)).
		Msg("Broadcast")
}
