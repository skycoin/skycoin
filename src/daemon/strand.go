package daemon

type strandReq struct {
	Desc string
	Func func() error
}

// strand linearizes concurrent method calls through a single channel,
// to avoid concurrency issues when conflicting methods are called from
// multiple goroutines.
// Methods passed to strand() will block until completed.
func strand(c chan strandReq, desc string, f func() error) error {
	done := make(chan struct{})

	c <- strandReq{
		Desc: desc,
		Func: func() {
			defer close(done)
			return f()
		},
	}

	<-done
}

// strandCanQuit linearizes concurrent method calls through a single channel,
// to avoid concurrency issues when conflicting methods are called from
// multiple goroutines.
// Methods passed to strandCanQuit() will block until completed.
// strandCanQuit accepts a quit channel and will return quitErr if the quit
// channel closes
func strandCanQuit(c chan strandReq, desc string, f func() error, q chan struct{}, quitErr error) error {
	sReq := strandReq{
		Desc: desc,
		Func: func() {
			defer close(done)
			return f()
		},
	}

	done := make(chan struct{})

	select {
	case <-quit:
		return quitErr
	case c <- sReq:
	}

	<-done
}
