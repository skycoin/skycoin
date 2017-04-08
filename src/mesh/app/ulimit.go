package app

import (
	"syscall"
)

func setLimit(limit uint64) {
	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		panic(err)
	}

	rlimit.Max, rlimit.Cur = limit, limit // ~ number of nodes * 2

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		panic(err)
	}
}
