package base

import "fmt"

func opI8Print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI8(fp, inp1))
}

func opI8Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) + ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Sub(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) - ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) * ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) / ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI8(fp, inp1) > ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI8(fp, inp1) >= ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI8(fp, inp1) < ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI8(fp, inp1) <= ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI8(fp, inp1) == ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI8(fp, inp1) != ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Bitand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) & ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Bitor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) | ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Bitxor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) ^ ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI8Bitclear(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI8(ReadI8(fp, inp1) &^ ReadI8(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
