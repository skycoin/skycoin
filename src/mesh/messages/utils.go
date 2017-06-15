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

func BytesToNum(bytes []uint8) int {
	n := len(bytes)
	result := 0
	multiplier := 1
	for i := 0; i < n; i++ {
		result += int(bytes[i]) * multiplier
		multiplier *= 256
	}
	return result
}

func NumToBytes(num, len int) []byte {
	result := make([]byte, len)
	multiplier := 1
	for i := 0; i < len; i++ {
		result[i] = (uint8)(num & (255 * multiplier) >> (uint32(i) * 8))
		multiplier *= 256
	}
	return result
}
