package evnet

import (
	"fmt"
	"testing"
)

func TestEventDispatcher_OnEvent(t *testing.T) {
	tests := []struct {
		name     string
		throttle int
	}{
		{
			name:     "test1",
			throttle: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evtdispatcher := NewAsyncEventDispatcher(tt.throttle, 10000)

			handler := NewCommonHandler("rpcmsg", func(evt *Event, evtret *EventRet) int {
				fmt.Println("rpcmsg")
				return EventHandled
			}, "rpcmsg")

			handler2 := NewCommonHandler("rpcmsg2", func(evt *Event, evtret *EventRet) int {
				fmt.Println("rpcmsg2")
				return EventHandled
			}, "rpcmsg2")

			evtdispatcher.PushHandler(handler)
			evtdispatcher.PushHandler(handler2)

			args := make(map[string]interface{})
			args["test"] = "aaa"
			testEvent := NewEvent("rpcmsg", args)

			args2 := make(map[string]interface{})
			args2["test2"] = "bbb"
			testEvent2 := NewEvent("rpcmsg2", args2)

			for i := 0; i < 10000; i++ {
				evtdispatcher.OnAsyncEvent(testEvent, nil)
				evtdispatcher.OnAsyncEvent(testEvent2, nil)
			}
		})
	}
}
