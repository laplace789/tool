package workerpool

import (
	"context"
	"sync"
	"time"
)

const (
	PoolSize     = 32
	inputChannel = 10000
	jobChannel   = 10000
)

var WorkerPoolInstance *WorkerPool

type WorkerPool struct {
	inputChan chan *TaskParam
	jobsChan  chan *TaskParam

	pq PriorityQueue

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func Init() {
	WorkerPoolInstance = NewWorkerPool()
}

func (c *WorkerPool) listen() {
	defer close(c.jobsChan)
	for {
		select {
		case job, ok := <-c.inputChan:
			if c.ctx.Err() != nil && !ok {
				return
			}
			//Enqueue job into priority queue
			c.pq.EnQueue(job)
			//c.jobsChan <- job
		}
	}
}

func (c *WorkerPool) worker(num int) {
	defer c.wg.Done()
	for {
		select {
		case job, ok := <-c.jobsChan:
			if c.ctx.Err() != nil && !ok {
				return
			}
			job.TaskMethod(job.TaskParam)
		case <-c.ctx.Done():
			//still have tasks
			if c.pq.Len() > 0 {
				continue
			}
			return
		}
	}
}

func (c *WorkerPool) deliver() {
	for {
		//has job in priority queue
		if c.pq.Len() > 0 {
			//todo check available worker
			taskCount := c.pq.Len()
			for i := 0; i < taskCount; i++ {
				task := c.pq.DeQueue()
				c.jobsChan <- task
			}
		} else {
			//wait 500 ms
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func NewWorkerPool() *WorkerPool {
	pool := &WorkerPool{
		inputChan: make(chan *TaskParam, inputChannel),
		jobsChan:  make(chan *TaskParam, jobChannel),
		wg:        &sync.WaitGroup{},
		pq:        NewPriorityQueue(),
	}

	pool.ctx, pool.cancel = context.WithCancel(context.Background())

	for i := 0; i < PoolSize; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}
	go pool.listen()
	go pool.deliver()

	return pool
}

func (c *WorkerPool) AddTask(task TaskMethod, params ...interface{}) {
	c.inputChan <- &TaskParam{TaskMethod: task, TaskParam: params}
}

func (c *WorkerPool) Name() string {
	return "WorkerPool"
}

func (c *WorkerPool) ShutdownPriority() int {
	return 0
}

func (c *WorkerPool) BeforeShutdown() {
	c.cancel()
	c.wg.Wait()
	close(c.inputChan)
}
