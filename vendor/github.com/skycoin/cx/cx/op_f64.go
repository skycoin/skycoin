package base

import (
	"fmt"
	"math"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opF64F64(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_STR:
		WriteObject(out1Offset, encoder.Serialize(strconv.FormatFloat(ReadF64(fp, inp1), 'f', -1, 64)))
	case TYPE_BYTE:
		WriteMemory(out1Offset, FromByte(byte(ReadF64(fp, inp1))))
	case TYPE_I32:
		WriteMemory(out1Offset, FromI32(int32(ReadF64(fp, inp1))))
	case TYPE_I64:
		WriteMemory(out1Offset, FromI64(int64(ReadF64(fp, inp1))))
	case TYPE_F32:
		WriteMemory(out1Offset, FromF32(float32(ReadF64(fp, inp1))))
	case TYPE_F64:
		WriteMemory(out1Offset, FromF64(ReadF64(fp, inp1)))
	}
}

// The print built-in function formats its arguments and prints them.
//
func opF64Print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadF64(fp, inp1))
}

// The built-in add function returns the sum of the two operands.
//
func opF64Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) + ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sub function returns the difference between the two operands.
//
func opF64Sub(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	var outB1 []byte
	if len(expr.Inputs) == 2 {
		inp2 := expr.Inputs[1]
		outB1 = FromF64(ReadF64(fp, inp1) - ReadF64(fp, inp2))
	} else {
		outB1 = FromF64(-ReadF64(fp, inp1))
	}
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in mul function returns the product of the two operands.
//
func opF64Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) * ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in div function returns the quotient between the two operands.
//
func opF64Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(ReadF64(fp, inp1) / ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in abs function returns the absolute value of the operand.
//
func opF64Abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Abs(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in pow function returns x**n for n>0 otherwise 1
//
func opF64Pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Pow(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gt function returns true if operand 1 is larger than operand 2.
//
func opF64Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) > ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
//
func opF64Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) >= ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
//
func opF64Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) < ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 2.
//
func opF64Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) <= ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
//
func opF64Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) == ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
//
func opF64Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF64(fp, inp1) != ReadF64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in cos function returns the cosine of the operand.
//
func opF64Cos(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Cos(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sin function returns the sine of the operand.
//
func opF64Sin(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Sin(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sqrt function returns the square root of the operand.
//
func opF64Sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Sqrt(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log function returns the natural logarithm of the operand.
//
func opF64Log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log2 function returns the 2-logarithm of the operand.
//
func opF64Log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log2(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log10 function returns the 10-logarithm of the operand.
//
func opF64Log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF64(math.Log10(ReadF64(fp, inp1)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in max function returns the largest value of the two operands.
//
func opF64Max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Max(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in min function returns the smallest value of the two operands.
//
func opF64Min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF64(math.Min(ReadF64(fp, inp1), ReadF64(fp, inp2)))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
