package base

import (
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var onMessages = map[string]string{
	"arg-arg-input":    "Replace %s.Input.%d with %s",
	"arg-arg-output":   "Replace %s.Output.%d with %s",
	"arg-expr":         "Add ",
	"prgrm-arg-input":  "Print %s.Input.%d's value",
	"prgrm-arg-output": "Print %s.Output.%d's value",
}
var ofMessages = map[string]string{
	"arg-arg-input":    "Replace %[3]s with %[1]s.Input.%[2]d",
	"arg-arg-output":   "Replace %[3]s with %[1]s.Output.%[2]d",
	"strct-arg-input":  "Add %[1]s.Input.%[2]d as a new field of %s",
	"strct-arg-output": "Add %[1]s.Output.%[2]d as a new field of %s",
	"prgrm-arg-input":  "Print %[1]s.Input.%[2]d's value",
	"prgrm-arg-output": "Print %[1]s.Output.%[2]d's value",
}

// GetInferActions ...
func GetInferActions(inp *CXArgument, fp int) []string {
	inpOffset := GetFinalOffset(fp, inp)

	var off int32
	encoder.DeserializeAtomic(PROGRAM.Memory[inpOffset:inpOffset+TYPE_POINTER_SIZE], &off)

	var l int32
	_l := PROGRAM.Memory[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]
	encoder.DeserializeAtomic(_l[4:], &l)

	result := make([]string, l)

	// for c := int(l); c > 0; c-- {
	for c := 0; c < int(l); c++ {
		var elOff int32
		// encoder.DeserializeAtomic(PROGRAM.Memory[int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + (c - 1) * TYPE_POINTER_SIZE : int(off) + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + c * STR_HEADER_SIZE], &elOff)
		encoder.DeserializeAtomic(PROGRAM.Memory[int(off)+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+c*TYPE_POINTER_SIZE:int(off)+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+(c+1)*STR_HEADER_SIZE], &elOff)

		var size int32
		encoder.DeserializeAtomic(PROGRAM.Memory[elOff:elOff+STR_HEADER_SIZE], &size)

		var res string
		encoder.DeserializeRaw(PROGRAM.Memory[elOff:elOff+STR_HEADER_SIZE+size], &res)

		// result[int(l) - c] = res
		result[c] = res
	}

	return result
}

func opAffPrint(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(GetInferActions(inp1, fp))
	// for _, aff := range GetInferActions(inp1, fp) {
	// 	fmt.Println(aff)
	// }
}

// CallAffPredicate ...
func CallAffPredicate(fn *CXFunction, predValue []byte) byte {
	prevCall := &PROGRAM.CallStack[PROGRAM.CallCounter]

	PROGRAM.CallCounter++
	newCall := &PROGRAM.CallStack[PROGRAM.CallCounter]
	newCall.Operator = fn
	newCall.Line = 0
	newCall.FramePointer = PROGRAM.StackPointer
	PROGRAM.StackPointer += newCall.Operator.Size

	newFP := newCall.FramePointer

	// wiping next mem frame (removing garbage)
	for c := 0; c < fn.Size; c++ {
		PROGRAM.Memory[newFP+c] = 0
	}

	// sending value to predicate function
	WriteMemory(
		GetFinalOffset(newFP, newCall.Operator.Inputs[0]),
		predValue)

	prevCC := PROGRAM.CallCounter
	for {
		call := &PROGRAM.CallStack[PROGRAM.CallCounter]
		call.ccall(PROGRAM)
		if PROGRAM.CallCounter < prevCC {
			break
		}
	}

	prevCall.Line--

	return ReadMemory(GetFinalOffset(
		newCall.FramePointer,
		newCall.Operator.Outputs[0]),
		newCall.Operator.Outputs[0])[0]
}

// This might not make sense, as we can use normal programming to create conditions on values
// func QueryValue (fn *CXFunction, argOffsetB []byte, affOffset *int) {
// 	for c := 0; c <= PROGRAM.CallCounter; c++ {
// 		inFP := 0
// 		op := PROGRAM.CallStack[c].Operator

// 		for _, expr := range op.Expressions {
// 			if expr.Operator == nil {
// 				for _, out := range expr.Outputs {
// 					if fn.Inputs[0].Type == out.Type && out.Name != "" {
// 						res := CallAffPredicate(fn, PROGRAM.Memory[inFP + out.Offset : inFP + out.Offset + out.TotalSize])

// 						if res == 1 {
// 							*affOffset = WriteToSlice(*affOffset, argOffsetB)

// 							outNameB := encoder.Serialize(out.Name)
// 							outNameOffset := AllocateSeq(len(outNameB))
// 							WriteMemory(outNameOffset, outNameB)
// 							outNameOffsetB := encoder.SerializeAtomic(int32(outNameOffset))

// 							*affOffset = WriteToSlice(*affOffset, outNameOffsetB)
// 						}
// 					}
// 				}
// 			}
// 		}

// 		inFP += op.Size
// 	}
// }

// Used by QueryArgument to query inputs and then outputs from expressions.
func queryParam(fn *CXFunction, args []*CXArgument, exprLbl string, argOffsetB []byte, affOffset *int) {
	for i, arg := range args {

		// Name
		argNameB := encoder.Serialize(arg.Name)
		argNameOffsetB := encoder.SerializeAtomic(int32(WriteObjectRetOff(argNameB)))

		argOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE + STR_SIZE)

		var typOffset int
		elt := GetAssignmentElement(arg)
		if elt.CustomType != nil {
			// then it's custom type
			typOffset = WriteObjectRetOff(encoder.Serialize(elt.CustomType.Package.Name + "." + elt.CustomType.Name))
		} else {
			// then it's native type
			typOffset = WriteObjectRetOff(encoder.Serialize(TypeNames[elt.Type]))
		}

		// Name
		WriteMemory(argOffset+OBJECT_HEADER_SIZE, argNameOffsetB)
		// Index
		WriteMemory(argOffset+OBJECT_HEADER_SIZE+STR_SIZE, encoder.SerializeAtomic(int32(i)))
		// Type
		WriteMemory(argOffset+OBJECT_HEADER_SIZE+STR_SIZE+I32_SIZE, encoder.SerializeAtomic(int32(typOffset)))

		res := CallAffPredicate(fn, PROGRAM.Memory[argOffset+OBJECT_HEADER_SIZE:argOffset+OBJECT_HEADER_SIZE+STR_SIZE+I32_SIZE+STR_SIZE])

		if res == 1 {
			*affOffset = WriteToSlice(*affOffset, argOffsetB)

			affNameB := encoder.Serialize(fmt.Sprintf("%s.%d", exprLbl, i))
			affNameOffset := AllocateSeq(len(affNameB))
			WriteMemory(affNameOffset, affNameB)

			*affOffset = WriteToSlice(*affOffset, encoder.SerializeAtomic(int32(affNameOffset)))
		}
	}
}

