package elapse

import (
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
)

// Elapser provides operations time elapsing
type Elapser struct {
	name             *string
	startTime        time.Time
	elapsedThreshold time.Duration
	Done             chan bool
	logger           *logging.Logger
}

// NewElapser creates an instance of Elapse
func NewElapser(elapsedThreshold time.Duration, logger *logging.Logger) *Elapser {
	elapser := &Elapser{
		elapsedThreshold: elapsedThreshold,
		Done:             make(chan bool, 100),
		logger:           logger,
	}
	return elapser
}

// CheckForDone checks if previous elapsing cycle is ready for finish.
func (e *Elapser) CheckForDone() {
	select {
	case <-e.Done:
		e.Elapsed()
	default:
	}
}

// Register registers new elapsing cycle name.
func (e *Elapser) Register(name string) {
	e.CheckForDone()
	e.name = &name
	e.startTime = time.Now()
	e.Done <- true
}

// ShowCurrentTime shows time from elapsing start to the current time.
func (e *Elapser) ShowCurrentTime(step string) {
	stopTime := time.Now()
	if e.name == nil {
		e.logger.Warning("no registered events for elapsing, but found Elapser.ShowCurrentTime calling")
		return
	}
	elapsed := stopTime.Sub(e.startTime)
	e.logger.Info("%s[%s] elapsed %s", *e.name, step, elapsed)

}

// Elapsed sets current elapsing cycle is finished and shows elapsed time if the time limit is reached.
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
