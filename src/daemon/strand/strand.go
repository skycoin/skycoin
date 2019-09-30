/*
Package strand is a utility for linearizing method calls, similar to locking.

The strand method is functionally similar to a lock, but operates on a queue
of method calls.
*/
package strand

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/SkycoinProject/skycoin/src/util/logging"
)

const (
	// logDurationThreshold is how long to wait before reporting a function call's time
	logDurationThreshold = time.Millisecond * 100
	// writeWait is how long to wait to write to a request channel before logging the delay
	logQueueRequestWaitThreshold = time.Second * 1
)

var (
	// Debug enables debug logging
	Debug = false
)

// Request is sent to the channel provided to Strand
type Request struct {
	Name string
	Func func() error
}

// Strand linearizes concurrent method calls through a single channel,
// to avoid concurrency issues when conflicting methods are called from
// multiple goroutines.
// Methods passed to Strand() will block until completed.
// Strand accepts a quit channel and will return quitErr if the quit
// channel closes.
func Strand(logger *logging.Logger, c chan Request, name string, f func() error, quit chan struct{}, quitErr error) error {
	if Debug {
		logger.WithField("operation", name).Debug("Strand precall")
	}

	done := make(chan struct{})
	var err error

	req := Request{
		Name: name,
		Func: func() error {
			defer close(done)

			// TODO: record time statistics in a data structure and expose stats via an API
			// logger.Debugf("%s begin", name)

			t := time.Now()

			// Log function duration at an exponential time interval,
			// this will notify us of any long running functions to look at.
			go func() {
				threshold := logDurationThreshold
				t := time.NewTimer(threshold)
				defer t.Stop()

				for {
					t0 := time.Now()
					select {
					case <-quit:
						return
					case <-done:
						return
					case <-t.C:
						logger.WithFields(logrus.Fields{
							"operation": name,
							"threshold": threshold,
						}).Warning("Strand operation exceeded threshold")
						threshold *= 10
						t.Reset(threshold)
					}
					t1 := time.Now()
					logger.WithField("elapsed", t1.Sub(t0)).Info()
				}
			}()

			if Debug {
				logger.WithField("operation", name).Debug("Stranding")
			}

			err = f()

			// Notify us if the function call took too long
			elapsed := time.Since(t)
			if elapsed > logDurationThreshold {
				logger.WithFields(logrus.Fields{
					"operation": name,
					"elapsed":   elapsed,
				}).Warning()
			} else if Debug {
				logger.WithFields(logrus.Fields{
					"operation": name,
					"elapsed":   elapsed,
				}).Debug()
			}

			return err
		},
	}

	// Log a message if waiting too long to write due to a full queue
	t := time.Now()
loop:
	for {
		select {
		case <-quit:
			return quitErr
		case c <- req:
			break loop
		case <-time.After(logQueueRequestWaitThreshold):
			logger.Warningf("Waited %s while trying to write %s to the strand request channel", time.Since(t), req.Name)
		}
	}

	t = time.Now()
	for {
		select {
		case <-quit:
			return quitErr
		case <-done:
			return err
		case <-time.After(logQueueRequestWaitThreshold):
			logger.Warningf("Waited %s while waiting for %s to be done or quit", time.Since(t), req.Name)
		}
	}
}
