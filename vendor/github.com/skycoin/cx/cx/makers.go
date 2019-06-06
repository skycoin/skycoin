package base

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Var
var (
	HeapOffset    int
	genSymCounter int
)

// MakeElementID ...
func MakeElementID() uuid.UUID {
	return uuid.NewV4()
}

// MakeGenSym ...
func MakeGenSym(name string) string {
	gensym := fmt.Sprintf("%s_%d", name, genSymCounter)
	genSymCounter++

	return gensym
}

// MakeDefaultValue Used only for native types
func MakeDefaultValue(typName string) *[]byte {
	var zeroVal []byte
	switch typName {
	case "byte":
		zeroVal = make([]byte, 1)
	case "i64", "f64":
		zeroVal = make([]byte, 8)
	default:
		zeroVal = make([]byte, 4)
	}
	return &zeroVal
}

// MakeValue ...
func MakeValue(value string) *[]byte {
	byts := encoder.Serialize(value)
	return &byts
}

// MakeCall ...
func MakeCall(op *CXFunction) CXCall {
	return CXCall{
		Operator:     op,
		Line:         0,
		FramePointer: 0,
		// Package:       pkg,
		// Program:       prgrm,
	}
}

// MakeIdentityOpName ...
func MakeIdentityOpName(typeName string) string {
	switch typeName {
	case "str":
		return "str.id"
	case "bool":
		return "bool.id"
	case "byte":
		return "byte.id"
	case "i32":
		return "i32.id"
	case "i64":
		return "i64.id"
	case "f32":
		return "f32.id"
	case "f64":
		return "f64.id"
	case "[]bool":
		return "[]bool.id"
	case "[]byte":
		return "[]byte.id"
	case "[]str":
		return "[]str.id"
	case "[]i32":
		return "[]i32.id"
	case "[]i64":
		return "[]i64.id"
	case "[]f32":
		return "[]f32.id"
	case "[]f64":
		return "[]f64.id"
	default:
		return ""
	}
}

// func MakeParameterCopy(param *CXArgument) *CXArgument {
// 	return &CXArgument{
// 		Name: param.Name,
// 		// Typ:  param.Typ,
// 	}
// }

// func MakeArgumentCopy(arg *CXArgument) *CXArgument {
// 	value := *arg.Value
// 	return &CXArgument{
// 		Typ:   arg.Typ,
// 		Value: &value,
// 	}
// }

// func MakeExpressionCopy(expr *CXExpression, fn *CXFunction, mod *CXPackage, cxt *CXProgram) *CXExpression {
// 	argsCopy := make([]*CXArgument, len(expr.Inputs))
// 	for i, arg := range expr.Inputs {
// 		argsCopy[i] = MakeArgumentCopy(arg)
// 	}
// 	return &CXExpression{
// 		Operator: expr.Operator,
// 		Inputs:   argsCopy,
// 		Outputs:  expr.Outputs,
// 		Line:     expr.Line,
// 		Function: fn,
// 		Package:  mod,
// 		Program:  cxt,
// 	}
// }

// func MakeFunctionCopy(fn *CXFunction, mod *CXPackage, cxt *CXProgram) *CXFunction {
// 	newFn := &CXFunction{}
// 	inputsCopy := make([]*CXArgument, len(fn.Inputs))
// 	outputsCopy := make([]*CXArgument, len(fn.Outputs))
// 	exprsCopy := make([]*CXExpression, len(fn.Expressions))
// 	for i, inp := range fn.Inputs {
// 		inputsCopy[i] = MakeParameterCopy(inp)
// 	}
// 	for i, out := range fn.Outputs {
// 		outputsCopy[i] = MakeParameterCopy(out)
// 	}

// 	for i, expr := range fn.Expressions {
// 		exprsCopy[i] = MakeExpressionCopy(expr, newFn, mod, cxt)
// 	}

// 	newFn.Name = fn.Name
// 	newFn.Inputs = inputsCopy
// 	newFn.Outputs = outputsCopy

// 	// if fn.Output != nil {
// 	// 	newFn.Output = MakeParameterCopy(fn.Output)
// 	// }
// 	newFn.Expressions = exprsCopy
// 	if len(exprsCopy) > 0 {
// 		newFn.CurrentExpression = exprsCopy[len(exprsCopy)-1]
// 	}
// 	newFn.Package = mod
// 	newFn.Program = cxt

// 	return newFn
// }

// func MakeFieldCopy(fld *CXArgument) *CXArgument {
// 	return &CXArgument{
// 		Name: fld.Name,
// 		Typ:  fld.Typ,
// 	}
// }

// func MakeStructCopy(strct *CXStruct, mod *CXPackage, cxt *CXProgram) *CXStruct {
// 	fldsCopy := make([]*CXArgument, len(strct.Fields))
// 	for i, fld := range strct.Fields {
// 		fldsCopy[i] = MakeFieldCopy(fld)
// 	}
// 	return &CXStruct{
// 		Name:    strct.Name,
// 		Fields:  fldsCopy,
// 		Package: mod,
// 		Program: cxt,
// 	}
// }

