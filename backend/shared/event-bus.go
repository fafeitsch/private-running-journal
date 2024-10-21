package shared

import (
	"context"
	"log"
	"reflect"
)

type Handler func(data ...any)

type TrackUpsertedEvent struct {
	*SaveTrack
}

type TrackDeletedEvent struct {
	Id string
}

var handlers = make(map[string][]Handler)

var Context context.Context

func RegisterHandler(name string, handler Handler) {
	_, ok := handlers[name]
	if !ok {
		handlers[name] = make([]Handler, 0)
	}
	handlers[name] = append(handlers[name], handler)
}

func Send(name string, data ...any) {
	if _, ok := handlers[name]; !ok {
		return
	}
	log.Printf("sending eving \"%s\" with params %v", name, data)
	for _, handler := range handlers[name] {
		handler(data...)
	}
}

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
