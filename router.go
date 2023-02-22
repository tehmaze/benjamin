package benjamin

type Route struct {
	Peripheral Peripheral
	EventHandler
}

type Router map[EventType][]Route

func (r Router) Run(d Device) {
	for e := range d.Events() {
		r.Handle(e)
	}
}

func (r Router) Handle(e Event) {
	//log.Printf("router: %s has %d handlers", e.Type, len(r[e.Type]))
	for _, h := range r[e.Type] {
		if h.Peripheral == e.Peripheral {
			h.Handle(e)
		}
	}
}

func (r Router) On(p Peripheral, t EventType, h EventHandler) {
	r[t] = append(r[t], Route{
		Peripheral:   p,
		EventHandler: h,
	})
}