// QueryArgument ...
func QueryArgument(fn *CXFunction, expr *CXExpression, argOffsetB []byte, affOffset *int) {
	for _, ex := range expr.Function.Expressions {
		if ex.Label == "" {
			// it's a non-labeled expression
			continue
		}

		queryParam(fn, ex.Inputs, ex.Label+".Input", argOffsetB, affOffset)
		queryParam(fn, ex.Outputs, ex.Label+".Output", argOffsetB, affOffset)
	}
}

// QueryExpressions ...
func QueryExpressions(fn *CXFunction, expr *CXExpression, exprOffsetB []byte, affOffset *int) {
	for _, ex := range expr.Function.Expressions {
		if ex.Operator == nil || ex.Label == "" {
			// then it's a variable declaration
			// or it's a non-labeled expression
			continue
		}

		var opNameB []byte
		if ex.Operator.IsNative {
			opNameB = encoder.Serialize(OpNames[ex.Operator.OpCode])
		} else {
			opNameB = encoder.Serialize(ex.Operator.Name)
		}

		opNameOffset := AllocateSeq(len(opNameB))
		WriteMemory(opNameOffset, opNameB)
		opNameOffsetB := encoder.SerializeAtomic(int32(opNameOffset))

		res := CallAffPredicate(fn, opNameOffsetB)

		if res == 1 {
			*affOffset = WriteToSlice(*affOffset, exprOffsetB)

			lblNameB := encoder.Serialize(ex.Label)
			lblNameOffset := AllocateSeq(len(lblNameB))
			WriteMemory(lblNameOffset, lblNameB)
			lblNameOffsetB := encoder.SerializeAtomic(int32(lblNameOffset))

			*affOffset = WriteToSlice(*affOffset, lblNameOffsetB)
		}
	}
}

