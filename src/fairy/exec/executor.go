package exec

import (
	"fairy/util"
	"sync"
)

var gExecutor *Executor

func GetExecutor() *Executor {
	util.Once(gExecutor, func() {
		gExecutor = NewExecutor()
	})

	return gExecutor
}

func NewExecutor() *Executor {
	e := &Executor{}
	e.mutex = &sync.Mutex{}
	e.waitGroup = &sync.WaitGroup{}
	return e
}

/*
一:场景需求:
1:单独主业务线程
2:主业务线程+多个子模块线程，消息分发到特定的线程上
二:线程池
启动一个goroutine，回调回来时要求在主线程执行一个回调
*/
const EVENT_QUEUE_MAIN = 0

type Executor struct {
	stop      chan bool
	workQueue []*EventQueue
	workCount int
	mutex     *sync.Mutex
	waitGroup *sync.WaitGroup
}

func (self *Executor) OnExit() {
	self.Stop()
}

func (self *Executor) Stop() {
	self.mutex.Lock()
	for _, wq := range self.workQueue {
		wq.Stop()
	}
	self.mutex.Unlock()

	self.waitGroup.Wait()
}

func (self *Executor) Wait() {
	self.waitGroup.Wait()
}

func (self *Executor) Dispatch(ev Event) {
	self.DispatchEx(ev, EVENT_QUEUE_MAIN)
}

func (self *Executor) DispatchEx(ev Event, queueId int) {
	queue := self.GetQueue(queueId)
	queue.Push(ev)
}

func (self *Executor) Go(goCB Callback, mainCB Callback) {
	if mainCB != nil {
		go func() {
			goCB()
			self.Dispatch(NewFuncEvent(mainCB))
		}()
	} else {
		go goCB()
	}
}

func (self *Executor) GetQueue(queueId int) *EventQueue {
	if queueId < self.workCount {
		return self.workQueue[queueId]
	}

	self.mutex.Lock()
	defer self.mutex.Unlock()

	if queueId < len(self.workQueue) {
		return self.workQueue[queueId]
	}

	if len(self.workQueue) == 0 {
		util.RegisterExit(self)
	}

	count := queueId - len(self.workQueue) + 1
	for i := 0; i < count; i++ {
		queue := NewEventQueue()
		queue.Start(self.waitGroup)
		self.workQueue = append(self.workQueue, queue)
	}

	self.workCount = len(self.workQueue)

	return self.workQueue[queueId]
}
