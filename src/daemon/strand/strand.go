package strand

import (
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	// logDurationThreshold is how long to wait before reporting a function call's time
	logDurationThreshold = time.Millisecond * 100
	// writeWait is how long to wait to write to a request channel before logging the delay
	logQueueRequestWaitThreshold = time.Second * 3
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
// Methods passed to strand() will block until completed.
func Strand(logger *logging.Logger, c chan Request, name string, f func() error) error {
	quit := make(chan struct{})
	return WithQuit(logger, c, name, f, quit, nil)
}

// WithQuit linearizes concurrent method calls through a single channel,
// to avoid concurrency issues when conflicting methods are called from
// multiple goroutines.
// Methods passed to WithQuit() will block until completed.
// WithQuit accepts a quit channel and will return quitErr if the quit
// channel closes.
func WithQuit(logger *logging.Logger, c chan Request, name string, f func() error, quit chan struct{}, quitErr error) error {
	if Debug {
		logger.Debug("Strand precall %s", name)
	}

	done := make(chan struct{})
	var err error

	req := Request{
		Name: name,
		Func: func() error {
			defer close(done)

			// TODO: record time statistics in a data structure and expose stats via an API
			// logger.Debug("%s begin", name)

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
						logger.Warning("%s is taking longer than %s", name, threshold)
						threshold *= 10
						t.Reset(threshold)
					}
					t1 := time.Now()
					logger.Info("ELAPSED: %s", t1.Sub(t0))
				}
			}()

			if Debug {
				logger.Debug("Stranding %s", name)
			}

			err = f()

			// Log the error here so that the Request channel consumer doesn't need to
			if err != nil {
				logger.Error("%s error: %v", name, err)
			}

			// Notify us if the function call took too long
			elapsed := time.Now().Sub(t)
			if elapsed > logDurationThreshold {
				logger.Warning("%s took %s", name, elapsed)
			} else {
				// logger.Debug("%s took %s", name, elapsed)
			}

			return err
		},
	}

	// Log a message if waiting too long to write due to a full queue
loop:
	for {
		select {
		case <-quit:
			return nil
		case c <- req:
			break loop
		case <-time.After(logQueueRequestWaitThreshold):
			logger.Warning("Waited %s while trying to write %s to the strand request channel", logQueueRequestWaitThreshold, req.Name)
		}
	}

	<-done
	return err
}