func getSignatureSlice(params []*CXArgument) int {
	var sliceOffset int
	for _, param := range params {

		var typOffset int
		if param.CustomType != nil {
			// then it's custom type
			typOffset = WriteObjectRetOff(encoder.Serialize(param.CustomType.Package.Name + "." + param.CustomType.Name))
		} else {
			// then it's native type
			typOffset = WriteObjectRetOff(encoder.Serialize(TypeNames[param.Type]))
		}

		sliceOffset = WriteToSlice(sliceOffset, encoder.SerializeAtomic(int32(typOffset)))
	}

	return sliceOffset
}

// Helper function for QueryStructure. Used to query all the structs in a particular package
func queryStructsInPackage(fn *CXFunction, expr *CXExpression, strctOffsetB []byte, affOffset *int, pkg *CXPackage) {
	for _, f := range pkg.Structs {
		strctNameB := encoder.Serialize(f.Name)

		strctNameOffsetB := encoder.SerializeAtomic(int32(WriteObjectRetOff(strctNameB)))

		strctOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE)
		// Name
		WriteMemory(strctOffset+OBJECT_HEADER_SIZE, strctNameOffsetB)

		val := PROGRAM.Memory[strctOffset+OBJECT_HEADER_SIZE : strctOffset+OBJECT_HEADER_SIZE+STR_SIZE]
		res := CallAffPredicate(fn, val)

		if res == 1 {
			*affOffset = WriteToSlice(*affOffset, strctOffsetB)
			*affOffset = WriteToSlice(*affOffset, strctNameOffsetB)
		}
	}
}

// QueryStructure ...
func QueryStructure(fn *CXFunction, expr *CXExpression, strctOffsetB []byte, affOffset *int) {
	queryStructsInPackage(fn, expr, strctOffsetB, affOffset, expr.Package)
	for _, imp := range expr.Package.Imports {
		queryStructsInPackage(fn, expr, strctOffsetB, affOffset, imp)
	}
}

// QueryFunction ...
func QueryFunction(fn *CXFunction, expr *CXExpression, fnOffsetB []byte, affOffset *int) {
	for _, f := range expr.Package.Functions {
		if f.Name == SYS_INIT_FUNC {
			continue
		}

		var opNameB []byte
		if f.IsNative {
			opNameB = encoder.Serialize(OpNames[f.OpCode])
		} else {
			opNameB = encoder.Serialize(f.Name)
		}

		opNameOffsetB := encoder.SerializeAtomic(int32(WriteObjectRetOff(opNameB)))

		inpSigOffset := getSignatureSlice(f.Inputs)
		outSigOffset := getSignatureSlice(f.Outputs)

		fnOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + TYPE_POINTER_SIZE + TYPE_POINTER_SIZE)
		// Name
		WriteMemory(fnOffset+OBJECT_HEADER_SIZE, opNameOffsetB)
		// InputSignature
		WriteMemory(fnOffset+OBJECT_HEADER_SIZE+TYPE_POINTER_SIZE, encoder.SerializeAtomic(int32(inpSigOffset)))
		// OutputSignature
		WriteMemory(fnOffset+OBJECT_HEADER_SIZE+TYPE_POINTER_SIZE+TYPE_POINTER_SIZE, encoder.SerializeAtomic(int32(outSigOffset)))

		val := PROGRAM.Memory[fnOffset+OBJECT_HEADER_SIZE : fnOffset+OBJECT_HEADER_SIZE+STR_SIZE+TYPE_POINTER_SIZE+TYPE_POINTER_SIZE]
		res := CallAffPredicate(fn, val)

		if res == 1 {
			*affOffset = WriteToSlice(*affOffset, fnOffsetB)
			*affOffset = WriteToSlice(*affOffset, opNameOffsetB)
		}
	}
}

