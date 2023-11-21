package evnet

import (
	"sort"
)

//EventHandlerInfo define event handler and the priority
type EventHandlerInfo struct {
	Handler  EventHandler
	Priority int
}

//EventDispatcher struct
type EventDispatcher struct {
	innerHandlers []*EventHandlerInfo
	asyncEvents   chan *AsyncEventInfo //channel AsyncEventInfo
	throttle      int                  //go routine size
}

type AsyncEventInfo struct {
	Evt      *Event
	Callback CallbackFun
}

//NewAsyncEventDispatcher receive AsyncEventInfo and send to right handler
func NewAsyncEventDispatcher(throttle int, chSize int) *EventDispatcher {
	evg := new(EventDispatcher)
	evg.asyncEvents = make(chan *AsyncEventInfo, chSize)
	evg.throttle = throttle

	for i := 0; i < throttle; i++ {
		//go routine handle asyncEvents
		go func() {
			for {
				asinfo := <-evg.asyncEvents
				if asinfo != nil {
					ret := evg.OnEvent(asinfo.Evt)

					if asinfo.Callback != nil {
						callbackf := func() {
							asinfo.Callback(asinfo.Evt, ret)
						}

						error.SafeCall(callbackf)
					}
				}
			}
		}()
	}

	return evg
}

func (evg *EventDispatcher) Len() int {
	return len(evg.innerHandlers)
}

func (evg *EventDispatcher) Less(i, j int) bool {
	ih := evg.innerHandlers[i]
	jh := evg.innerHandlers[j]

	return ih.Priority > jh.Priority
}

func (evg *EventDispatcher) Swap(i, j int) {
	evg.innerHandlers[i], evg.innerHandlers[j] = evg.innerHandlers[j], evg.innerHandlers[i]
}

//AddHandler and sort
func (evg *EventDispatcher) AddHandler(h EventHandler, sortIndex int) *EventDispatcher {
	if h == nil {
		panic("handler can't be nil")
	}

	hi := &EventHandlerInfo{
		Handler:  h,
		Priority: sortIndex,
	}
	evg.innerHandlers = append(evg.innerHandlers, hi)

	//排序
	sort.Sort(evg)

	return evg
}

//PushHandler add handler with Priority
func (evg *EventDispatcher) PushHandler(h EventHandler) *EventDispatcher {
	return evg.AddHandler(h, 0)
}

func (evg *EventDispatcher) RangeSort(f func(int, EventHandler) bool) {
	for _, v := range evg.innerHandlers {
		ret := f(v.Priority, v.Handler)

		if !ret {
			return
		}
	}
}

func (evg *EventDispatcher) Range(f func(EventHandler) bool) {
	for _, v := range evg.innerHandlers {
		ret := f(v.Handler)

		if !ret {
			return
		}
	}
}

//OnEvent iteration innerHandler and use the correct handler to handler event
func (evg *EventDispatcher) OnEvent(evt *Event) (ret *EventRet) {
	//ret = NewEventRet()

	defer func() {
		if err := recover(); err != nil {
			error.GetStackInfo()
			ret.Error = err
		}
	}()

	for _, h := range evg.innerHandlers {
		handleRet := h.Handler.OnEvent(evt, ret)

		switch handleRet {
		case EventNotHook:
			continue
		case EventHandled:
			return
		}
	}

	return
}

//OnAsyncEvent packet normal to AsyncEventInfo and send it to channel
func (evg *EventDispatcher) OnAsyncEvent(evt *Event, callback CallbackFun) {
	asinfo := new(AsyncEventInfo)
	asinfo.Evt = evt
	asinfo.Callback = callback

	evg.asyncEvents <- asinfo
}
