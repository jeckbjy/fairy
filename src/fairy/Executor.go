package fairy

import (
	"fairy/timer"
	"sync"
)

var gExecutor *Executor

func GetExecutor() *Executor {
	if gExecutor == nil {
		gExecutor = NewExecutor()
	}

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
二:多线程
启动一个多线程，回调回来时要求在主线程执行一个回调
三:定时器触发时机
*/
const EVENT_QUEUE_MAIN = 0

type Executor struct {
	stop        chan bool
	workQueue   []*EventQueue
	workCount   int
	mutex       *sync.Mutex
	waitGroup   *sync.WaitGroup
	timerEngine *timer.TimerEngine
}

func (self *Executor) Stop() {
	self.Wait()
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
	go func() {
		goCB()
		self.Dispatch(NewFuncEvent(mainCB))
	}()
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

	count := queueId - len(self.workQueue) + 1
	for i := 0; i < count; i++ {
		queue := NewEventQueue()
		queue.Start(self.waitGroup)
		self.workQueue = append(self.workQueue, queue)
	}

	self.workCount = len(self.workQueue)

	return self.workQueue[queueId]
}
