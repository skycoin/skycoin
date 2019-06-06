package base

import (
	"fmt"
)

func opBoolPrint(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadBool(fp, inp1))
}

func opBoolEqual(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) == ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolUnequal(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) != ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolNot(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromBool(!ReadBool(fp, inp1))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolAnd(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) && ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opBoolOr(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadBool(fp, inp1) || ReadBool(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
