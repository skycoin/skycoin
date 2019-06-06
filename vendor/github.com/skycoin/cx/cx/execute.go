package base

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// It "un-runs" a program
// func (prgrm *CXProgram) Reset() {
// 	prgrm.CallStack = MakeCallStack(0)
// 	prgrm.Steps = make([][]CXCall, 0)
// 	prgrm.Outputs = make([]*CXArgument, 0)
// 	//prgrm.ProgramSteps = nil
// }

// UnRun ...
func (cxt *CXProgram) UnRun(nCalls int) {
	if nCalls >= 0 || cxt.CallCounter < 0 {
		return
	}

	call := &cxt.CallStack[cxt.CallCounter]

	for c := nCalls; c < 0; c++ {
		if call.Line >= c {
			// then we stay in this call counter
			call.Line += c
			c -= c
		} else {

			if cxt.CallCounter == 0 {
				call.Line = 0
				return
			}
			c += call.Line
			call.Line = 0
			cxt.CallCounter--
			call = &cxt.CallStack[cxt.CallCounter]
		}
	}
}

// ToCall ...
func (cxt *CXProgram) ToCall() *CXExpression {
	for c := cxt.CallCounter - 1; c >= 0; c-- {
		if cxt.CallStack[c].Line+1 >= len(cxt.CallStack[c].Operator.Expressions) {
			// then it'll also return from this function call; continue
			continue
		}
		return cxt.CallStack[c].Operator.Expressions[cxt.CallStack[c].Line+1]
		// prgrm.CallStack[c].Operator.Expressions[prgrm.CallStack[prgrm.CallCounter-1].Line + 1]
	}
	// error
	return &CXExpression{Operator: MakeFunction("", "", -1)}
	// panic("")
}

// Run ...
func (cxt *CXProgram) Run(untilEnd bool, nCalls *int, untilCall int) error {
	defer RuntimeError()
	var err error

	for !cxt.Terminated && (untilEnd || *nCalls != 0) && cxt.CallCounter > untilCall {
		call := &cxt.CallStack[cxt.CallCounter]

		// checking if enough memory in stack
		if cxt.StackPointer > STACK_SIZE {
			panic(STACK_OVERFLOW_ERROR)
		}

		if !untilEnd {
			var inName string
			var toCallName string
			var toCall *CXExpression

			if call.Line >= call.Operator.Length && cxt.CallCounter == 0 {
				cxt.Terminated = true
				cxt.CallStack[0].Operator = nil
				cxt.CallCounter = 0
				fmt.Println("in:terminated")
				return err
			}

			if call.Line >= call.Operator.Length && cxt.CallCounter != 0 {
				toCall = cxt.ToCall()
				// toCall = prgrm.CallStack[prgrm.CallCounter-1].Operator.Expressions[prgrm.CallStack[prgrm.CallCounter-1].Line + 1]
				inName = cxt.CallStack[cxt.CallCounter-1].Operator.Name
			} else {
				toCall = call.Operator.Expressions[call.Line]
				inName = call.Operator.Name
			}

			if toCall.Operator == nil {
				// then it's a declaration
				toCallName = "declaration"
			} else if toCall.Operator.IsNative {
				toCallName = OpNames[toCall.Operator.OpCode]
			} else {
				if toCall.Operator.Name != "" {
					toCallName = toCall.Operator.Package.Name + "." + toCall.Operator.Name
				} else {
					// then it's the end of the program got from nested function calls
					cxt.Terminated = true
					cxt.CallStack[0].Operator = nil
					cxt.CallCounter = 0
					fmt.Println("in:terminated")
					return err
				}
			}

			fmt.Printf("in:%s, expr#:%d, calling:%s()\n", inName, call.Line+1, toCallName)
			*nCalls--
		}

		err = call.ccall(cxt)
		if err != nil {
			return err
		}
	}

	return nil
}

