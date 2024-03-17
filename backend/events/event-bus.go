package events

type Handler func(data any)

var handlers = make(map[string][]Handler)

func RegisterHandler(name string, handler Handler) {
	_, ok := handlers[name]
	if !ok {
		handlers[name] = make([]Handler, 0)
	}
	handlers[name] = append(handlers[name], handler)
}

func Send(name string, data any) {
	if _, ok := handlers[name]; !ok {
		return
	}
	for _, handler := range handlers[name] {
		go handler(data)
	}
}
