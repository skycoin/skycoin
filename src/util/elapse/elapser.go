package elapse

import (
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
)

type Elapser struct {
	name             *string
	startTime        time.Time
	elapsedThreshold time.Duration
	Done             chan bool
	logger           *logging.Logger
}

func NewElapser(elapsedThreshold time.Duration, logger *logging.Logger) *Elapser {
	elapser := &Elapser{
		elapsedThreshold: elapsedThreshold,
		Done:             make(chan bool, 100),
		logger:           logger,
	}
	return elapser
}

func (e *Elapser) CheckForDone() {
	select {
	case <-e.Done:
		e.Elapsed()
	default:
	}
}

func (e *Elapser) Register(name string) {
	e.CheckForDone()
	e.name = &name
	e.startTime = time.Now()
	e.Done <- true
}

func (e *Elapser) ShowCurrentTime(step string) {
	stopTime := time.Now()
	if e.name == nil {
		e.logger.Warning("no registered events for elapsing, but found Elapser.ShowCurrentTime calling")
		return
	}
	elapsed := stopTime.Sub(e.startTime)
	e.logger.Info("%s[%s] elapsed %s", *e.name, step, elapsed)

}

func (e *Elapser) Elapsed() {
	stopTime := time.Now()
	if e.name == nil {
		e.logger.Warning("no registered events for elapsing, but found Elapser.Elapsed calling")
		return
	}
	elapsed := stopTime.Sub(e.startTime)
	if elapsed >= e.elapsedThreshold {
		e.logger.Warning("%s elapsed %s", *e.name, elapsed)
	}
	e.name = nil
}