// QueryCaller ...
func QueryCaller(fn *CXFunction, expr *CXExpression, callerOffsetB []byte, affOffset *int) {
	if PROGRAM.CallCounter == 0 {
		// then it's entry point
		return
	}

	call := PROGRAM.CallStack[PROGRAM.CallCounter-1]

	var opNameB []byte
	if call.Operator.IsNative {
		opNameB = encoder.Serialize(OpNames[call.Operator.OpCode])
	} else {
		opNameB = encoder.Serialize(call.Operator.Package.Name + "." + call.Operator.Name)
	}

	opNameOffsetB := encoder.SerializeAtomic(int32(WriteObjectRetOff(opNameB)))

	callOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE)
	// FnName
	WriteMemory(callOffset+OBJECT_HEADER_SIZE, opNameOffsetB)
	// FnSize
	WriteMemory(callOffset+OBJECT_HEADER_SIZE+STR_SIZE, encoder.SerializeAtomic(int32(call.Operator.Size)))

	res := CallAffPredicate(fn, PROGRAM.Memory[callOffset+OBJECT_HEADER_SIZE:callOffset+OBJECT_HEADER_SIZE+STR_SIZE+I32_SIZE])

	if res == 1 {
		*affOffset = WriteToSlice(*affOffset, callerOffsetB)
	}
}

// QueryProgram ...
func QueryProgram(fn *CXFunction, expr *CXExpression, prgrmOffsetB []byte, affOffset *int) {
	prgrmOffset := AllocateSeq(OBJECT_HEADER_SIZE + I32_SIZE + I64_SIZE + STR_SIZE + I32_SIZE)
	// Callcounter
	WriteMemory(prgrmOffset+OBJECT_HEADER_SIZE, encoder.SerializeAtomic(int32(PROGRAM.CallCounter)))
	// HeapUsed
	WriteMemory(prgrmOffset+OBJECT_HEADER_SIZE+I32_SIZE, encoder.Serialize(int64(PROGRAM.HeapPointer)))

	// Caller
	if PROGRAM.CallCounter != 0 {
		// then it's not just entry point
		call := PROGRAM.CallStack[PROGRAM.CallCounter-1]

		var opNameB []byte
		if call.Operator.IsNative {
			opNameB = encoder.Serialize(OpNames[call.Operator.OpCode])
		} else {
			opNameB = encoder.Serialize(call.Operator.Package.Name + "." + call.Operator.Name)
		}

		opNameOffsetB := encoder.SerializeAtomic(int32(WriteObjectRetOff(opNameB)))

		// callOffset := AllocateSeq(OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE)
		// FnName
		WriteMemory(prgrmOffset+OBJECT_HEADER_SIZE+I32_SIZE+I64_SIZE, opNameOffsetB)
		// FnSize
		WriteMemory(prgrmOffset+OBJECT_HEADER_SIZE+I32_SIZE+I64_SIZE+STR_SIZE, encoder.SerializeAtomic(int32(call.Operator.Size)))

		// res := CallAffPredicate(fn, PROGRAM.Memory[callOffset + OBJECT_HEADER_SIZE : callOffset + OBJECT_HEADER_SIZE + STR_SIZE + I32_SIZE])

		// if res == 1 {
		// 	*affOffset = WriteToSlice(*affOffset, callerOffsetB)
		// }
	}

	res := CallAffPredicate(fn, PROGRAM.Memory[prgrmOffset+OBJECT_HEADER_SIZE:prgrmOffset+OBJECT_HEADER_SIZE+I32_SIZE+I64_SIZE+STR_SIZE+I32_SIZE])

	if res == 1 {
		*affOffset = WriteToSlice(*affOffset, prgrmOffsetB)
		*affOffset = WriteToSlice(*affOffset, prgrmOffsetB)
	}
}

