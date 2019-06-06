package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// EscapeAnalysis ...
func EscapeAnalysis(fp int, inpOffset, outOffset int, arg *CXArgument) {
	heapOffset := AllocateSeq(arg.TotalSize + OBJECT_HEADER_SIZE)

	byts := ReadMemory(inpOffset, arg)

	// creating a header for this object
	size := encoder.SerializeAtomic(int32(len(byts)))

	var header = make([]byte, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteMemory(outOffset, off)
}

func opIdentity(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := GetFinalOffset(fp, inp1)
	out1Offset := GetFinalOffset(fp, out1)

	var elt *CXArgument
	if len(out1.Fields) > 0 {
		elt = out1.Fields[len(out1.Fields)-1]
	} else {
		elt = out1
	}

	if elt.DoesEscape {
		EscapeAnalysis(fp, inp1Offset, out1Offset, inp1)
	} else {
		switch elt.PassBy {
		case PASSBY_VALUE:
			WriteMemory(out1Offset, ReadMemory(inp1Offset, inp1))
		case PASSBY_REFERENCE:
			WriteMemory(out1Offset, encoder.SerializeAtomic(int32(inp1Offset)))
		}
	}
}

func opJmp(expr *CXExpression, fp int, call *CXCall) {
	inp1 := expr.Inputs[0]
	var predicate bool

	if expr.Label != "" {
		// then it's a goto
		call.Line = call.Line + expr.ThenLines
	} else {
		inp1Offset := GetFinalOffset(fp, inp1)

		predicateB := PROGRAM.Memory[inp1Offset : inp1Offset+inp1.Size]
		encoder.DeserializeAtomic(predicateB, &predicate)

		if predicate {
			call.Line = call.Line + expr.ThenLines
		} else {
			call.Line = call.Line + expr.ElseLines
		}
	}
}
