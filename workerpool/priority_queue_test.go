package workerpool

import (
	"testing"
)

func TestPriorityQueueImp_EnQueue(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []*TaskParam
		wantTasks []string
	}{
		{
			name: "test1",
			tasks: []*TaskParam{
				&TaskParam{
					TaskName:     "a",
					TaskPriority: 1,
				},
				&TaskParam{
					TaskName:     "b",
					TaskPriority: 5,
				},
				&TaskParam{
					TaskName:     "c",
					TaskPriority: 3,
				},
				&TaskParam{
					TaskName:     "d",
					TaskPriority: 2,
				},
				&TaskParam{
					TaskName:     "f",
					TaskPriority: 99999,
				},
			},
			wantTasks: []string{"f", "b", "c", "d", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewPriorityQueue()

			for i := 0; i < len(tt.tasks); i++ {
				pq.EnQueue(tt.tasks[i])
			}

			if pq.Len() == len(tt.wantTasks) {
				for i := 0; i < len(tt.wantTasks); i++ {
					got := pq.DeQueue()
					if got.TaskName != tt.wantTasks[i] {
						t.Errorf("want task = %v,got task = %v", tt.wantTasks[i], got.TaskName)
					}
				}
			} else {
				t.Errorf("want len = %v,got len = %v", len(tt.wantTasks), pq.Len())
			}

		})
	}
}
