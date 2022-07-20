package workerpool

type TaskMethod func(params []interface{}) interface{}

type TaskParam struct {
	TaskName     string
	TaskMethod   TaskMethod
	TaskParam    []interface{}
	TaskPriority int
}
