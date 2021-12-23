package world

import (
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
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

const (
	updatesBufferSize = 10
)

type SubscriberID int32

type World struct {
	listeners     map[SubscriberID]chan<- interface{}
	listenersLock sync.RWMutex
}

func NewWorld() *World {
	return &World{
		listeners: make(map[SubscriberID]chan<- interface{}),
	}
}

// Subscribe creates a new subscription for the specified SubscribedID. It returns a channel that will receive the
// updates. The SubscriberID provided must be unique, it is not allowed to call this function multiple times with the
// same SubscriberID unless it has been unregistered using Unsubscribe first.
func (w *World) Subscribe(id SubscriberID) <-chan interface{} {
	w.listenersLock.Lock()
	defer w.listenersLock.Unlock()

	ch := make(chan interface{}, updatesBufferSize)
	w.listeners[id] = ch

	return ch
}

// Unsubscribe cancels a previous subscription if it exists, otherwise it does nothing. The channel associated with the
// subscription will be closed.
func (w *World) Unsubscribe(id SubscriberID) {
	w.listenersLock.Lock()
	defer w.listenersLock.Unlock()

	ch, ok := w.listeners[id]
	if ok {
		close(ch)
		delete(w.listeners, id)
	}
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
