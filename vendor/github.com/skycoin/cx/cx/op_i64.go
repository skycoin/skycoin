package base

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opI64I64(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_STR:
		WriteObject(out1Offset, encoder.Serialize(strconv.Itoa(int(ReadI64(fp, inp1)))))
	case TYPE_BYTE:
		WriteMemory(out1Offset, FromByte(byte(ReadI64(fp, inp1))))
	case TYPE_I32:
		WriteMemory(out1Offset, FromI32(int32(ReadI64(fp, inp1))))
	case TYPE_I64:
		WriteMemory(out1Offset, FromI64(ReadI64(fp, inp1)))
	case TYPE_F32:
		WriteMemory(out1Offset, FromF32(float32(ReadI64(fp, inp1))))
	case TYPE_F64:
		WriteMemory(out1Offset, FromF64(float64(ReadI64(fp, inp1))))
	}
}

// The built-in print function formats its arguments in an
// implementation-specific

func opI64Print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadI64(fp, inp1))
}

// The built-in add function returns the sum of the two operands.

func opI64Add(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) + ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sub function returns the difference between the two operands.

func opI64Sub(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	var outB1 []byte
	if len(expr.Inputs) == 2 {
		inp2 := expr.Inputs[1]
		outB1 = FromI64(ReadI64(fp, inp1) - ReadI64(fp, inp2))
	} else {
		outB1 = FromI64(-ReadI64(fp, inp1))
	}
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in mul function returns the product of the two operands.

func opI64Mul(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) * ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in div function returns the quotient of the two operands.

func opI64Div(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) / ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in abs function returns the absolute value of the operand.

func opI64Abs(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Abs(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in pow function returns x**n for n>0 otherwise 1

func opI64Pow(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(math.Pow(float64(ReadI64(fp, inp1)), float64(ReadI64(fp, inp2)))))

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gt function returns true if operand 1 is greater than operand 2.

func opI64Gt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) > ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in gteq function returns true if operand 1 is greater than or
// equal to operand 2.

func opI64Gteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) >= ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lt function returns true if operand 1 is less than oeprand 2.

func opI64Lt(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) < ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in lteq function returns true if operand 1 is less than or
// equal to operand 2.

func opI64Lteq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) <= ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in eq function returns true if operand 1 is equal to operand 2.

func opI64Eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) == ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in uneq function returns true if operand 1 is different from operand 2.

func opI64Uneq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadI64(fp, inp1) != ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Mod(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) % ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Rand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	minimum := ReadI64(fp, inp1)
	maximum := ReadI64(fp, inp2)

	outB1 := FromI64(int64(rand.Intn(int(maximum-minimum)) + int(minimum)))

	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Bitand(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) & ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Bitor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) | ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Bitxor(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) ^ ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Bitclear(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(ReadI64(fp, inp1) &^ ReadI64(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Bitshl(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(uint64(ReadI64(fp, inp1)) << uint64(ReadI64(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func opI64Bitshr(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(uint64(ReadI64(fp, inp1)) >> uint64(ReadI64(fp, inp2))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in sqrt function returns the square root of the operand.

func opI64Sqrt(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Sqrt(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log function returns the natural logarithm of the operand.

func opI64Log(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Log(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log2 function returns the 2-logarithm of the operand.

func opI64Log2(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Log2(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in log10 function returns the 10-logarithm of the operand.

func opI64Log10(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	outB1 := FromI64(int64(math.Log10(float64(ReadI64(fp, inp1)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in max function returns the greatest value of the two operands.

func opI64Max(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(math.Max(float64(ReadI64(fp, inp1)), float64(ReadI64(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

// The built-in min function returns the smallest value of the two operands.

func opI64Min(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromI64(int64(math.Min(float64(ReadI64(fp, inp1)), float64(ReadI64(fp, inp2)))))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}
