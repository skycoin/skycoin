package daemon

type strandReq struct {
	Desc string
	Func func()
}

// strand linearizes concurrent method calls through a single channel,
// to avoid concurrency issues when conflicting methods are called from
// multiple goroutines.
// Methods passed to strand() will block until completed.
func strand(c chan strandReq, desc string, f func()) {
	done := make(chan struct{})
	c <- strandReq{
		Desc: desc,
		Func: func() {
			f()
			close(done)
		},
	}
	<-done
}
