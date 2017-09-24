package strand

import (
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
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
func Strand(logger *logging.Logger, c chan Request, req Request) error {
	done := make(chan struct{})
	var err error

	c <- Request{
		Name: req.Name,
		Func: func() error {
			defer close(done)

			// TODO: record time statistics in a data structure and expose stats via an API

			// logger.Debug("%s begin", req.Name)

			t := time.Now()

			wait := make(chan struct{})

			minThreshold := time.Millisecond * 10

			go func() {
				select {
				case <-wait:
				case <-time.After(minThreshold):
					logger.Warning("%s is taking longer than %s", req.Name, minThreshold)
				case <-time.After(minThreshold * 10):
					logger.Warning("%s is taking longer than %s", req.Name, minThreshold*10)
				case <-time.After(minThreshold * 100):
					logger.Warning("%s is taking longer than %s", req.Name, minThreshold*100)
				case <-time.After(minThreshold * 1000):
					logger.Warning("%s is taking longer than %s", req.Name, minThreshold*1000)
				}
			}()

			err = req.Func()
			if err != nil {
				logger.Error("%s error: %v", req.Name, err)
			}

			close(wait)

			elapsed := time.Now().Sub(t)
			if elapsed > minThreshold {
				logger.Warning("%s took %s", req.Name, elapsed)
			} else {
				// logger.Debug("%s took %s", req.Name, elapsed)
			}

			return err
		},
	}

	<-done
	return err
}

// StrandCanQuit linearizes concurrent method calls through a single channel,
// to avoid concurrency issues when conflicting methods are called from
// multiple goroutines.
// Methods passed to StrandCanQuit() will block until completed.
// StrandCanQuit accepts a quit channel and will return quitErr if the quit
// channel closes.
func StrandCanQuit(logger *logging.Logger, c chan Request, req Request, q chan struct{}, quitErr error) error {
	done := make(chan struct{})
	var err error

	select {
	case <-quit:
		return quitErr
	case c <- Request{
		Name: req.Name,
		Func: func() error {
			defer close(done)

			t := time.Now()

			logger.Debug("%s begin", req.Name)

			err = req.Func()
			if err != nil {
				logger.Error("%s error: %v", req.Name, err)
			}

			elapsed := time.Now().Sub(t)
			if elapsed > time.Second {
				logger.Warning("%s took %s", req.Name, elapsed)
			} else {
				logger.Debug("%s took %s", req.Name, elapsed)
			}

			return err
		},
	}:
	}

	<-done
	return err
}
