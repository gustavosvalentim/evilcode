package api

type Mediator struct {
	handlers map[string][]func(interface{})
}

func (m *Mediator) AddEventHandler(e string, handler func(interface{})) {
	if _, ok := m.handlers[e]; !ok {
		m.handlers[e] = make([]func(interface{}), 0)
	}
	m.handlers[e] = append(m.handlers[e], handler)
}

func (m *Mediator) Dispatch(e string, data interface{}) {
	eventHandlers, ok := m.handlers[e]
	if !ok {
		return
	}
	for _, h := range eventHandlers {
		h(data)
	}
}
