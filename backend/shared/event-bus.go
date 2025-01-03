package shared

import (
	"context"
	"log"
	"reflect"
	"time"
)

type Handler func(data ...any)

type TrackUpsertedEvent struct {
	*SaveTrack
}

type TrackDeletedEvent struct {
	Id string
}

type JournalEntryUpsertedEvent struct {
	*JournalEntry
	OldTrackId string
	OldDate    *time.Time
}

type JournalEntryDeletedEvent struct {
	*JournalEntry
}

type GitEnablementChangedEvent struct {
	NewValue bool
}

type GitPushChangedEvent struct {
	NewValue bool
}

type SettingsChangedEvent struct{}

type TileServerChangedEvent struct {
	NewValue string
}

type TileServerCacheEnabledEvent struct {
	NewValue bool
}

type MigrationEvent struct {
	OldVersion int
	NewVersion int
}

var handlers = make(map[string][]Handler)

var Context context.Context

func SendEvent(data any) {
	namedType := reflect.TypeOf(data)
	if _, ok := handlers[namedType.String()]; !ok {
		return
	}
	for _, handler := range handlers[namedType.String()] {
		handler(data)
	}
}

func Listen[K any](event K, handler func(k K)) {
	namedType := reflect.TypeOf(event)
	_, ok := handlers[namedType.String()]
	if !ok {
		handlers[namedType.String()] = make([]Handler, 0)
	}
	handlers[namedType.String()] = append(
		handlers[namedType.String()], func(data ...any) {
			payload := data[0]
			payloadType := reflect.TypeOf(payload)
			if payloadType != namedType {
				log.Fatalf("handler for %s received parameter of type %s", namedType, payloadType.String())
			}
			handler(payload.(K))
		},
	)
}
