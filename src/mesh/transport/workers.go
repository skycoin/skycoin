package transport

import (
	"github.com/skycoin/skycoin/src/mesh/messages"
)

var (
	maxWorkers = 10
	maxQueue   = 1024
)

type Job struct {
	msg *messages.TransportDatagramTransfer
}

type Worker struct {
	transport  *Transport
	workerPool chan chan *Job
	jobChannel chan *Job
	quit       chan bool
}

func NewWorker(workerPool chan chan *Job, tr *Transport) *Worker {
	return &Worker{
		transport:  tr,
		workerPool: workerPool,
		jobChannel: make(chan *Job),
		quit:       make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel
			select {
			case job := <-w.jobChannel:
				w.transport.packetsSent++
				job.msg.Sequence = w.transport.packetsSent
				w.transport.sendPacket(job.msg)
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	transport  *Transport
	workerPool chan chan *Job
	maxWorkers int
}

func NewDispatcher(tr *Transport, maxWorkers int) *Dispatcher {
	pool := make(chan chan *Job, maxWorkers)
	d := &Dispatcher{tr, pool, maxWorkers}
	return d
}

func (dispatcher *Dispatcher) Run() {
	for i := 0; i < dispatcher.maxWorkers; i++ {
		worker := NewWorker(dispatcher.workerPool, dispatcher.transport)
		worker.Start()
	}
	go dispatcher.dispatch()
}

func (dispatcher *Dispatcher) dispatch() {
	for {
		select {
		case job := <-dispatcher.transport.pendingOut:
			go func(job *Job) {
				jobChannel := <-dispatcher.workerPool
				jobChannel <- job
			}(job)
		}
	}
}
