package base

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Debug ...
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

// ExprOpName ...
func ExprOpName(expr *CXExpression) string {
	if expr.Operator.IsNative {
		return OpNames[expr.Operator.OpCode]
	}
	return expr.Operator.Name

}

// func limitString (str string) string {
// 	if len(str) > 3
// }

func stackValueHeader(fileName string, fileLine int) string {
	return fmt.Sprintf("%s:%d", fileName, fileLine)
}

// PrintStack ...
func (prgrm *CXProgram) PrintStack() {
	fmt.Println()
	fmt.Println("===Callstack===")

	// we're going backwards in the stack
	fp := prgrm.StackPointer

	for c := prgrm.CallCounter; c >= 0; c-- {
		op := prgrm.CallStack[c].Operator
		fp -= op.Size

		var dupNames []string

		fmt.Printf(">>> %s()\n", op.Name)

		for _, inp := range op.Inputs {
			fmt.Println("Inputs")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(inp.FileName, inp.FileLine), op.Name, GetPrintableValue(fp, inp))

			dupNames = append(dupNames, inp.Package.Name+inp.Name)
		}

		for _, out := range op.Outputs {
			fmt.Println("Outputs")
			fmt.Printf("\t%s : %s() : %s\n", stackValueHeader(out.FileName, out.FileLine), op.Name, GetPrintableValue(fp, out))

			dupNames = append(dupNames, out.Package.Name+out.Name)
		}

		// fmt.Println("Expressions")
		exprs := ""
		for _, expr := range op.Expressions {
			for _, inp := range expr.Inputs {
				if inp.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == inp.Package.Name+inp.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", inp.Name, "\t", ":", "\t", GetPrintableValue(fp, inp))
				// exprs += fmt.Sprintln("\t", stackValueHeader(inp.FileName, inp.FileLine), "\t", ":", "\t", GetPrintableValue(fp, inp))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(inp.FileName, inp.FileLine), ExprOpName(expr), GetPrintableValue(fp, inp))

				dupNames = append(dupNames, inp.Package.Name+inp.Name)
			}

			for _, out := range expr.Outputs {
				if out.Name == "" || expr.Operator == nil {
					continue
				}
				var dup bool
				for _, name := range dupNames {
					if name == out.Package.Name+out.Name {
						dup = true
						break
					}
				}
				if dup {
					continue
				}

				// fmt.Println("\t", out.Name, "\t", ":", "\t", GetPrintableValue(fp, out))
				// exprs += fmt.Sprintln("\t", stackValueHeader(out.FileName, out.FileLine), ":", GetPrintableValue(fp, out))

				exprs += fmt.Sprintf("\t%s : %s() : %s\n", stackValueHeader(out.FileName, out.FileLine), ExprOpName(expr), GetPrintableValue(fp, out))

				dupNames = append(dupNames, out.Package.Name+out.Name)
			}
		}

		if len(exprs) > 0 {
			fmt.Println("Expressions\n", exprs)
		}
	}
}

