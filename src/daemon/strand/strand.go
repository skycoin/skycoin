package strand

import (
	"log"
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
func Strand(logger *logging.Logger, c chan Request, name string, f func() error) error {
	done := make(chan struct{})
	var err error

	req := Request{
		Name: name,
		Func: func() error {
			defer close(done)

			// TODO: record time statistics in a data structure and expose stats via an API
			// logger.Debug("%s begin", name)

			t := time.Now()
			// minThreshold is how long to wait before reporting a function call's time
			minThreshold := time.Millisecond * 10

			// Log function duration at an exponential time interval,
			// this will notify us of any long running functions to look at.
			go func() {
				threshold := minThreshold
				t := time.NewTimer(threshold)
				defer t.Stop()

				for {
					t0 := time.Now()
					select {
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

			err = f()

			// Log the error here so that the Request channel consumer doesn't need to
			if err != nil {
				logger.Error("%s error: %v", name, err)
			}

			// Notify us if the function call took too long
			elapsed := time.Now().Sub(t)
			if elapsed > minThreshold {
				logger.Warning("%s took %s", name, elapsed)
			} else {
				// logger.Debug("%s took %s", name, elapsed)
			}

			return err
		},
	}

	// Log a message if waiting too long to write due to a full queue
	writeWait := time.Second * 3
	select {
	case c <- req:
	case <-time.After(writeWait):
		log.Println("Waited %s while trying to write %s to the strand request channel", writeWait, req.Name)
		c <- req
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