func getTarget(inp2 *CXArgument, fp int, tgtElt *string, tgtArgType *string, tgtArgIndex *int,
	tgtPkg *CXPackage, tgtStrct *CXStruct, tgtFn *CXFunction, tgtExpr *CXExpression) {
	for _, aff := range GetInferActions(inp2, fp) {
		switch aff {
		case "prgrm":
			*tgtElt = "prgrm"
		case "pkg":
			*tgtElt = "pkg"
		case "strct":
			*tgtElt = "strct"
		case "fn":
			*tgtElt = "fn"
		case "expr":
			*tgtElt = "expr"
		case "rec":
			*tgtElt = "rec"
		case "inp":
			*tgtElt = "inp"
		case "out":
			*tgtElt = "out"
		default:
			switch *tgtElt {
			case "pkg":
				if pkg, err := PROGRAM.GetPackage(aff); err == nil {
					*tgtPkg = *pkg
				} else {
					panic(err)
				}
			case "fn":
				if fn, err := tgtPkg.GetFunction(aff); err == nil {
					*tgtFn = *fn
				} else {
					panic(err)
				}
			case "expr":
				if expr, err := tgtFn.GetExpressionByLabel(aff); err == nil {
					*tgtExpr = *expr
				} else {
					panic(err)
				}
			case "inp":
				*tgtArgType = "inp"
				*tgtElt = "arg"
				i, err := strconv.ParseInt(aff, 10, 32)

				*tgtArgIndex = int(i)

				if err != nil {
					panic(err)
				}
			case "out":
				*tgtArgType = "out"
				*tgtElt = "arg"
				i, err := strconv.ParseInt(aff, 10, 32)

				*tgtArgIndex = int(i)

				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func getAffordances(inp1 *CXArgument, fp int,
	tgtElt string, tgtArgType string, tgtArgIndex int,
	tgtPkg *CXPackage, tgtStrct *CXStruct, tgtFn *CXFunction, tgtExpr *CXExpression,
	affMsgs map[string]string,
	affs *[]string) {
	var fltrElt string
	elts := GetInferActions(inp1, fp)
	// for _, elt := range elts {
	for c := 0; c < len(elts); c++ {
		elt := elts[c]
		switch elt {
		case "arg":
			fltrElt = "arg"
		case "expr":
			fltrElt = "expr"
		case "fn":
			fltrElt = "fn"
		case "strct":
			fltrElt = "strct"
		case "pkg":
			fltrElt = "pkg"
		case "prgrm":
			fltrElt = "prgrm"
			// skipping the extra "prgrm"
			c++
			switch tgtElt {
			case "arg":
				if tgtArgType == "inp" {
					if msg, ok := affMsgs["prgrm-arg-input"]; ok {
						*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex))
					}
				} else {
					if msg, ok := affMsgs["prgrm-arg-output"]; ok {
						*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex))
					}
				}
			case "prgrm":
				*affs = append(*affs, "Run program")
			}
		default:
			switch fltrElt {
			case "arg":
				switch tgtElt {
				case "arg":
					if tgtArgType == "inp" {
						if msg, ok := affMsgs["arg-arg-input"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					} else {
						if msg, ok := affMsgs["arg-arg-output"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					}
				case "prgrm":
					*affs = append(*affs, "Print FA's value")
				}
			case "expr":
				if expr, err := tgtFn.GetExpressionByLabel(elt); err == nil {
					_ = expr
					switch tgtElt {
					case "arg":
						*affs = append(*affs, "Replace TA by FE")
					case "fn":
						*affs = append(*affs, "Add FE to TF")
					case "prgrm":
						*affs = append(*affs, "Call FE")
					}
				} else {
					panic(err)
				}
			case "fn":
				if fn, err := tgtPkg.GetFunction(elt); err == nil {
					_ = fn
					switch tgtElt {
					case "arg":
						*affs = append(*affs, "Replace TA by a call to FF")
					case "expr":
						*affs = append(*affs, "Change TE's operator to FF")
					case "pkg":
						// affs = append(affs, fmt.Sprintf("[%s.Operator = %s]", tgtExpr.Label, fn.Name))
						*affs = append(*affs, "Move FF to TP")
					case "prgrm":
						*affs = append(*affs, "Call FF")
					}
				} else {
					panic(err)
				}
			case "strct":
				switch tgtElt {
				case "arg":
					if msg, ok := affMsgs["arg-strct"]; ok {
						*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
					}
					if tgtArgType == "inp" {
						if msg, ok := affMsgs["strct-arg-input"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					} else {
						if msg, ok := affMsgs["strct-arg-output"]; ok {
							*affs = append(*affs, fmt.Sprintf(msg, tgtExpr.Label, tgtArgIndex, elt))
						}
					}
				case "fn":
					*affs = append(*affs, "Add or change TF's receiver to FS")
				case "pkg":
					*affs = append(*affs, "Move FS to TP")
				}
			case "pkg":
				if pkg, err := PROGRAM.GetPackage(elt); err == nil {
					_ = pkg
					switch tgtElt {
					case "pkg":
						*affs = append(*affs, "Make TP import FP")
					}
				} else {
					panic(err)
				}
				// case "prgrm":
				// 	switch tgtElt {
				// 	case "prgrm":
				// 		affs = append(affs, "Run program")
				// 	}
			}
		}
	}
}

func opAffOn(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]

	prevPkg := PROGRAM.CurrentPackage
	prevStrct := prevPkg.CurrentStruct
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	var tgtPkg = CXPackage(*prevPkg)
	var tgtStrct CXStruct
	if prevStrct != nil {
		tgtStrct = *prevStrct
	}
	var tgtFn = CXFunction(*expr.Function)
	var tgtExpr = CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp2, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtStrct, &tgtFn, &tgtExpr)

	// var affPkg *CXPackage = prevPkg
	// var affFn *CXFunction = prevFn
	// var affExpr *CXExpression = prevExpr

	// processing the affordances
	var affs []string
	getAffordances(inp1, fp, tgtElt, tgtArgType, tgtArgIndex, &tgtPkg, &tgtStrct, &tgtFn, &tgtExpr, onMessages, &affs)

	// returning to previous state
	PROGRAM.CurrentPackage = prevPkg
	PROGRAM.CurrentPackage.CurrentFunction = prevFn
	PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr

	for i, aff := range affs {
		fmt.Println(fmt.Sprintf("%d - %s", i, aff))
	}
}

func opAffOf(expr *CXExpression, fp int) {
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]

	prevPkg := PROGRAM.CurrentPackage
	prevStrct := prevPkg.CurrentStruct
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	var tgtPkg = CXPackage(*expr.Package)
	var tgtStrct CXStruct
	if prevStrct != nil {
		tgtStrct = *prevStrct
	}
	var tgtFn = CXFunction(*expr.Function)
	var tgtExpr = CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp2, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtStrct, &tgtFn, &tgtExpr)

	// processing the affordances
	var affs []string
	getAffordances(inp1, fp, tgtElt, tgtArgType, tgtArgIndex, &tgtPkg, &tgtStrct, &tgtFn, &tgtExpr, ofMessages, &affs)

	// returning to previous state
	PROGRAM.CurrentPackage = prevPkg
	PROGRAM.CurrentPackage.CurrentFunction = prevFn
	PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr

	for i, aff := range affs {
		fmt.Println(fmt.Sprintf("%d - %s", i, aff))
	}
}