// PrintProgram ...
func (prgrm *CXProgram) PrintProgram() {
	fmt.Println("Program")

	var currentFunction *CXFunction
	var currentPackage *CXPackage

	_ = currentFunction
	_ = currentPackage

	// saving current program state because PrintProgram uses SelectXXX
	if pkg, err := prgrm.GetCurrentPackage(); err == nil {
		currentPackage = pkg
	}

	if fn, err := prgrm.GetCurrentFunction(); err == nil {
		currentFunction = fn
	}

	i := 0

	for _, mod := range prgrm.Packages {
		if IsCorePackage(mod.Name) {
			continue
		}

		fmt.Printf("%d.- Package: %s\n", i, mod.Name)

		if len(mod.Imports) > 0 {
			fmt.Println("\tImports")
		}

		j := 0
		for _, imp := range mod.Imports {
			fmt.Printf("\t\t%d.- Import: %s\n", j, imp.Name)
			j++
		}

		if len(mod.Globals) > 0 {
			fmt.Println("\tGlobals")
		}

		j = 0
		for _, v := range mod.Globals {
			var arrayStr string
			if v.IsArray {
				for _, l := range v.Lengths {
					arrayStr += fmt.Sprintf("[%d]", l)
				}
			}
			fmt.Printf("\t\t%d.- Global: %s %s%s\n", j, v.Name, arrayStr, TypeNames[v.Type])
			j++
		}

		if len(mod.Structs) > 0 {
			fmt.Println("\tStructs")
		}

		j = 0
		for _, strct := range mod.Structs {
			fmt.Printf("\t\t%d.- Struct: %s\n", j, strct.Name)

			for k, fld := range strct.Fields {
				fmt.Printf("\t\t\t%d.- Field: %s %s\n",
					k, fld.Name, TypeNames[fld.Type])
			}

			j++
		}

		if len(mod.Functions) > 0 {
			fmt.Println("\tFunctions")
		}

		j = 0
		for _, fn := range mod.Functions {
			mod.SelectFunction(fn.Name)

			var inps bytes.Buffer
			for i, inp := range fn.Inputs {
				var isPointer string
				if inp.IsPointer {
					isPointer = "*"
				}

				var arrayStr string
				if inp.IsArray {
					for _, l := range inp.Lengths {
						arrayStr += fmt.Sprintf("[%d]", l)
					}
				}

				var typeName string
				elt := GetAssignmentElement(inp)
				if elt.CustomType != nil {
					// then it's custom type
					typeName = elt.CustomType.Name
				} else {
					// then it's native type
					typeName = TypeNames[elt.Type]
				}

				if i == len(fn.Inputs)-1 {
					inps.WriteString(fmt.Sprintf("%s %s%s%s", inp.Name, isPointer, arrayStr, typeName))
				} else {
					inps.WriteString(fmt.Sprintf("%s %s%s%s, ", inp.Name, isPointer, arrayStr, typeName))
				}
			}

			var outs bytes.Buffer
			for i, out := range fn.Outputs {
				var isPointer string
				if out.IsPointer {
					isPointer = "*"
				}

				var arrayStr string
				if out.IsArray {
					for _, l := range out.Lengths {
						arrayStr += fmt.Sprintf("[%d]", l)
					}
				}

				var typeName string
				elt := GetAssignmentElement(out)
				if elt.CustomType != nil {
					// then it's custom type
					typeName = elt.CustomType.Name
				} else {
					// then it's native type
					typeName = TypeNames[elt.Type]
				}

				if i == len(fn.Outputs)-1 {
					outs.WriteString(fmt.Sprintf("%s %s%s%s", out.Name, isPointer, arrayStr, typeName))
				} else {
					outs.WriteString(fmt.Sprintf("%s %s%s%s, ", out.Name, isPointer, arrayStr, typeName))
				}
			}

			fmt.Printf("\t\t%d.- Function: %s (%s) (%s)\n",
				j, fn.Name, inps.String(), outs.String())

			k := 0
			for _, expr := range fn.Expressions {
				// if expr.Operator == nil {
				//      continue
				// }
				//Arguments
				var args bytes.Buffer

				for i, arg := range expr.Inputs {
					var name string
					var dat []byte

					if arg.Offset > STACK_SIZE {
						dat = prgrm.Memory[arg.Offset : arg.Offset+arg.Size]
					} else {
						name = arg.Name
					}

					if dat != nil {
						switch TypeNames[arg.Type] {
						case "str":
							encoder.DeserializeRaw(dat, &name)
							name = "\"" + name + "\""
						case "i32":
							var i32 int32
							encoder.DeserializeAtomic(dat, &i32)
							name = fmt.Sprintf("%v", i32)
						case "i64":
							var i64 int64
							encoder.DeserializeRaw(dat, &i64)
							name = fmt.Sprintf("%v", i64)
						case "f32":
							var f32 float32
							encoder.DeserializeRaw(dat, &f32)
							name = fmt.Sprintf("%v", f32)
						case "f64":
							var f64 float64
							encoder.DeserializeRaw(dat, &f64)
							name = fmt.Sprintf("%v", f64)
						case "bool":
							var b bool
							encoder.DeserializeRaw(dat, &b)
							name = fmt.Sprintf("%v", b)
						case "byte":
							var b bool
							encoder.DeserializeRaw(dat, &b)
							name = fmt.Sprintf("%v", b)
						}
					}

					if arg.Name != "" {
						name = arg.Name
						for _, fld := range arg.Fields {
							name += "." + fld.Name
						}
					}

					var derefLevels string
					if arg.DereferenceLevels > 0 {
						for c := 0; c < arg.DereferenceLevels; c++ {
							derefLevels += "*"
						}
					}

					var isReference string
					if arg.PassBy == PASSBY_REFERENCE {
						isReference = "&"
					}

					var arrayStr string
					if arg.IsArray {
						for _, l := range arg.Lengths {
							arrayStr += fmt.Sprintf("[%d]", l)
						}
					}

					var typeName string
					elt := GetAssignmentElement(arg)
					if elt.CustomType != nil {
						// then it's custom type
						typeName = elt.CustomType.Name
					} else {
						// then it's native type
						typeName = TypeNames[elt.Type]
					}

					if i == len(expr.Inputs)-1 {
						args.WriteString(fmt.Sprintf("%s%s%s %s%s", isReference, derefLevels, name, arrayStr, typeName))
					} else {
						args.WriteString(fmt.Sprintf("%s%s%s %s%s, ", isReference, derefLevels, name, arrayStr, typeName))
					}
				}

				var opName string
				if expr.Operator != nil {
					if expr.Operator.IsNative {
						opName = OpNames[expr.Operator.OpCode]
					} else {
						opName = expr.Operator.Name
					}
				}

				if len(expr.Outputs) > 0 {
					var outNames bytes.Buffer
					for i, outName := range expr.Outputs {
						out := GetAssignmentElement(outName)

						var derefLevels string
						if outName.DereferenceLevels > 0 {
							for c := 0; c < outName.DereferenceLevels; c++ {
								derefLevels += "*"
							}
						}

						var arrayStr string
						if outName.IsArray {
							for _, l := range outName.Lengths {
								arrayStr += fmt.Sprintf("[%d]", l)
							}
						}

						var typeName string
						if out.CustomType != nil {
							// then it's custom type
							typeName = out.CustomType.Name
						} else {
							// then it's native type
							typeName = TypeNames[out.Type]
						}

						fullName := outName.Name

						for _, fld := range outName.Fields {
							fullName += "." + fld.Name
						}

						if i == len(expr.Outputs)-1 {
							outNames.WriteString(fmt.Sprintf("%s%s%s %s", derefLevels, fullName, arrayStr, typeName))
						} else {
							outNames.WriteString(fmt.Sprintf("%s%s%s %s, ", derefLevels, fullName, arrayStr, typeName))
						}
					}

					var lbl string
					if expr.Label != "" {
						lbl = " <<" + expr.Label + ">>"
					} else {
						lbl = ""
					}

					if expr.Operator != nil {
						fmt.Printf("\t\t\t%d.- Expression%s: %s = %s(%s)\n",
							k,
							lbl,
							outNames.String(),
							opName,
							args.String(),
						)
					} else {
						if len(expr.Outputs) > 0 {
							var typ string

							out := expr.Outputs[len(expr.Outputs)-1]

							// NOTE: this only adds a single *, regardless of how many
							// dereferences we have. won't be fixed atm, as the whole
							// PrintProgram needs to be refactored soon.
							if out.IsPointer {
								typ = "*"
							}

							if GetAssignmentElement(out).CustomType != nil {
								// then it's custom type
								typ += GetAssignmentElement(out).CustomType.Name
							} else {
								// then it's native type
								typ += TypeNames[GetAssignmentElement(out).Type]
							}

							fmt.Printf("\t\t\t%d.- Declaration%s: %s %s\n",
								k,
								lbl,
								expr.Outputs[0].Name,
								typ)
						}
					}

				} else {
					var lbl string

					if expr.Label != "" {
						lbl = " <<" + expr.Label + ">>"
					} else {
						lbl = ""
					}

					fmt.Printf("\t\t\t%d.- Expression%s: %s(%s)\n",
						k,
						lbl,
						opName,
						args.String(),
					)
				}
				k++
			}
			j++
		}
		i++
	}

	if currentPackage != nil {
		prgrm.SelectPackage(currentPackage.Name)
	}
	if currentFunction != nil {
		prgrm.SelectFunction(currentFunction.Name)
	}

	prgrm.CurrentPackage = currentPackage
	currentPackage.CurrentFunction = currentFunction
}