// RunCompiled ...
func (cxt *CXProgram) RunCompiled(nCalls int, args []string) error {
	PROGRAM = cxt
	// prgrm.PrintProgram()
	rand.Seed(time.Now().UTC().UnixNano())

	var untilEnd bool
	if nCalls == 0 {
		untilEnd = true
	}
	mod, err := cxt.SelectPackage(MAIN_PKG)
	if err == nil {
		// initializing program resources
		// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))

		if cxt.CallStack[0].Operator == nil {
			// then the program is just starting and we need to run the SYS_INIT_FUNC
			if fn, err := mod.SelectFunction(SYS_INIT_FUNC); err == nil {
				// *init function
				mainCall := MakeCall(fn)
				cxt.CallStack[0] = mainCall
				cxt.StackPointer = fn.Size

				var err error

				for !cxt.Terminated {
					call := &cxt.CallStack[cxt.CallCounter]
					err = call.ccall(cxt)
					if err != nil {
						return err
					}
				}
				// we reset call state
				cxt.Terminated = false
				cxt.CallCounter = 0
				cxt.CallStack[0].Operator = nil
			} else {
				return err
			}
		}

		if fn, err := mod.SelectFunction(MAIN_FUNC); err == nil {
			if len(fn.Expressions) < 1 {
				return nil
			}

			if cxt.CallStack[0].Operator == nil {
				// main function
				mainCall := MakeCall(fn)
				mainCall.FramePointer = cxt.StackPointer
				// initializing program resources
				cxt.CallStack[0] = mainCall

				// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))
				cxt.StackPointer += fn.Size

				// feeding os.Args
				if osPkg, err := PROGRAM.SelectPackage(OS_PKG); err == nil {
					argsOffset := 0
					if osGbl, err := osPkg.GetGlobal(OS_ARGS); err == nil {
						for _, arg := range args {
							argBytes := encoder.Serialize(arg)
							argOffset := AllocateSeq(len(argBytes))
							WriteMemory(argOffset, argBytes)
							argOffsetBytes := encoder.SerializeAtomic(int32(argOffset))
							argsOffset = WriteToSlice(argsOffset, argOffsetBytes)
						}
						WriteMemory(GetFinalOffset(0, osGbl), FromI32(int32(argsOffset)))
					}
				}
				cxt.Terminated = false
			}

			if err = cxt.Run(untilEnd, &nCalls, -1); err != nil {
				return err
			}

			if cxt.Terminated {
				cxt.Terminated = false
				cxt.CallCounter = 0
				cxt.CallStack[0].Operator = nil
			}

			// debugging memory
			if len(cxt.Memory) < 2000 {
				fmt.Println("prgrm.Memory", cxt.Memory)
			}

			return err
		}
		return err

	}
	return err

}

func (cxt *CXProgram) ccallback(expr *CXExpression, functionName string, packageName string, inputs [][]byte) {
	if fn, err := cxt.GetFunction(functionName, packageName); err == nil {
		line := cxt.CallStack[cxt.CallCounter].Line
		previousCall := cxt.CallCounter
		cxt.CallCounter++
		newCall := &cxt.CallStack[cxt.CallCounter]
		newCall.Operator = fn
		newCall.Line = 0
		newCall.FramePointer = cxt.StackPointer
		cxt.StackPointer += newCall.Operator.Size
		newFP := newCall.FramePointer

		// wiping next mem frame (removing garbage)
		for c := 0; c < expr.Operator.Size; c++ {
			cxt.Memory[newFP+c] = 0
		}

		for i, inp := range inputs {
			WriteMemory(GetFinalOffset(newFP, newCall.Operator.Inputs[i]), inp)
		}

		var nCalls = 0
		if err := cxt.Run(true, &nCalls, previousCall); err != nil {
			os.Exit(CX_INTERNAL_ERROR)
		}

		cxt.CallCounter = previousCall
		cxt.CallStack[cxt.CallCounter].Line = line
	}
}

