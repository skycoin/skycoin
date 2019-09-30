/*
Package elapse provides time measuring instruments
*/
package elapse

import (
	"time"

	"github.com/SkycoinProject/skycoin/src/util/logging"
)

// Elapser measures time elapsed for an operation. It is not thread-safe, use a different elapser per thread.
type Elapser struct {
	name             *string
	startTime        time.Time
	elapsedThreshold time.Duration
	Done             chan bool
	logger           *logging.Logger
}

// NewElapser creates an Elapser
func NewElapser(elapsedThreshold time.Duration, logger *logging.Logger) *Elapser {
	elapser := &Elapser{
		elapsedThreshold: elapsedThreshold,
		Done:             make(chan bool, 100),
		logger:           logger,
	}
	return elapser
}

// CheckForDone checks if the elapser has triggered and records the elapsed time
func (e *Elapser) CheckForDone() {
	select {
	case <-e.Done:
		e.Elapsed()
	default:
	}
}

// Register begins an operation to measure
func (e *Elapser) Register(name string) {
	e.CheckForDone()
	e.name = &name
	e.startTime = time.Now()
	e.Done <- true
}

// ShowCurrentTime logs the elapsed time so far
func (e *Elapser) ShowCurrentTime(step string) {
	stopTime := time.Now()
	if e.name == nil {
		e.logger.Warning("no registered events for elapsing, but found Elapser.ShowCurrentTime calling")
		return
	}
	elapsed := stopTime.Sub(e.startTime)
	e.logger.Infof("%s[%s] elapsed %s", *e.name, step, elapsed)

}

// Elapsed stops measuring an operation and logs the elapsed time if it exceeds the configured threshold
func (e *Elapser) Elapsed() {
	stopTime := time.Now()
	if e.name == nil {
		e.logger.Warning("no registered events for elapsing, but found Elapser.Elapsed calling")
		return
	}
	elapsed := stopTime.Sub(e.startTime)
	if elapsed >= e.elapsedThreshold {
		e.logger.Warningf("%s elapsed %s", *e.name, elapsed)
	}
	e.name = nil
}