// CheckArithmeticOp ...
func CheckArithmeticOp(expr *CXExpression) bool {
	if expr.Operator.IsNative {
		switch expr.Operator.OpCode {
		case OP_I32_MUL, OP_I32_DIV, OP_I32_MOD, OP_I32_ADD,
			OP_I32_SUB, OP_I32_NEG, OP_I32_BITSHL, OP_I32_BITSHR, OP_I32_LT,
			OP_I32_GT, OP_I32_LTEQ, OP_I32_GTEQ, OP_I32_EQ, OP_I32_UNEQ,
			OP_I32_BITAND, OP_I32_BITXOR, OP_I32_BITOR, OP_STR_EQ:
			return true
		}
	}
	return false
}

// IsCorePackage ...
func IsCorePackage(ident string) bool {
	for _, core := range CorePackages {
		if core == ident {
			return true
		}
	}
	return false
}

// IsTempVar ...
func IsTempVar(name string) bool {
	if len(name) >= len(LOCAL_PREFIX) && name[:len(LOCAL_PREFIX)] == LOCAL_PREFIX {
		return true
	}
	return false
}

// GetArgSize ...
func GetArgSize(typ int) int {
	switch typ {
	case TYPE_BOOL, TYPE_BYTE:
		return 1
	case TYPE_STR, TYPE_I32, TYPE_F32, TYPE_AFF:
		return 4
	case TYPE_I64, TYPE_F64:
		return 8
	default:
		return 4
	}
}

