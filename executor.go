package fairy

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/jeckbjy/fairy/exit"
)

// QueueMainID 主队列ID
const QueueMainID = 0

// Callback 回调函数
type Callback func()

var gExecutor = NewExecutor()

// GetExecutor 获取全局调度器
func GetExecutor() *Executor {
	return gExecutor
}

// NewExecutor 创建调度器
func NewExecutor() *Executor {
	e := &Executor{}
	exit.Add(e.Stop)
	return e
}

// 任务队列
type queue struct {
	tasks   *list.List
	cond    *sync.Cond
	mutex   sync.Mutex
	stopped bool
}

func (q *queue) init() {
	q.tasks = list.New()
	q.cond = sync.NewCond(&q.mutex)
	q.stopped = false
}

func (q *queue) wait(tasks *list.List) {
	q.mutex.Lock()
	for !q.stopped && q.tasks.Len() == 0 {
		q.cond.Wait()
	}
	*tasks = *q.tasks
	q.tasks.Init()
	q.mutex.Unlock()
}

func (q *queue) process(tasks *list.List) {
	for iter := tasks.Front(); iter != nil; iter = iter.Next() {
		cb := iter.Value.(Callback)
		cb()
	}
}

func (q *queue) push(task Callback) {
	q.mutex.Lock()
	q.tasks.PushBack(task)
	q.cond.Signal()
	q.mutex.Unlock()
}

func (q *queue) stop() {
	q.mutex.Lock()
	q.stopped = true
	q.cond.Signal()
	q.mutex.Unlock()
}

// Executor 任务调度器
// 一个Main,多个Work队列，main和work队列不会同时执行消息
// 如果不需要Main与Work互斥,则可以将其中一个work队列作为主循环
type Executor struct {
	queues  []*queue
	mutex   sync.Mutex
	wg      sync.WaitGroup
	rwmux   sync.RWMutex
	stopped bool
}

// EnsureQueue 确保队列存在
func (exec *Executor) EnsureQueue(queueId uint) {
	exec.mutex.Lock()
	exec.ensure(queueId)
	exec.mutex.Unlock()
}

func (exec *Executor) ensure(queueId uint) *queue {
	if int(queueId) < len(exec.queues) {
		return exec.queues[queueId]
	}

	// main
	if queueId == 0 {
		q := &queue{}
		exec.queues = append(exec.queues, q)
		go exec.mainLoop(q)
		return q
	}

	// work
	count := int(queueId) - len(exec.queues) + 1
	for i := 0; i < count; i++ {
		q := &queue{}
		exec.queues = append(exec.queues, q)
		go exec.workLoop(q)
	}

	return exec.queues[queueId]
}

// Dispatch 分发到相应队列执行任务
func (exec *Executor) Dispatch(queueId uint, task Callback) error {
	exec.mutex.Lock()
	defer exec.mutex.Unlock()
	if exec.stopped {
		return fmt.Errorf("exec has stopped")
	}

	q := exec.ensure(queueId)
	q.push(task)

	return nil
}

// Stop 关闭队列,等待所有任务结束
func (exec *Executor) Stop() {
	exec.mutex.Lock()
	defer exec.mutex.Unlock()
	if exec.stopped {
		return
	}
	exec.stopped = true
	for _, q := range exec.queues {
		q.stop()
	}

	exec.wg.Wait()
}

// 运行主队列
func (exec *Executor) mainLoop(q *queue) {
	exec.wg.Add(1)
	defer exec.wg.Done()
	q.init()
	for {
		tasks := list.List{}
		q.wait(&tasks)

		exec.rwmux.Lock()
		q.process(&tasks)
		exec.rwmux.Unlock()

		if q.stopped {
			break
		}
	}
}

// 运行
func (exec *Executor) workLoop(q *queue) {
	exec.wg.Add(1)
	defer exec.wg.Done()
	q.init()
	for {
		tasks := list.List{}
		q.wait(&tasks)

		exec.rwmux.RLock()
		q.process(&tasks)
		exec.rwmux.RUnlock()
		if exec.stopped {
			break
		}
	}
}
