package world

import (
	"math/rand"
	"sync"
	"time"
)

// Default is a global default world for testing, will be removed in the future
var Default = &World{
	listeners: make(map[Listener]struct{}),
}

func init() {
	go func() {
		for range time.Tick(time.Second) {
			Default.Broadcast(rand.Int())
		}
	}()
}

type Listener interface {
	// OnWorldUpdate is currently invoked with an arbitrary int value, in the future this function will receive
	// data about things such as entity movement, block changes, etc.
	OnWorldUpdate(v int)
}

type World struct {
	listeners     map[Listener]struct{}
	listenersLock sync.RWMutex
}

func (w *World) Subscribe(l Listener) {
	w.listenersLock.Lock()
	defer w.listenersLock.Unlock()

	w.listeners[l] = struct{}{}
}

func (w *World) Unsubscribe(l Listener) {
	w.listenersLock.Lock()
	defer w.listenersLock.Unlock()

	delete(w.listeners, l)
}

// Broadcast currently just broadcasts an int value to each Listener
func (w *World) Broadcast(v int) {
	w.listenersLock.RLock()
	defer w.listenersLock.RUnlock()

	for l := range w.listeners {
		l.OnWorldUpdate(v)
	}
}