func checkForEscapedChars(str string) []byte {
	var res []byte
	var lenStr = int(len(str))
	for c := 0; c < len(str); c++ {
		var nextCh byte
		ch := str[c]
		if c < lenStr-1 {
			nextCh = str[c+1]
		}
		if ch == '\\' {
			switch nextCh {
			case '%':
				c++
				res = append(res, nextCh)
				continue
			case 'n':
				c++
				res = append(res, '\n')
				continue
			default:
				res = append(res, ch)
				continue
			}

		} else {
			res = append(res, ch)
		}
	}

	return res
}

// GetAssignmentElement ...
func GetAssignmentElement(arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1]
	}
	return arg

}

// IsValidSliceIndex ...
func IsValidSliceIndex(offset int, index int, sizeofElement int) bool {
	sliceLen := GetSliceLen(int32(offset))
	bytesLen := sliceLen * int32(sizeofElement)
	index -= OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + offset

	if index >= 0 && index < int(bytesLen) && (index%sizeofElement) == 0 {
		return true
	}
	return false
}

// GetPointerOffset ...
func GetPointerOffset(pointer int32) int32 {
	var offset int32
	encoder.DeserializeAtomic(PROGRAM.Memory[pointer:pointer+TYPE_POINTER_SIZE], &offset)
	return offset
}

// GetSliceOffset ...
func GetSliceOffset(fp int, arg *CXArgument) int32 {
	element := GetAssignmentElement(arg)
	if element.IsSlice {
		return GetPointerOffset(int32(GetFinalOffset(fp, arg)))
	}

	return -1
}

// GetObjectHeader ...
func GetObjectHeader(offset int32) []byte {
	return PROGRAM.Memory[offset : offset+OBJECT_HEADER_SIZE]
}

