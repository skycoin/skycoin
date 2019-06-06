package base

// import (
// 	"fmt"
// 	"regexp"
// 	"bytes"
// 	"github.com/skycoin/skycoin/src/cipher/encoder"
// )

// // This method removes 1 or more expressions
// // Then adds the equivalent of removed expressions

// func (expr *CXExpression) BuildExpression (outNames []*CXArgument) *CXExpression {
// 	numInps := len(expr.Operator.Inputs)
// 	numOuts := len(expr.Operator.Outputs)

// 	for i := 0; i < numInps; i++ {
// 		affs := FilterAffordances(expr.GetAffordances(nil), "Argument")
// 		r := random(0, len(affs))
// 		affs[r].ApplyAffordance()
// 	}

// 	if len(outNames) == 0 {
// 		for i := 0; i < numOuts; i++ {
// 			affs := FilterAffordances(expr.GetAffordances(nil), "OutputName")
// 			r := random(0, len(affs))
// 			affs[r].ApplyAffordance()
// 		}
// 	} else {
// 		expr.OutputNames = outNames
// 	}

// 	return expr
// }

// func (cxt *CXProgram) MutateSolution (solutionName string, fnBag string, numberExprs int) {
// 	cxt.SelectFunction(solutionName)
// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
// 		removeIndex := random(0, len(fn.Expressions))

// 		removedVars := fn.Expressions[removeIndex].OutputNames

// 		indexesToRemove := make([]int, 0)

// 		for i, expr := range fn.Expressions {
// 			for _, arg := range expr.Arguments {
// 				if i == removeIndex {
// 					indexesToRemove = append([]int{i}, indexesToRemove...)
// 					break
// 				}

// 				argIdentRemoved := false
// 				for _, removedVar := range removedVars {
// 					var identName string
// 					encoder.DeserializeRaw(*arg.Value, &identName)
// 					if i > removeIndex && identName == removedVar.Name {
// 						indexesToRemove = append([]int{i}, indexesToRemove...)
// 						argIdentRemoved = true
// 					}
// 				}
// 				if argIdentRemoved {
// 					break
// 				}

// 				for _, arg := range fn.Expressions[i].Arguments {
// 					broke := false
// 					for _, j := range indexesToRemove {
// 						argIsOutput := false
// 						for _, outName := range fn.Expressions[j].OutputNames {
// 							var identName string
// 							encoder.DeserializeRaw(*arg.Value, &identName)
// 							if identName == outName.Name {
// 								indexesToRemove = append([]int{i}, indexesToRemove...)
// 								broke = true
// 								argIsOutput = true
// 							}
// 						}
// 						if argIsOutput {
// 							break
// 						}
// 					}
// 					if broke {
// 						break
// 					}
// 				}
// 			}
// 		}

// 		indexesToRemove = removeDuplicatesInt(indexesToRemove)

// 		for _, i := range indexesToRemove {
// 			if i != len(fn.Expressions) {
// 				fn.Expressions = append(fn.Expressions[:i], fn.Expressions[i+1:]...)
// 			} else {
// 				fn.Expressions = fn.Expressions[:i]
// 			}
// 		}

// 		// re-indexing expression line numbers
// 		for i, expr := range fn.Expressions {
// 			expr.Line = i
// 		}

// 		hasOutputs := false
// 		fnOutputs := fn.Outputs
// 		for _, expr := range fn.Expressions {
// 			exprOutNames := expr.OutputNames

// 			if len(exprOutNames) != len(fnOutputs) {
// 				break
// 			}

// 			sameOutputs := true

// 			for i, out := range fnOutputs {
// 				if out.Name != exprOutNames[i].Name {
// 					sameOutputs = false
// 				}
// 			}

// 			if sameOutputs {
// 				hasOutputs = true
// 				break
// 			}
// 		}

// 		lenFnExprs := len(fn.Expressions)
// 		for i := 0; len(fn.Expressions) < numberExprs; i++ {
// 			var affs []*CXAffordance
// 			if !hasOutputs && i == numberExprs - lenFnExprs - 1 {

// 				var outs bytes.Buffer
// 				for j, out := range fn.Outputs {
// 					if j == len(fn.Outputs) - 1 {
// 						outs.WriteString(concat(out.Typ))
// 					} else {
// 						outs.WriteString(concat(out.Typ, ", "))
// 					}
// 				}

// 				affs = FilterAffordances(fn.GetAffordances(), "Expression", concat("\\(", regexp.QuoteMeta(outs.String()), "\\)$"), fnBag)
// 			} else {
// 				affs = FilterAffordances(fn.GetAffordances(), "Expression", fnBag)
// 			}

// 			r := random(0, len(affs))
// 			affs[r].ApplyAffordance()

