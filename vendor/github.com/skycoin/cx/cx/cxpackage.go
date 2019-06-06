package base

import (
	"errors"
	"fmt"

	. "github.com/satori/go.uuid" //nolint golint
)

// CXPackage is used to represent a CX package.
//
type CXPackage struct {
	// Metadata
	Name      string // Name of the package
	ElementID UUID

	// Contents
	Imports   []*CXPackage  // imported packages
	Functions []*CXFunction // declared functions in this package
	Structs   []*CXStruct   // declared structs in this package
	Globals   []*CXArgument // declared global variables in this package

	// Used by the REPL and parser
	CurrentFunction *CXFunction
	CurrentStruct   *CXStruct
}

// MakePackage creates a new empty CXPackage.
//
// It can be filled in later with imports, structs, globals and functions.
//
func MakePackage(name string) *CXPackage {
	return &CXPackage{
		ElementID: MakeElementID(),
		Name:      name,
		Globals:   make([]*CXArgument, 0, 10),
		Imports:   make([]*CXPackage, 0),
		Structs:   make([]*CXStruct, 0),
		Functions: make([]*CXFunction, 0, 10),
	}
}

// ----------------------------------------------------------------
//                             Getters

// GetImport ...
func (pkg *CXPackage) GetImport(impName string) (*CXPackage, error) {
	for _, imp := range pkg.Imports {
		if imp.Name == impName {
			return imp, nil
		}
	}
	return nil, fmt.Errorf("package '%s' not imported", impName)
}

// GetFunctions ...
func (pkg *CXPackage) GetFunctions() ([]*CXFunction, error) {
	// going from map to slice
	if pkg.Functions != nil {
		return pkg.Functions, nil
	}
	return nil, fmt.Errorf("package '%s' has no functions", pkg.Name)

}