// GetSliceHeader ...
func GetSliceHeader(offset int32) []byte {
	return PROGRAM.Memory[offset+OBJECT_HEADER_SIZE : offset+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]
}

// GetSliceLen ...
func GetSliceLen(offset int32) int32 {
	var sliceLen int32
	sliceHeader := GetSliceHeader(offset)
	encoder.DeserializeAtomic(sliceHeader[4:8], &sliceLen)
	return sliceLen
}

// GetSlice ...
func GetSlice(offset int32, sizeofElement int) []byte {
	if offset > 0 {
		sliceLen := GetSliceLen(offset)
		if sliceLen > 0 {
			dataOffset := offset + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE - 4
			dataLen := 4 + sliceLen*int32(sizeofElement)
			return PROGRAM.Memory[dataOffset : dataOffset+dataLen]
		}
	}
	return nil
}

// GetSliceData ...
func GetSliceData(offset int32, sizeofElement int) []byte {
	if slice := GetSlice(offset, sizeofElement); slice != nil {
		return slice[4:]
	}
	return nil
}

// SliceResize ...
func SliceResize(outputSliceOffset int32, inputSliceOffset int32, count int32, sizeofElement int) int {
	if count < 0 {
		panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE) // TODO : should use uint32
	}

	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
		//inputSliceHeader := GetSliceHeader(inputSliceOffset)
		//encoder.DeserializeAtomic(inputSliceHeader[4:8], &inputSliceLen)
	}

	var outputSliceHeader []byte
	var outputSliceLen int32
	var outputSliceCap int32

	if outputSliceOffset > 0 {
		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		encoder.DeserializeAtomic(outputSliceHeader[0:4], &outputSliceCap)
		encoder.DeserializeAtomic(outputSliceHeader[4:8], &outputSliceLen)
	}

	var newLen = count
	var newCap = outputSliceCap
	if newLen > newCap {
		if newCap <= 0 {
			newCap = newLen
		} else {
			newCap *= 2
		}
		var outputObjectSize = OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE + newCap*int32(sizeofElement)
		outputSliceOffset = int32(AllocateSeq(int(outputObjectSize)))
		copy(GetObjectHeader(outputSliceOffset)[5:9], encoder.SerializeAtomic(outputObjectSize))

		outputSliceHeader = GetSliceHeader(outputSliceOffset)
		copy(outputSliceHeader[0:4], encoder.SerializeAtomic(newCap))
	}

	if outputSliceOffset > 0 {
		copy(outputSliceHeader[4:8], encoder.SerializeAtomic(newLen))
		outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)

		if (outputSliceOffset != inputSliceOffset) && inputSliceLen > 0 {
			copy(outputSliceData, GetSliceData(inputSliceOffset, sizeofElement))
		}
	}

	return int(outputSliceOffset)
}

// SliceAppend ...
func SliceAppend(outputSliceOffset int32, inputSliceOffset int32, object []byte) int {
	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
		//inputSliceHeader := GetSliceHeader(inputSliceOffset)
		//encoder.DeserializeAtomic(inputSliceHeader[4:8], &inputSliceLen)
	}

	sizeofElement := len(object)
	outputSliceOffset = int32(SliceResize(outputSliceOffset, inputSliceOffset, inputSliceLen+1, sizeofElement))
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[int(inputSliceLen)*sizeofElement:], object)
	return int(outputSliceOffset)
}