// 			if expr, err := fn.GetCurrentExpression(); err == nil {
// 				if !hasOutputs && i == numberExprs - lenFnExprs - 1 {
// 					outNames := make([]*CXDefinition, len(fn.Outputs))
// 					for i, out := range fn.Outputs {
// 						outDef := MakeDefinition(
// 							out.Name,
// 							MakeDefaultValue(out.Typ),
// 							out.Typ)
// 						outNames[i] = outDef
// 					}
// 					expr.BuildExpression(outNames)
// 				} else {
// 					expr.BuildExpression(nil)
// 				}
// 			}
// 		}
// 	}
// }

// func (cxt *CXProgram) adaptPreEvolution (solutionName string) {
// 	if mod, err := cxt.GetModule("main"); err == nil {
// 		if _, err := cxt.GetFunction("main", "main"); err == nil {
// 			idx := -1
// 			for i, fn := range mod.Functions {
// 				if fn.Name == "main" {
// 					idx = i
// 					break
// 				}
// 			}
// 			mod.Functions = append(mod.Functions[:idx], mod.Functions[idx+1:]...)
// 		} else {
// 			panic(err)
// 		}

// 		if solMod, err := cxt.GetCurrentModule(); err == nil {
// 			if mod, err := cxt.GetModule("main"); err == nil {
// 				mod.AddFunction(MakeFunction("main"))
// 				if fn, err := cxt.GetCurrentFunction(); err == nil {
// 					fn.AddOutput(MakeParameter("out", "f64"))
// 					if sol, err := cxt.GetFunction(solutionName, solMod.Name); err == nil {
// 						fn.AddExpression(MakeExpression(sol))