func readStrctAff(aff string, tgtPkg *CXPackage) *CXStruct {
	strct, err := tgtPkg.GetStruct(aff)
	if err != nil {
		panic(err)
	}

	return strct
}

func readArgAff(aff string, tgtFn *CXFunction) *CXArgument {
	var affExpr *CXExpression
	var lIdx int
	var rIdx int
	var ch rune

	for _, ch = range aff {
		if ch == '.' {
			exprLbl := aff[lIdx:rIdx]
			for _, expr := range tgtFn.Expressions {
				if exprLbl == expr.Label {
					affExpr = expr
					rIdx++
					break
				}
			}
		}

		if affExpr != nil {
			break
		}

		rIdx++
	}

	lIdx = rIdx

	var argType string

	for _, ch = range aff[lIdx:] {
		if ch == '.' {
			argType = aff[lIdx:rIdx]
			rIdx++
			break
		}

		if argType != "" {
			break
		}

		rIdx++
	}

	lIdx = rIdx

	argIdx, err := strconv.ParseInt(aff[lIdx:], 10, 32)
	if err != nil {
		panic(err)
	}

	if argType == "Input" {
		return affExpr.Inputs[argIdx]
	}
	return affExpr.Outputs[argIdx]

}

func opAffInform(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]

	prevPkg := PROGRAM.CurrentPackage
	prevStrct := prevPkg.CurrentStruct
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	var tgtPkg = CXPackage(*prevPkg)
	var tgtStrct CXStruct
	if prevStrct != nil {
		tgtStrct = *prevStrct
	}
	var tgtFn = CXFunction(*expr.Function)
	var tgtExpr = CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp3, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtStrct, &tgtFn, &tgtExpr)

	elts := GetInferActions(inp1, fp)
	eltIdx := ReadI32(fp, inp2)
	eltType := elts[eltIdx*2]
	elt := elts[eltIdx*2+1]

	switch eltType {
	case "arg":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				tgtExpr.Inputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
			} else {
				tgtExpr.Outputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
			}
		case "strct":

		case "prgrm":

		}
	case "expr":
		if expr, err := tgtFn.GetExpressionByLabel(elt); err == nil {
			_ = expr
			switch tgtElt {
			case "arg":

			case "fn":

			case "prgrm":

			}
		} else {
			panic(err)
		}
	case "fn":
		if fn, err := tgtPkg.GetFunction(elt); err == nil {
			_ = fn
			switch tgtElt {
			case "arg":

			case "expr":

			case "pkg":

			case "prgrm":

			}
		} else {
			panic(err)
		}
	case "strct":
		switch tgtElt {
		case "arg":

		case "fn":

		case "pkg":

		}
	case "pkg":
		if pkg, err := PROGRAM.GetPackage(elt); err == nil {
			_ = pkg
			switch tgtElt {
			case "pkg":

			}
		} else {
			panic(err)
		}
		// case "prgrm":
		// 	switch tgtElt {
		// 	case "prgrm":
		// 		affs = append(affs, "Run program")
		// 	}
	}

	// returning to previous state
	PROGRAM.CurrentPackage = prevPkg
	PROGRAM.CurrentPackage.CurrentFunction = prevFn
	PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr
}

