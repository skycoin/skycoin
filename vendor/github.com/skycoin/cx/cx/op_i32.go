package base

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opI32I32(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_STR:
		WriteObject(out1Offset, encoder.Serialize(strconv.Itoa(int(ReadI32(fp, inp1)))))
	case TYPE_BYTE:
		WriteMemory(out1Offset, FromByte(byte(ReadI32(fp, inp1))))
	case TYPE_I32:
		WriteMemory(out1Offset, FromI32(ReadI32(fp, inp1)))
	case TYPE_I64:
		WriteMemory(out1Offset, FromI64(int64(ReadI32(fp, inp1))))
	case TYPE_F32:
		WriteMemory(out1Offset, FromF32(float32(ReadI32(fp, inp1))))
	case TYPE_F64:
		WriteMemory(out1Offset, FromF64(float64(ReadI32(fp, inp1))))
	}
}

func opI32Print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI32(fp, inp1))
}

// The built-in add function returns the sum of two i32 numbers
func opI32Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) + ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sub function returns the difference of two i32 numbers
func opI32Sub(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	var outB1 []byte
	if len(expr.Inputs) == 2 {
		inp2 := expr.Inputs[1]
		outB1 = FromI32(ReadI32(fp, inp1) - ReadI32(fp, inp2))
	} else {
		outB1 = FromI32(-ReadI32(fp, inp1))
	}
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in mul function returns the product of two i32 numbers
func opI32Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) * ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in div function returns the quotient of two i32 numbers
func opI32Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) / ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in abs function returns the absolute number of the number
func opI32Abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Abs(float64(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The div built-in function returns x**n for n>0 otherwise 1
func opI32Pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(math.Pow(float64(ReadI32(fp, inp1)), float64(ReadI32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.
func opI32Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) > ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.
func opI32Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) >= ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lt function returns true if operand 1 is less than operand 2.
func opI32Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) < ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lteq function returns true if operand 1 is less than or equal
// to operand 1.
func opI32Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) <= ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.
func opI32Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) == ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.
func opI32Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI32(fp, inp1) != ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Mod(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) % ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Rand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	minimum := ReadI32(fp, inp1)
	maximum := ReadI32(fp, inp2)

	outB1 := FromI32(int32(rand.Intn(int(maximum-minimum)) + int(minimum)))

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Bitand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) & ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Bitor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) | ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Bitxor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) ^ ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Bitclear(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(ReadI32(fp, inp1) &^ ReadI32(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Bitshl(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(uint32(ReadI32(fp, inp1)) << uint32(ReadI32(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI32Bitshr(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(uint32(ReadI32(fp, inp1)) >> uint32(ReadI32(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sqrt function returns the square root of the operand.
func opI32Sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Sqrt(float64(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log function returns the natural logarithm of the operand.
func opI32Log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Log(float64(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log2 function returns the 2-logarithm of the operand.
func opI32Log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Log2(float64(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log10 function returns the 10-logarithm of the operand
func opI32Log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI32(int32(math.Log10(float64(ReadI32(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in max function returns the biggest of the two operands.
func opI32Max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(math.Max(float64(ReadI32(fp, inp1)), float64(ReadI32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in min function returns the smallest of the two operands.
func opI32Min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI32(int32(math.Min(float64(ReadI32(fp, inp1)), float64(ReadI32(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
