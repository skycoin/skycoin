package messages

func Increase(throttle uint32) uint32 {
	if throttle == 0 {
		return 1
	} else {
		return throttle * 2
	}
}

func Decrease(throttle uint32) uint32 {
	if throttle <= 1 {
		return 0
	} else {
		return throttle / 2
	}
}