func opAffRequest(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]

	prevPkg := PROGRAM.CurrentPackage
	prevStrct := prevPkg.CurrentStruct
	prevFn := prevPkg.CurrentFunction
	prevExpr := prevFn.CurrentExpression

	var tgtPkg = CXPackage(*prevPkg)
	var tgtStrct CXStruct
	if prevStrct != nil {
		tgtStrct = *prevStrct
	}
	var tgtFn = CXFunction(*expr.Function)
	var tgtExpr = CXExpression(*prevExpr)

	// processing the target
	var tgtElt string
	var tgtArgType string
	var tgtArgIndex int

	getTarget(inp3, fp, &tgtElt, &tgtArgType, &tgtArgIndex, &tgtPkg, &tgtStrct, &tgtFn, &tgtExpr)

	// var affs []string

	elts := GetInferActions(inp1, fp)
	eltIdx := ReadI32(fp, inp2)
	eltType := elts[eltIdx*2]
	elt := elts[eltIdx*2+1]

	switch eltType {
	case "arg":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				// tgtExpr.Inputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
				*readArgAff(elt, &tgtFn) = *tgtExpr.Inputs[tgtArgIndex]
			} else {
				// tgtExpr.Outputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
				*readArgAff(elt, &tgtFn) = *tgtExpr.Outputs[tgtArgIndex]
			}
		case "strct":

		case "prgrm":
			fmt.Println(GetPrintableValue(fp, readArgAff(elt, &tgtFn)))
		}
	case "expr":
		if expr, err := tgtFn.GetExpressionByLabel(elt); err == nil {
			_ = expr
			switch tgtElt {
			case "arg":

			case "fn":

			case "prgrm":

			}
		} else {
			panic(err)
		}
	case "fn":
		if fn, err := tgtPkg.GetFunction(elt); err == nil {
			_ = fn
			switch tgtElt {
			case "arg":

			case "expr":

			case "pkg":

			case "prgrm":

			}
		} else {
			panic(err)
		}
	case "strct":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				// tgtExpr.Inputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
				readStrctAff(elt, &tgtPkg).AddField(tgtExpr.Inputs[tgtArgIndex])
			} else {
				// tgtExpr.Outputs[tgtArgIndex] = readArgAff(elt, &tgtFn)
				readStrctAff(elt, &tgtPkg).AddField(tgtExpr.Outputs[tgtArgIndex])
			}
		case "fn":

		case "pkg":

		}
	case "pkg":
		if pkg, err := PROGRAM.GetPackage(elt); err == nil {
			_ = pkg
			switch tgtElt {
			case "pkg":

			}
		} else {
			panic(err)
		}
	case "prgrm":
		switch tgtElt {
		case "arg":
			if tgtArgType == "inp" {
				fmt.Println(GetPrintableValue(fp, tgtExpr.Inputs[tgtArgIndex]))
			} else {
				fmt.Println(GetPrintableValue(fp, tgtExpr.Outputs[tgtArgIndex]))
			}
		case "prgrm":
			// affs = append(affs, "Run program")
		}
	}

	// returning to previous state
	PROGRAM.CurrentPackage = prevPkg
	PROGRAM.CurrentPackage.CurrentFunction = prevFn
	PROGRAM.CurrentPackage.CurrentFunction.CurrentExpression = prevExpr
}

