package base

import (
	"errors"
	// "fmt"
	. "github.com/satori/go.uuid" //nolint golint
)

// CXExpression is used represent a CX expression.
//
// All statements in CX are expressions, including for loops and other control
// flow.
//
type CXExpression struct {
	// Metadata
	ElementID UUID

	// Contents
	Inputs   []*CXArgument
	Outputs  []*CXArgument
	Label    string
	Operator *CXFunction
	Function *CXFunction
	Package  *CXPackage

	// debugging
	FileName string
	FileLine int

	// used for jmp statements
	ThenLines int
	ElseLines int

	// 1 = start new scope; -1 = end scope; 0 = just regular expression
	ScopeOperation int

	IsMethodCall    bool
	IsStructLiteral bool
	IsArrayLiteral  bool
	IsUndType       bool
	IsBreak         bool
	IsContinue      bool
}

// MakeExpression ...
func MakeExpression(op *CXFunction, fileName string, fileLine int) *CXExpression {
	return &CXExpression{
		ElementID: MakeElementID(),
		Operator:  op,
		FileLine:  fileLine,
		FileName:  fileName}
}

// ----------------------------------------------------------------
//                             Getters

// GetInputs ...
func (expr *CXExpression) GetInputs() ([]*CXArgument, error) {
	if expr.Inputs != nil {
		return expr.Inputs, nil
	}
	return nil, errors.New("expression has no arguments")

}

// ----------------------------------------------------------------
//                     Member handling

// AddInput ...
func (expr *CXExpression) AddInput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Inputs = append(expr.Inputs, param)
	if param.Package == nil {
		param.Package = expr.Package
	}
	return expr
}

// RemoveInput ...
func (expr *CXExpression) RemoveInput() {
	if len(expr.Inputs) > 0 {
		expr.Inputs = expr.Inputs[:len(expr.Inputs)-1]
	}
}

// AddOutput ...
func (expr *CXExpression) AddOutput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Outputs = append(expr.Outputs, param)
	param.Package = expr.Package
	return expr
}

// RemoveOutput ...
func (expr *CXExpression) RemoveOutput() {
	if len(expr.Outputs) > 0 {
		expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]
	}
}

// AddLabel ...
func (expr *CXExpression) AddLabel(lbl string) *CXExpression {
	expr.Label = lbl
	return expr
}