// 						if expr, err := cxt.GetCurrentExpression(); err == nil {
// 							expr.AddOutputName("out")
// 							tmpVal := encoder.Serialize(float64(0))
// 							expr.AddArgument(MakeArgument(&tmpVal, "f64"))
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func (cxt *CXProgram) adaptInput (testValue float64) {
// 	if fn, err := cxt.GetFunction("main", "main"); err == nil {
// 		val := encoder.Serialize(testValue)
// 		arg := MakeArgument(&val, "f64")

// 		fn.Expressions[0].Arguments[0] = arg
// 	}
// }

// func (fromCxt *CXProgram) transferSolution (solutionName string, toCxt *CXProgram) {
// 	if fromMod, err := fromCxt.GetCurrentModule(); err == nil {
// 		if fromFn, err := fromCxt.GetFunction(solutionName, fromMod.Name); err == nil {
// 			if toMod, err := toCxt.GetCurrentModule(); err == nil {
// 				for _, fn := range toMod.Functions {
// 					if fn.Name == solutionName {
// 						//idx = i
// 						fn.Expressions = fromFn.Expressions
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func (cxt *CXProgram) Evolve (solutionName string, fnBag string, inputs, outputs []float64, numberExprs, iterations int, epsilon float64, expr *CXExpression, call *CXCall) error {
// 	cxt.SelectFunction(solutionName)

// 	// Initializing expressions
// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
// 		preExistingExpressions := len(fn.Expressions)
// 		for i := 0; i < numberExprs - preExistingExpressions; i++ {

// 			if i != numberExprs - preExistingExpressions - 1 {
// 				affs := FilterAffordances(fn.GetAffordances(), "Expression", fnBag)

// 				r := random(0, len(affs))
// 				affs[r].ApplyAffordance()

// 				if expr, err := fn.GetCurrentExpression(); err == nil {
// 					expr.BuildExpression(nil)
// 				}
// 			} else {

// 				fnOutTypNames := make([]string, 0)
// 				for _, out := range fn.Outputs {
// 					fnOutTypNames = append(fnOutTypNames, out.Typ)
// 				}

// 				var possibleOps string
// 				coreModule, _ := cxt.GetModule("core")
// 				for _, op := range coreModule.Functions {
// 					possibleOp := true
// 					for i, out := range op.Outputs {
// 						if len(op.Outputs) != len(fnOutTypNames) || out.Typ != fnOutTypNames[i] {
// 							possibleOp = false
// 							break
// 						}
// 					}

// 					if possibleOp {
// 						possibleOps = concat(possibleOps, concat(regexp.QuoteMeta(op.Name), "|"))
// 					}
// 				}
// 				for _, op := range fn.Module.Functions {
// 					possibleOp := true
// 					for i, out := range op.Outputs {
// 						if len(op.Outputs) != len(fnOutTypNames) || out.Typ != fnOutTypNames[i] {
// 							possibleOp = false
// 							break
// 						}
// 					}

// 					if possibleOp {
// 						possibleOps = concat(possibleOps, concat(regexp.QuoteMeta(op.Name), "|"))
// 					}
// 				}

// 				possibleOps = string([]byte(possibleOps)[:len(possibleOps)-1])

// 				affs := FilterAffordances(fn.GetAffordances(), "Expression", possibleOps, fnBag)

// 				r := random(0, len(affs))
// 				affs[r].ApplyAffordance()

// 				//making sure last expression assigns to output
// 				outNames := make([]*CXDefinition, len(fn.Outputs))
// 				for i, out := range fn.Outputs {
// 					outDef := MakeDefinition(
// 						out.Name,
// 						MakeDefaultValue(out.Typ),
// 						out.Typ)
// 					outNames[i] = outDef
// 				}

// 				if expr, err := fn.GetCurrentExpression(); err == nil {
// 					expr.BuildExpression(outNames)
// 				} else {
// 					fmt.Println(err)
// 				}
// 			}
// 		}
// 	}

// 	cxtCopy := MakeContextCopy(cxt, -1)

// 	best := cxtCopy
// 	best.adaptPreEvolution(solutionName)

// 	//best.PrintProgram(false)

// 	var finalError float64

// 	if mod, err := cxt.GetCurrentModule(); err == nil {
// 		if fn, err := cxt.GetFunction(solutionName, mod.Name); err == nil {
// 			if len(fn.Inputs) > 0 && len(fn.Outputs) > 0 {
// 				//fmt.Printf("Evolving function '%s'\n", solutionName)
// 				for i := 0; i < iterations; i++ {
// 					programs := make([]*CXProgram, 5) // it's always 5, because of the 1+4 strategy
// 					errors := make([]float64, 5)
// 					programs[0] = best // the best solution will always be at index 0

// 					for i := 1; i < 5; i++ {
// 						programs[i] = MakeContextCopy(programs[0], 0)
// 						// we need to mutate these 4
// 						programs[i].MutateSolution(solutionName, fnBag, numberExprs)
// 						//programs[i].PrintProgram(false)
// 					}

// 					// we need to evaluate all of them with each of the inputs and get the error
// 					for i, program := range programs {
// 						//fmt.Printf("Program:%d\n", i)
// 						var error float64 = 0
// 						for i, inp := range inputs {
// 							program.adaptInput(inp)
// 							//fmt.Println(program.Modules[0].Functions[0].Expressions[0].Arguments[0].Value)

// 							//program.PrintProgram(false)
// 							//fmt.Println(inp)

// 							program.Reset()
// 							program.Run(false, -1)

// 							// getting the simulated output
// 							var result float64
// 							// We'll always take the first value
// 							// The algorithm shouldn't work with multiple value returns yet

// 							if len(program.Outputs) < 1 {
// 								continue
// 							}
// 							output := program.Outputs[0].Value
// 							encoder.DeserializeRaw(*output, &result)

// 							//fmt.Println(program.Outputs[0].Value)

// 							// I don't want to import Math, so I will hardcode abs
// 							diff := float64(result - outputs[i])

// 							// fmt.Println(result)
// 							//fmt.Println(program.Outputs[0].Value)
// 							// fmt.Println(outputs[i])
// 							// fmt.Println()

// 							if diff >= 0 {
// 								error += diff
// 							} else {
// 								error += (diff * float64(-1))
// 							}
// 							//fmt.Println(error)

// 						}

// 						//fmt.Println(len(program.CallStack[len(program.CallStack) - 1].State))
// 						//fmt.Println(error / float64(len(inputs)))
// 						errors[i] = error / float64(len(inputs))
// 					}

// 					//fmt.Println(errors)

// 					// the program with the lowest error becomes the best
// 					bestIndex := 0
// 					for i, _ := range programs {
// 						if errors[i] <= errors[bestIndex] && errors[i] >= 0 {
// 							bestIndex = i
// 						}
// 					}
// 					//fmt.Println(errors)
// 					//fmt.Println(bestIndex)

// 					// print error each iteration
// 					//fmt.Println(errors[bestIndex])
// 					//fmt.Println(errors)

// 					//best.PrintProgram(false)
// 					best = programs[bestIndex]
// 					//best.PrintProgram(false)
// 					finalError = errors[bestIndex]

// 					if errors[bestIndex] < epsilon {
// 						break
// 					}
// 				}
// 				//fmt.Printf("Finished evolving function '%s'\n", solutionName)
// 			}
// 		}
// 	}

// 	best.transferSolution(solutionName, cxt)

// 	sFinalError := encoder.Serialize(finalError)

// 	for _, def := range call.State {
// 		if def.Name == expr.OutputNames[0].Name {
// 			def.Value = &sFinalError
// 			return nil
// 		}
// 	}

// 	call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sFinalError, "f64"))
// 	return nil
// }