func opAffQuery(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]

	out1Offset := GetFinalOffset(fp, out1)

	var affOffset int

	var cmd string
	for _, rule := range GetInferActions(inp1, fp) {
		switch rule {
		case "filter":
			cmd = "filter"
		case "sort":
			cmd = "sort"
		default:
			switch cmd {
			case "filter":
				if fn, err := inp1.Package.GetFunction(rule); err == nil {

					// arg keyword
					argB := encoder.Serialize("arg")
					argOffset := AllocateSeq(len(argB))
					WriteMemory(argOffset, argB)
					argOffsetB := encoder.SerializeAtomic(int32(argOffset))

					// expr keyword
					exprB := encoder.Serialize("expr")
					exprOffset := AllocateSeq(len(exprB))
					WriteMemory(exprOffset, exprB)
					exprOffsetB := encoder.SerializeAtomic(int32(exprOffset))

					// fn keyword
					fnB := encoder.Serialize("fn")
					fnOffset := AllocateSeq(len(fnB))
					WriteMemory(fnOffset, fnB)
					fnOffsetB := encoder.SerializeAtomic(int32(fnOffset))

					// strct keyword
					strctB := encoder.Serialize("strct")
					strctOffset := AllocateSeq(len(strctB))
					WriteMemory(strctOffset, strctB)
					strctOffsetB := encoder.SerializeAtomic(int32(strctOffset))

					// caller keyword
					callerB := encoder.Serialize("caller")
					callerOffset := AllocateSeq(len(callerB))
					WriteMemory(callerOffset, callerB)
					callerOffsetB := encoder.SerializeAtomic(int32(callerOffset))

					// program keyword
					prgrmB := encoder.Serialize("prgrm")
					prgrmOffset := AllocateSeq(len(prgrmB))
					WriteMemory(prgrmOffset, prgrmB)
					prgrmOffsetB := encoder.SerializeAtomic(int32(prgrmOffset))

					predInp := fn.Inputs[0]

					if predInp.Type == TYPE_CUSTOM {
						if predInp.CustomType != nil {
							switch predInp.CustomType.Name {
							case "Argument":
								QueryArgument(fn, expr, argOffsetB, &affOffset)
							case "Expression":
								QueryExpressions(fn, expr, exprOffsetB, &affOffset)
							case "Function":
								QueryFunction(fn, expr, fnOffsetB, &affOffset)
							case "Structure":
								QueryStructure(fn, expr, strctOffsetB, &affOffset)
							case "Caller":
								QueryCaller(fn, expr, callerOffsetB, &affOffset)
							case "Program":
								QueryProgram(fn, expr, prgrmOffsetB, &affOffset)
							}
						}
					}

				} else {
					panic(err)
				}
			case "sort":

			}
		}
	}

	WriteMemory(out1Offset, FromI32(int32(affOffset)))
}