// func MakeDefinitionCopy(def *CXArgument, mod *CXPackage, cxt *CXProgram) *CXArgument {
// 	valCopy := *def.Value
// 	return &CXArgument{
// 		Name:    def.Name,
// 		Typ:     def.Typ,
// 		Value:   &valCopy,
// 		Package: mod,
// 		Program: cxt,
// 	}
// }

// func MakeModuleCopy(mod *CXPackage, cxt *CXProgram) *CXPackage {
// 	newMod := &CXPackage{Program: cxt}
// 	fnsCopy := make([]*CXFunction, len(mod.Functions))
// 	strctsCopy := make([]*CXStruct, len(mod.Structs))
// 	defsCopy := make([]*CXArgument, len(mod.Globals))

// 	for k, fn := range mod.Functions {
// 		fnsCopy[k] = MakeFunctionCopy(fn, newMod, cxt)
// 	}
// 	for k, strct := range mod.Structs {
// 		strctsCopy[k] = MakeStructCopy(strct, newMod, cxt)
// 	}
// 	for k, def := range mod.Globals {
// 		defsCopy[k] = MakeDefinitionCopy(def, newMod, cxt)
// 	}

// 	// Setting current function in copy
// 	for _, fn := range fnsCopy {
// 		if fn.Name == mod.CurrentFunction.Name {
// 			newMod.CurrentFunction = fn
// 		}
// 	}

// 	newMod.Name = mod.Name
// 	newMod.Imports = mod.Imports
// 	newMod.Functions = fnsCopy
// 	newMod.Structs = strctsCopy
// 	newMod.Globals = defsCopy
// 	newMod.Program = cxt

// 	return newMod
// }

// func MakeCallCopy(call *CXCall, mod *CXPackage, cxt *CXProgram) *CXCall {
// 	stateCopy := make([]*CXArgument, len(call.State))
// 	for k, v := range call.State {
// 		stateCopy[k] = MakeDefinitionCopy(v, mod, cxt)
// 	}
// 	return &CXCall{
// 		Operator:      call.Operator,
// 		Line:          call.Line,
// 		State:         stateCopy,
// 		ReturnAddress: call.ReturnAddress,
// 		Package:       mod,
// 		Program:       cxt,
// 	}
// }

// MakeCallStack ...
func MakeCallStack(size int) []CXCall {
	return make([]CXCall, 0)
	// return &CXCallStack{
	// 	Calls: make([]*CXCall, size),
	// }
}

// func MakeContextCopy(cxt *CXProgram, stepNumber int) *CXProgram {
// 	newContext := &CXProgram{}

// 	modsCopy := make([]*CXPackage, len(cxt.Packages))
// 	if stepNumber >= len(cxt.Steps) || stepNumber < 0 {
// 		stepNumber = len(cxt.Steps) - 1
// 	}

// 	for k, mod := range cxt.Packages {
// 		modsCopy[k] = MakeModuleCopy(mod, newContext)
// 	}

// 	// Setting current module in copy
// 	for _, mod := range modsCopy {
// 		if mod.Name == cxt.CurrentPackage.Name {
// 			newContext.CurrentPackage = mod
// 		}
// 	}

// 	newContext.Packages = modsCopy

// 	// Making imports copies
// 	for _, mod := range modsCopy {
// 		for impKey, _ := range mod.Imports {
// 			mod.Imports[impKey] = modsCopy[impKey]
// 		}
// 	}

// 	// Making expressions/operators
// 	for _, mod := range modsCopy {
// 		for _, fn := range mod.Functions {
// 			for _, expr := range fn.Expressions {
// 				if op, err := newContext.GetFunction(expr.Operator.Name, expr.Package.Name); err == nil {
// 					expr.Operator = op
// 				}
// 			}
// 		}
// 	}

// 	if len(cxt.Steps) > 0 {
// 		reqStep := cxt.Steps[stepNumber]
// 		newStep := MakeCallStack(len(reqStep))

// 		var lastCall *CXCall
// 		for j, call := range reqStep {
// 			var callModule *CXPackage
// 			for _, mod := range modsCopy {
// 				if call.Package.Name == mod.Name {
// 					callModule = mod
// 				}
// 			}

// 			newCall := *MakeCallCopy(&call, callModule, newContext)
// 			if callOp, err := newContext.GetFunction(call.Operator.Name, call.Operator.Package.Name); err == nil {
// 				newCall.Operator = callOp
// 			}
// 			newCall.ReturnAddress = lastCall
// 			lastCall = &newCall
// 			newStep[j] = newCall
// 		}

// 		newContext.CallStack = newStep
// 		newContext.Steps = nil
// 	}

// 	return newContext
// }
