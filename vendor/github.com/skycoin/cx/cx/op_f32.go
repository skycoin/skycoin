package base

import (
	"fmt"
	"math"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opF32F32(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_STR:
		WriteObject(out1Offset, encoder.Serialize(strconv.FormatFloat(float64(ReadF32(fp, inp1)), 'f', -1, 32)))
	case TYPE_BYTE:
		WriteMemory(out1Offset, FromByte(byte(ReadF32(fp, inp1))))
	case TYPE_I32:
		WriteMemory(out1Offset, FromI32(int32(ReadF32(fp, inp1))))
	case TYPE_I64:
		WriteMemory(out1Offset, FromI64(int64(ReadF32(fp, inp1))))
	case TYPE_F32:
		WriteMemory(out1Offset, FromF32(float32(ReadF32(fp, inp1))))
	case TYPE_F64:
		WriteMemory(out1Offset, FromF64(float64(ReadF32(fp, inp1))))
	}
}

func opF32Isnan(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromBool(math.IsNaN(float64(ReadF32(fp, inp1))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opF32Print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadF32(fp, inp1))
}

// The built-in add function returns the sum of the two operands.
//
func opF32Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(ReadF32(fp, inp1) + ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sub function returns the difference between the two operands.
//
func opF32Sub(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	var outB1 []byte
	if len(expr.Inputs) == 2 {
		inp2 := expr.Inputs[1]
		outB1 = FromF32(ReadF32(fp, inp1) - ReadF32(fp, inp2))
	} else {
		outB1 = FromF32(-ReadF32(fp, inp1))
	}
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in mul function returns the product of the two operands.
//
func opF32Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(ReadF32(fp, inp1) * ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in div function returns the quotient between the two operands.
//
func opF32Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(ReadF32(fp, inp1) / ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in abs function returns the absolute value of the operand.
//
func opF32Abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Abs(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in pow function returns x**n for n>0 otherwise 1
//
func opF32Pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(float32(math.Pow(float64(ReadF32(fp, inp1)), float64(ReadF32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
//
func opF32Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) > ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gteq function returns true if the operand 1 is greater than or
// equal to operand 2.
//
func opF32Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) >= ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lt function returns true if operand 1 is less than operand 2.

func opF32Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) < ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lteq function returns true if operand 1 is less than or
// equal to operand 2.
//
func opF32Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) <= ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
//
func opF32Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) == ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in uneq function returns true operand1 is different from operand 2.
//
func opF32Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadF32(fp, inp1) != ReadF32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in cos function returns the cosine of the operand.
//
func opF32Cos(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Cos(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sin function returns the sine of the operand.
//
func opF32Sin(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Sin(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sqrt function returns the square root of the operand.
//
func opF32Sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Sqrt(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log function returns the natural logarithm of the operand.
//
func opF32Log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Log(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log2 function returns the 2-logarithm of the operand.
//
func opF32Log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Log2(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log10 function returns the 10-logarithm of the operand.
//
func opF32Log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromF32(float32(math.Log10(float64(ReadF32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in max function returns the largest value of the two operands.
//
func opF32Max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(float32(math.Max(float64(ReadF32(fp, inp1)), float64(ReadF32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in min function returns the smallest value of the two operands.
//
func opF32Min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromF32(float32(math.Min(float64(ReadF32(fp, inp1)), float64(ReadF32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