// SliceInsert ...
func SliceInsert(outputSliceOffset int32, inputSliceOffset int32, index int32, object []byte) int {
	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index > inputSliceLen {
		panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	var newLen = inputSliceLen + 1
	sizeofElement := len(object)
	outputSliceOffset = int32(SliceResize(outputSliceOffset, inputSliceOffset, newLen, sizeofElement))
	outputSliceData := GetSliceData(outputSliceOffset, sizeofElement)
	copy(outputSliceData[int(index+1)*sizeofElement:], outputSliceData[int(index)*sizeofElement:])
	copy(outputSliceData[int(index)*sizeofElement:], object)
	return int(outputSliceOffset)
}

// SliceRemove ...
func SliceRemove(outputSliceOffset int32, inputSliceOffset int32, index int32, sizeofElement int32) int {
	var inputSliceLen int32
	if inputSliceOffset != 0 {
		inputSliceLen = GetSliceLen(inputSliceOffset)
	}

	if index < 0 || index >= inputSliceLen {
		panic(CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	}

	outputSliceData := GetSliceData(outputSliceOffset, int(sizeofElement))
	copy(outputSliceData[index*sizeofElement:], outputSliceData[(index+1)*sizeofElement:])
	outputSliceOffset = int32(SliceResize(outputSliceOffset, inputSliceOffset, inputSliceLen-1, int(sizeofElement)))
	return int(outputSliceOffset)
}

// WriteToSlice ...
func WriteToSlice(off int, inp []byte) int {
	return SliceAppend(int32(off), int32(off), inp)
}

// refactoring reuse in WriteObject and WriteObjectRetOff
func writeObj(obj []byte) int {
	size := len(obj)
	sizeB := encoder.SerializeAtomic(int32(size))
	// heapOffset := AllocateSeq(size + OBJECT_HEADER_SIZE + SLICE_HEADER_SIZE)
	heapOffset := AllocateSeq(size + OBJECT_HEADER_SIZE)

	var finalObj = make([]byte, OBJECT_HEADER_SIZE+size)

	for c := OBJECT_GC_HEADER_SIZE; c < OBJECT_HEADER_SIZE; c++ {
		finalObj[c] = sizeB[c-OBJECT_GC_HEADER_SIZE]
	}
	for c := OBJECT_HEADER_SIZE; c < size+OBJECT_HEADER_SIZE; c++ {
		finalObj[c] = obj[c-OBJECT_HEADER_SIZE]
	}

	WriteMemory(heapOffset, finalObj)
	return heapOffset + OBJECT_HEADER_SIZE
}

// WriteObject ...
func WriteObject(out1Offset int, obj []byte) {
	off := encoder.SerializeAtomic(int32(writeObj(obj)))

	WriteMemory(out1Offset, off)
}

// WriteObjectRetOff ...
func WriteObjectRetOff(obj []byte) int {
	return writeObj(obj)
}

// ErrorHeader ...
func ErrorHeader(currentFile string, lineNo int) string {
	return "error: " + currentFile + ":" + strconv.FormatInt(int64(lineNo), 10)
}

// CompilationError is a helper function that concatenates the `currentFile` and `lineNo` data to a error header and returns the full error string.
func CompilationError(currentFile string, lineNo int) string {
	FoundCompileErrors = true
	return ErrorHeader(currentFile, lineNo)
}

// ErrorString ...
func ErrorString(code int) string {
	if str, found := ErrorStrings[code]; found {
		return str
	}
	return ErrorStrings[CX_RUNTIME_ERROR]
}

func errorCode(r interface{}) int {
	switch v := r.(type) {
	case int:
		return int(v)
	default:
		return CX_RUNTIME_ERROR
	}
}

func runtimeErrorInfo(r interface{}, printStack bool, defaultError int) {
	call := PROGRAM.CallStack[PROGRAM.CallCounter]
	expr := call.Operator.Expressions[call.Line]
	code := errorCode(r)
	if code == CX_RUNTIME_ERROR {
		code = defaultError
	}

	fmt.Printf("%s, %s, %v", ErrorHeader(expr.FileName, expr.FileLine), ErrorString(code), r)

	if printStack {
		PROGRAM.PrintStack()
	}

	if DBG_GOLANG_STACK_TRACE {
		debug.PrintStack()
	}

	os.Exit(code)
}

// RuntimeError ...
func RuntimeError() {
	if r := recover(); r != nil {
		switch r {
		case STACK_OVERFLOW_ERROR:
			call := PROGRAM.CallStack[PROGRAM.CallCounter]
			if PROGRAM.CallCounter > 0 {
				PROGRAM.CallCounter--
				PROGRAM.StackPointer = call.FramePointer
				runtimeErrorInfo(r, true, CX_RUNTIME_STACK_OVERFLOW_ERROR)
			} else {
				// error at entry point
				runtimeErrorInfo(r, false, CX_RUNTIME_STACK_OVERFLOW_ERROR)
			}
		case HEAP_EXHAUSTED_ERROR:
			runtimeErrorInfo(r, true, CX_RUNTIME_HEAP_EXHAUSTED_ERROR)
		default:
			runtimeErrorInfo(r, true, CX_RUNTIME_ERROR)
		}
	}
}

func getNonCollectionValue(fp int, arg, elt *CXArgument, typ string) string {
	switch typ {
	case "bool":
		return fmt.Sprintf("%v", ReadBool(fp, elt))
	case "byte":
		return fmt.Sprintf("%v", ReadByte(fp, elt))
	case "str":
		return fmt.Sprintf("%v", ReadStr(fp, elt))
	case "i32":
		return fmt.Sprintf("%v", ReadI32(fp, elt))
	case "i64":
		return fmt.Sprintf("%v", ReadI64(fp, elt))
	case "f32":
		return fmt.Sprintf("%v", ReadF32(fp, elt))
	case "f64":
		return fmt.Sprintf("%v", ReadF64(fp, elt))
	default:
		// then it's a struct
		var val string
		val = "{"
		// for _, fld := range elt.CustomType.Fields {
		lFlds := len(elt.CustomType.Fields)
		off := 0
		for c := 0; c < lFlds; c++ {
			fld := elt.CustomType.Fields[c]
			if c == lFlds-1 {
				val += fmt.Sprintf("%s: %s", fld.Name, GetPrintableValue(fp+arg.Offset+off, fld))
			} else {
				val += fmt.Sprintf("%s: %s, ", fld.Name, GetPrintableValue(fp+arg.Offset+off, fld))
			}
			off += fld.TotalSize
		}
		val += "}"
		return val
	}
}

// GetPrintableValue ...
func GetPrintableValue(fp int, arg *CXArgument) string {
	var typ string
	elt := GetAssignmentElement(arg)
	if elt.CustomType != nil {
		// then it's custom type
		typ = elt.CustomType.Name
	} else {
		// then it's native type
		typ = TypeNames[elt.Type]
	}

	if len(elt.Lengths) > 0 {
		var val string
		if len(elt.Lengths) == 1 {
			val = "["
			for c := 0; c < elt.Lengths[0]; c++ {
				if c == elt.Lengths[0]-1 {
					val += getNonCollectionValue(fp+c*elt.Size, arg, elt, typ)
				} else {
					val += getNonCollectionValue(fp+c*elt.Size, arg, elt, typ) + ", "
				}

			}
			val += "]"
		} else {
			// 5, 4, 1
			val = ""

			finalSize := 1
			for _, l := range elt.Lengths {
				finalSize *= l
			}

			lens := make([]int, len(elt.Lengths))
			copy(lens, elt.Lengths)

			for c := 0; c < len(lens); c++ {
				for i := 0; i < len(lens[c+1:]); i++ {
					lens[c] *= lens[c+i]
				}
			}

			for range lens {
				val += "["
			}

			// adding first element because of formatting reasons
			val += getNonCollectionValue(fp, arg, elt, typ)
			for c := 1; c < finalSize; c++ {
				closeCount := 0
				for _, l := range lens {
					if c%l == 0 && c != 0 {
						// val += "] ["
						closeCount++
					}
				}

				if closeCount > 0 {
					for c := 0; c < closeCount; c++ {
						val += "]"
					}
					val += " "
					for c := 0; c < closeCount; c++ {
						val += "["
					}

					val += getNonCollectionValue(fp+c*elt.Size, arg, elt, typ)
				} else {
					val += " " + getNonCollectionValue(fp+c*elt.Size, arg, elt, typ)
				}
			}
			for range lens {
				val += "]"
			}
		}

		return val
	}

	return getNonCollectionValue(fp, arg, elt, typ)
}
