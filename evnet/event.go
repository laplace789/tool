package evnet

type EventArgs map[string]interface{}

type CallbackFun func(evt *Event, ret *EventRet)

//Event define basic struct for an async event
type Event struct {
	Name string    //event name
	Args EventArgs //args in event
}

func NewEvent(name string, args EventArgs) *Event {
	evt := new(Event)
	evt.Name = name
	evt.Args = args

	return evt
}

type EventRet struct {
	Ret   map[string]interface{}
	Error interface{}
}
