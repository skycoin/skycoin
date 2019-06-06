package base

import (
	"errors"
	"fmt"

	. "github.com/satori/go.uuid" //nolint golint
)

// CXFunction is used to represent a CX function.
//
type CXFunction struct {
	// Metadata
	Name      string     // Name of the function
	Package   *CXPackage // The package it's a member of
	ElementID UUID
	IsNative  bool // True if the function is native to CX, e.g. int32.add()
	OpCode    int  // opcode if IsNative = true

	// Contents
	Inputs      []*CXArgument   // Input parameters to the function
	Outputs     []*CXArgument   // Output parameters from the function
	Expressions []*CXExpression // Expressions, including control flow statements, in the function
	Length      int             // number of expressions, pre-computed for performance
	Size        int             // automatic memory size

	// Debugging
	FileName string
	FileLine int

	// Used by the GC
	ListOfPointers []*CXArgument // Root pointers for the GC algorithm

	// Used by the REPL and parser
	CurrentExpression *CXExpression
}

// MakeFunction creates an empty function.
//
// Later, parameters and contents can be added.
//
func MakeFunction(name string, fileName string, fileLine int) *CXFunction {
	return &CXFunction{
		ElementID: MakeElementID(),
		Name:      name,
		FileName:  fileName,
		FileLine:  fileLine,
	}
}

// MakeNativeFunction creates a native function such as i32.add()
//
func MakeNativeFunction(opCode int, inputs []*CXArgument, outputs []*CXArgument) *CXFunction {
	fn := &CXFunction{
		ElementID: MakeElementID(),
		IsNative:  true,
		OpCode:    opCode,
	}

	offset := 0
	for _, inp := range inputs {
		// for _, typCode := range inputs {
		// inp := MakeArgument("", "", -1).AddType(TypeNames[typCode])
		inp.Offset = offset
		offset += inp.Size
		fn.Inputs = append(fn.Inputs, inp)
	}
	for _, out := range outputs {
		// for _, typCode := range outputs {
		// fn.Outputs = append(fn.Outputs, MakeArgument("", "", -1).AddType(TypeNames[typCode]))
		// out := MakeArgument("", "", -1).AddType(TypeNames[typCode])
		fn.Outputs = append(fn.Outputs, out)
		out.Offset = offset
		offset += out.Size
	}

	return fn
}

// ----------------------------------------------------------------
//                             Getters

// GetExpressions ...
func (fn *CXFunction) GetExpressions() ([]*CXExpression, error) {
	if fn.Expressions != nil {
		return fn.Expressions, nil
	}
	return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)

}

// GetExpressionByLabel ...
func (fn *CXFunction) GetExpressionByLabel(lbl string) (*CXExpression, error) {
	if fn.Expressions != nil {
		for _, expr := range fn.Expressions {
			if expr.Label == lbl {
				return expr, nil
			}
		}

		return nil, fmt.Errorf("expression '%s' not found in function '%s'", lbl, fn.Name)
	}
	return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)

}

// GetExpressionByLine ...
func (fn *CXFunction) GetExpressionByLine(line int) (*CXExpression, error) {
	if fn.Expressions != nil {
		if line <= len(fn.Expressions) {
			return fn.Expressions[line], nil
		}
		return nil, fmt.Errorf("expression line number '%d' exceeds number of expressions in function '%s'", line, fn.Name)

	}
	return nil, fmt.Errorf("function '%s' has no expressions", fn.Name)

}

// GetCurrentExpression ...
func (fn *CXFunction) GetCurrentExpression() (*CXExpression, error) {
	if fn.CurrentExpression != nil {
		return fn.CurrentExpression, nil
	} else if fn.Expressions != nil {
		return fn.Expressions[0], nil
	} else {
		return nil, errors.New("current expression is nil")
	}
}

// ----------------------------------------------------------------
//                     Member handling

// ----------------------------------------------------------------
//                     Member handling

// AddInput ...
func (fn *CXFunction) AddInput(param *CXArgument) *CXFunction {
	found := false
	for _, inp := range fn.Inputs {
		if inp.Name == param.Name {
			found = true
			break
		}
	}
	if !found {
		fn.Inputs = append(fn.Inputs, param)
	}

	return fn
}

// RemoveInput ...
func (fn *CXFunction) RemoveInput(inpName string) {
	if len(fn.Inputs) > 0 {
		lenInps := len(fn.Inputs)
		for i, inp := range fn.Inputs {
			if inp.Name == inpName {
				if i == lenInps {
					fn.Inputs = fn.Inputs[:len(fn.Inputs)-1]
				} else {
					fn.Inputs = append(fn.Inputs[:i], fn.Inputs[i+1:]...)
				}
				break
			}
		}
	}
}

// AddOutput ...
func (fn *CXFunction) AddOutput(param *CXArgument) *CXFunction {
	found := false
	for _, out := range fn.Outputs {
		if out.Name == param.Name {
			found = true
			break
		}
	}
	if !found {
		fn.Outputs = append(fn.Outputs, param)
	}

	param.Package = fn.Package

	return fn
}

// RemoveOutput ...
func (fn *CXFunction) RemoveOutput(outName string) {
	if len(fn.Outputs) > 0 {
		lenOuts := len(fn.Outputs)
		for i, out := range fn.Outputs {
			if out.Name == outName {
				if i == lenOuts {
					fn.Outputs = fn.Outputs[:len(fn.Outputs)-1]
				} else {
					fn.Outputs = append(fn.Outputs[:i], fn.Outputs[i+1:]...)
				}
				break
			}
		}
	}
}

// AddExpression ...
func (fn *CXFunction) AddExpression(expr *CXExpression) *CXFunction {
	// expr.Program = fn.Program
	expr.Package = fn.Package
	expr.Function = fn
	fn.Expressions = append(fn.Expressions, expr)
	fn.CurrentExpression = expr
	fn.Length++
	return fn
}

// RemoveExpression ...
func (fn *CXFunction) RemoveExpression(line int) {
	if len(fn.Expressions) > 0 {
		lenExprs := len(fn.Expressions)
		if line >= lenExprs-1 || line < 0 {
			fn.Expressions = fn.Expressions[:len(fn.Expressions)-1]
		} else {
			fn.Expressions = append(fn.Expressions[:line], fn.Expressions[line+1:]...)
		}
		// for i, expr := range fn.Expressions {
		// 	expr.Index = i
		// }
	}
}

// ----------------------------------------------------------------
//                             Selectors

// SelectExpression ...
func (fn *CXFunction) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if mod, err := cxt.GetCurrentPackage(); err == nil {
	// 			if fn, err := mod.GetCurrentFunction(); err == nil {
	// 				fn.SelectExpression(line)
	// 			}
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, fn.Context)
	if len(fn.Expressions) == 0 {
		return nil, errors.New("There are no expressions in this function")
	}

	if line >= len(fn.Expressions) {
		line = len(fn.Expressions) - 1
	}

	if line < 0 {
		line = 0
	}

	expr := fn.Expressions[line]
	fn.CurrentExpression = expr

	return expr, nil
}