// GetFunction ...
func (pkg *CXPackage) GetFunction(fnName string) (*CXFunction, error) {
	var found bool
	for _, fn := range pkg.Functions {
		if fn.Name == fnName {
			return fn, nil
		}
	}

	// now checking in imported packages
	if !found {
		for _, imp := range pkg.Imports {
			for _, fn := range imp.Functions {
				if fn.Name == fnName {
					return fn, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("function '%s' not found in package '%s' or its imports", fnName, pkg.Name)
}

// GetMethod ...
func (pkg *CXPackage) GetMethod(fnName string, receiverType string) (*CXFunction, error) {
	for _, fn := range pkg.Functions {
		if fn.Name == fnName && len(fn.Inputs) > 0 && fn.Inputs[0].CustomType != nil && fn.Inputs[0].CustomType.Name == receiverType {
			return fn, nil
		}
	}

	return nil, fmt.Errorf("method '%s' not found in package '%s'", fnName, pkg.Name)
}

// GetStruct ...
func (pkg *CXPackage) GetStruct(strctName string) (*CXStruct, error) {
	var foundStrct *CXStruct
	for _, strct := range pkg.Structs {
		if strct.Name == strctName {
			foundStrct = strct
			break
		}
	}

	if foundStrct == nil {
		//looking in imports
		for _, imp := range pkg.Imports {
			for _, strct := range imp.Structs {
				if strct.Name == strctName {
					foundStrct = strct
					break
				}
			}
		}
	}

	if foundStrct != nil {
		return foundStrct, nil
	}
	return nil, fmt.Errorf("struct '%s' not found in package '%s'", strctName, pkg.Name)

}

// GetGlobal ...
func (pkg *CXPackage) GetGlobal(defName string) (*CXArgument, error) {
	var foundDef *CXArgument
	for _, def := range pkg.Globals {
		if def.Name == defName {
			foundDef = def
			break
		}
	}

	if foundDef != nil {
		return foundDef, nil
	}
	return nil, fmt.Errorf("global '%s' not found in package '%s'", defName, pkg.Name)

}

// GetCurrentFunction ...
func (pkg *CXPackage) GetCurrentFunction() (*CXFunction, error) {
	if pkg.CurrentFunction != nil {
		return pkg.CurrentFunction, nil
	}

	return nil, errors.New("current function is nil")
}

// GetCurrentStruct ...
func (pkg *CXPackage) GetCurrentStruct() (*CXStruct, error) {
	if pkg.CurrentStruct != nil {
		return pkg.CurrentStruct, nil
	}

	return nil, errors.New("current struct is nil")
}

// ----------------------------------------------------------------
//                     Member handling

// AddImport ...
func (pkg *CXPackage) AddImport(imp *CXPackage) *CXPackage {
	found := false
	for _, im := range pkg.Imports {
		if im.Name == imp.Name {
			found = true
			break
		}
	}
	if !found {
		pkg.Imports = append(pkg.Imports, imp)
	}

	return pkg
}

// RemoveImport ...
func (pkg *CXPackage) RemoveImport(impName string) {
	lenImps := len(pkg.Imports)
	for i, imp := range pkg.Imports {
		if imp.Name == impName {
			if i == lenImps-1 {
				pkg.Imports = pkg.Imports[:len(pkg.Imports)-1]
			} else {
				pkg.Imports = append(pkg.Imports[:i], pkg.Imports[i+1:]...)
			}
			break
		}
	}
}

// AddFunction ...
func (pkg *CXPackage) AddFunction(fn *CXFunction) *CXPackage {
	fn.Package = pkg

	found := false
	for i, f := range pkg.Functions {
		if f.Name == fn.Name {
			pkg.Functions[i].Name = fn.Name
			pkg.Functions[i].Inputs = fn.Inputs
			pkg.Functions[i].Outputs = fn.Outputs
			pkg.Functions[i].Expressions = fn.Expressions
			pkg.Functions[i].CurrentExpression = fn.CurrentExpression
			pkg.Functions[i].Package = fn.Package
			pkg.CurrentFunction = pkg.Functions[i]
			found = true
			break
		}
	}
	if found && !InREPL {
		println(CompilationError(fn.FileName, fn.FileLine), "function redeclaration")
	}
	if !found {
		pkg.Functions = append(pkg.Functions, fn)
		pkg.CurrentFunction = fn
	}

	return pkg
}

// RemoveFunction ...
func (pkg *CXPackage) RemoveFunction(fnName string) {
	lenFns := len(pkg.Functions)
	for i, fn := range pkg.Functions {
		if fn.Name == fnName {
			if i == lenFns-1 {
				pkg.Functions = pkg.Functions[:len(pkg.Functions)-1]
			} else {
				pkg.Functions = append(pkg.Functions[:i], pkg.Functions[i+1:]...)
			}
			break
		}
	}
}

// AddStruct ...
func (pkg *CXPackage) AddStruct(strct *CXStruct) *CXPackage {
	found := false
	for i, s := range pkg.Structs {
		if s.Name == strct.Name {
			pkg.Structs[i] = strct
			found = true
			break
		}
	}
	if !found {
		pkg.Structs = append(pkg.Structs, strct)
	}

	strct.Package = pkg
	pkg.CurrentStruct = strct

	return pkg
}

// RemoveStruct ...
func (pkg *CXPackage) RemoveStruct(strctName string) {
	lenStrcts := len(pkg.Structs)
	for i, strct := range pkg.Structs {
		if strct.Name == strctName {
			if i == lenStrcts-1 {
				pkg.Structs = pkg.Structs[:len(pkg.Structs)-1]
			} else {
				pkg.Structs = append(pkg.Structs[:i], pkg.Structs[i+1:]...)
			}
			break
		}
	}
}

// AddGlobal ...
func (pkg *CXPackage) AddGlobal(def *CXArgument) *CXPackage {
	def.Package = pkg
	found := false
	for i, df := range pkg.Globals {
		if df.Name == def.Name {
			pkg.Globals[i] = def
			found = true
			break
		}
	}
	if !found {
		pkg.Globals = append(pkg.Globals, def)
	}

	return pkg
}

// RemoveGlobal ...
func (pkg *CXPackage) RemoveGlobal(defName string) {
	lenGlobals := len(pkg.Globals)
	for i, def := range pkg.Globals {
		if def.Name == defName {
			if i == lenGlobals-1 {
				pkg.Globals = pkg.Globals[:len(pkg.Globals)-1]
			} else {
				pkg.Globals = append(pkg.Globals[:i], pkg.Globals[i+1:]...)
			}
			break
		}
	}
}

// ----------------------------------------------------------------
//                             Selectors

// SelectFunction ...
func (pkg *CXPackage) SelectFunction(name string) (*CXFunction, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {

	// 		if pkg, err := cxt.GetCurrentPackage(); err == nil {
	// 			pkg.SelectFunction(name)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, pkg.Context)

	var found *CXFunction
	for _, fn := range pkg.Functions {
		if fn.Name == name {
			pkg.CurrentFunction = fn
			found = fn
		}
	}

	if found == nil {
		return nil, fmt.Errorf("function '%s' does not exist", name)
	}

	return found, nil
}

// SelectStruct ...
func (pkg *CXPackage) SelectStruct(name string) (*CXStruct, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if pkg, err := cxt.GetCurrentPackage(); err == nil {
	// 			pkg.SelectStruct(name)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, pkg.Context)

	var found *CXStruct
	for _, strct := range pkg.Structs {
		if strct.Name == name {
			pkg.CurrentStruct = strct
			found = strct
		}
	}

	if found == nil {
		return nil, errors.New("Desired structure does not exist")
	}

	return found, nil
}

// SelectExpression ...
func (pkg *CXPackage) SelectExpression(line int) (*CXExpression, error) {
	// prgrmStep := &CXProgramStep{
	// 	Action: func(cxt *CXProgram) {
	// 		if pkg, err := cxt.GetCurrentPackage(); err == nil {
	// 			pkg.SelectExpression(line)
	// 		}
	// 	},
	// }
	// saveProgramStep(prgrmStep, pkg.Context)
	fn, err := pkg.GetCurrentFunction()
	if err == nil {
		return fn.SelectExpression(line)
	}
	return nil, err
}
