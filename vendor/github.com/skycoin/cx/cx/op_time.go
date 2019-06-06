// +build base extra full

package base

import (
	"time"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func op_time_UnixMilli(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI64(makeTimestamp())
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_time_UnixNano(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	outB1 := FromI64(time.Now().UnixNano())
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_time_Sleep(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	time.Sleep(time.Duration(ReadI32(fp, inp1)) * time.Millisecond)
}
