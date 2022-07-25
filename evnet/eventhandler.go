package evnet

const (
	EventHandled = 1
	EventNotHook = 2 //event is not hook on this handler
)

type EventHandler interface {
	OnEvent(evt *Event, ret *EventRet) int
	HandlerName() string
}

type CommonHandler struct {
	handlerName string
	handler     func(evt *Event, ret *EventRet) int
	hookEvents  map[string]bool //only add hook event when generate don't add hook event otherwise will need sync.Mutex
}

func (ch *CommonHandler) OnEvent(evt *Event, ret *EventRet) int {

	_, exists := ch.hookEvents[evt.Name]
	if !exists {
		return EventNotHook
	}
	return ch.handler(evt, ret)
}

func (ch *CommonHandler) HandlerName() string {
	return ch.handlerName
}

func NewCommonHandler(name string, f func(evt *Event, ret *EventRet) int, evts ...string) *CommonHandler {
	ch := new(CommonHandler)
	ch.handlerName = name
	ch.handler = f
	hookEvent := make(map[string]bool)
	for _, evt := range evts {
		hookEvent[evt] = true
	}
	ch.hookEvents = hookEvent

	return ch
}