func (call *CXCall) ccall(prgrm *CXProgram) error {
	// CX is still single-threaded, so only one stack
	if call.Line >= call.Operator.Length {
		/*
		   popping the stack
		*/
		// going back to the previous call
		prgrm.CallCounter--
		if prgrm.CallCounter < 0 {
			// then the program finished
			prgrm.Terminated = true
		} else {
			// copying the outputs to the previous stack frame
			returnAddr := &prgrm.CallStack[prgrm.CallCounter]
			returnOp := returnAddr.Operator
			returnLine := returnAddr.Line
			returnFP := returnAddr.FramePointer
			fp := call.FramePointer

			expr := returnOp.Expressions[returnLine]

			// lenOuts := len(expr.Outputs)
			for i, out := range call.Operator.Outputs {
				WriteMemory(
					GetFinalOffset(returnFP, expr.Outputs[i]),
					ReadMemory(
						GetFinalOffset(fp, out),
						out))
			}

			// return the stack pointer to its previous state
			prgrm.StackPointer = call.FramePointer
			// we'll now execute the next command
			prgrm.CallStack[prgrm.CallCounter].Line++
			// calling the actual command
			// prgrm.CallStack[prgrm.CallCounter].ccall(prgrm)
		}
	} else {
		/*
		   continue with call operator's execution
		*/
		fn := call.Operator
		expr := fn.Expressions[call.Line]
		// if it's a native, then we just process the arguments with execNative
		if expr.Operator == nil {
			// then it's a declaration
			// wiping this declaration's memory (removing garbage)
			newCall := &prgrm.CallStack[prgrm.CallCounter]
			newFP := newCall.FramePointer
			for c := 0; c < expr.Outputs[0].Size; c++ {
				prgrm.Memory[newFP+expr.Outputs[0].Offset+c] = 0
			}
			call.Line++
		} else if expr.Operator.IsNative {
			execNative(prgrm)
			call.Line++
		} else {
			/*
			   It was not a native, so we need to create another call
			   with the current expression's operator
			*/
			// we're going to use the next call in the callstack
			prgrm.CallCounter++
			if prgrm.CallCounter >= CALLSTACK_SIZE {
				panic(STACK_OVERFLOW_ERROR)
			}
			newCall := &prgrm.CallStack[prgrm.CallCounter]
			// setting the new call
			newCall.Operator = expr.Operator
			newCall.Line = 0
			newCall.FramePointer = prgrm.StackPointer
			// the stack pointer is moved to create room for the next call
			// prgrm.MemoryPointer += fn.Size
			prgrm.StackPointer += newCall.Operator.Size

			// checking if enough memory in stack
			if prgrm.StackPointer > STACK_SIZE {
				panic(STACK_OVERFLOW_ERROR)
			}

			fp := call.FramePointer
			newFP := newCall.FramePointer

			// wiping next stack frame (removing garbage)
			for c := 0; c < expr.Operator.Size; c++ {
				prgrm.Memory[newFP+c] = 0
			}

			for i, inp := range expr.Inputs {
				var byts []byte
				// finalOffset := inp.Offset
				finalOffset := GetFinalOffset(fp, inp)
				// finalOffset := fp + inp.Offset

				// if inp.Indexes != nil {
				// 	finalOffset = GetFinalOffset(&prgrm.Stacks[0], fp, inp)
				// }
				if inp.PassBy == PASSBY_REFERENCE {
					// If we're referencing an inner element, like an element of a slice (&slc[0])
					// or a field of a struct (&struct.fld) we no longer need to add
					// the OBJECT_HEADER_SIZE to the offset
					if inp.IsInnerReference {
						finalOffset -= OBJECT_HEADER_SIZE
					}
					byts = encoder.Serialize(int32(finalOffset))
				} else {
					byts = prgrm.Memory[finalOffset : finalOffset+inp.TotalSize]
				}

				// writing inputs to new stack frame
				WriteMemory(
					GetFinalOffset(newFP, newCall.Operator.Inputs[i]),
					// newFP + newCall.Operator.Inputs[i].Offset,
					// GetFinalOffset(prgrm.Memory, newFP, newCall.Operator.Inputs[i], MEM_WRITE),
					byts)
			}
		}
	}
	return nil
}
