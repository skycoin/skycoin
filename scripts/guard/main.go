// Command guard checks the Naming, interface and init rules of AGENTS.md.
// Existing violations are listed in an allowlist that may only shrink.
//
//	go run ./scripts/guard names|ifaces|inits
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

type finding struct {
	key string
	why string
}

type guard struct {
	allowlist string
	rule      string
	scan      func(f *ast.File, fset *token.FileSet, file string) []finding
}

var guards = map[string]guard{
	"names": {
		allowlist: "names-allowlist.txt",
		rule:      `the "Naming" section of AGENTS.md`,
		scan:      scanNames,
	},
	"ifaces": {
		allowlist: "iface-allowlist.txt",
		rule:      "AGENTS.md, list the methods instead of embedding",
		scan:      scanIfaces,
	},
	"inits": {
		allowlist: "init-allowlist.txt",
		rule:      "AGENTS.md, call an explicit setup function from main instead",
		scan:      scanInits,
	},
}

func main() {
	dir := flag.String("dir", "scripts/guard", "directory holding the allowlists")
	list := flag.Bool("list", false, "print the current violations instead of checking")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: guard [-dir d] [-list] names|ifaces|inits")
		flag.PrintDefaults()
	}
	flag.Parse()
	g, ok := guards[flag.Arg(0)]
	if flag.NArg() != 1 || !ok {
		flag.Usage()
		os.Exit(2)
	}
	found, err := scanRepo(g.scan)
	if err != nil {
		fmt.Fprintln(os.Stderr, "guard:", err)
		os.Exit(1)
	}
	keys := make([]string, 0, len(found))
	for k := range found {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if *list {
		for _, k := range keys {
			fmt.Println(k)
		}
		return
	}
	allowPath := path.Join(*dir, g.allowlist)
	allowed, err := readLines(allowPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "guard:", err)
		os.Exit(1)
	}
	failed := false
	for _, l := range subtract(keys, allowed) {
		fmt.Printf("%s: %s, see %s\n", l, found[l], g.rule)
		failed = true
	}
	for _, l := range subtract(allowed, keys) {
		fmt.Printf("%s: fixed, remove it from %s\n", l, allowPath)
		failed = true
	}
	if !sort.StringsAreSorted(allowed) {
		fmt.Printf("%s is not sorted, run LC_ALL=C sort on it\n", allowPath)
		failed = true
	}
	fmt.Printf("%s: %d allowlisted\n", flag.Arg(0), len(allowed))
	if failed {
		os.Exit(1)
	}
}

// scanRepo runs scan over every Go file git would commit, skipping tests,
// generated code, vendor, node_modules, testdata and dot directories.
func scanRepo(scan func(*ast.File, *token.FileSet, string) []finding) (map[string]string, error) {
	out, err := exec.Command("git", "ls-files", "-z", "--cached", "--others", "--exclude-standard", "--", "*.go").Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w", err)
	}
	seen := map[string]string{}
	for _, file := range strings.Split(string(out), "\x00") {
		if skipFile(file) {
			continue
		}
		src, err := os.ReadFile(file) //nolint:gosec // paths come from git ls-files
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, src, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		if ast.IsGenerated(f) {
			continue
		}
		for _, fd := range scan(f, fset, file) {
			seen[fd.key] = fd.why
		}
	}
	return seen, nil
}

func skipFile(file string) bool {
	if file == "" || strings.HasSuffix(file, "_test.go") {
		return true
	}
	for _, d := range strings.Split(path.Dir(file), "/") {
		if d == "vendor" || d == "node_modules" || d == "testdata" || (strings.HasPrefix(d, ".") && d != ".") {
			return true
		}
	}
	return false
}

func readLines(file string) ([]string, error) {
	b, err := os.ReadFile(file) //nolint:gosec // allowlist path is a flag
	if err != nil {
		return nil, err
	}
	var lines []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" {
			lines = append(lines, l)
		}
	}
	return lines, sc.Err()
}

// subtract returns the lines of a that are not in b.
func subtract(a, b []string) []string {
	in := make(map[string]bool, len(b))
	for _, l := range b {
		in[l] = true
	}
	var out []string
	for _, l := range a {
		if !in[l] {
			out = append(out, l)
		}
	}
	return out
}

var bannedNames = map[string]bool{
	"Data": true, "Info": true, "Tmp": true, "Temp": true, "Item": true,
	"Thing": true, "Stuff": true, "Misc": true, "Do": true, "Handle": true,
	"Handler": true, "Process": true, "Manager": true, "Helper": true,
	"Util": true, "Get": true,
}

var bannedPackages = map[string]bool{
	"util": true, "utils": true, "common": true, "helpers": true, "helper": true, "misc": true,
}

// scanNames reports exported funcs, methods and types named only by a
// generic word, generic package names, and bool fields, params and package
// level vars with an is or has prefix. Lines are "file:Name".
func scanNames(f *ast.File, _ *token.FileSet, file string) []finding {
	var out []finding
	add := func(name, why string) { out = append(out, finding{file + ":" + name, why}) }
	dir := path.Dir(file)
	if bannedPackages[f.Name.Name] {
		out = append(out, finding{dir + ":package " + f.Name.Name, "generic package name"})
	} else if bannedPackages[path.Base(dir)] && dir != "." {
		out = append(out, finding{dir + ":package " + f.Name.Name, "generic directory name"})
	}
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			name := funcName(d)
			if bannedNames[d.Name.Name] {
				add(name, "generic func name")
			}
			boolParams(d.Type, name, add)
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if bannedNames[s.Name.Name] {
						add(s.Name.Name, "generic type name")
					}
				case *ast.ValueSpec:
					if d.Tok != token.VAR {
						continue
					}
					for i, n := range s.Names {
						if prefixed(n.Name) && boolValue(s, i) {
							add(n.Name, "bool with is/has prefix")
						}
					}
				}
			}
		}
	}
	// Struct fields and func literal params anywhere in the file, keyed by
	// the nearest enclosing declaration.
	owner := ""
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.FuncDecl:
			owner = funcName(t)
		case *ast.TypeSpec:
			owner = t.Name.Name
		case *ast.FuncLit:
			boolParams(t.Type, owner+".func", add)
		case *ast.InterfaceType:
			for _, m := range t.Methods.List {
				if ft, ok := m.Type.(*ast.FuncType); ok && len(m.Names) > 0 {
					boolParams(ft, owner+"."+m.Names[0].Name, add)
				}
			}
		case *ast.StructType:
			for _, fld := range t.Fields.List {
				if !boolType(fld.Type) {
					continue
				}
				for _, n := range fld.Names {
					if prefixed(n.Name) {
						add(owner+"."+n.Name, "bool with is/has prefix")
					}
				}
			}
		}
		return true
	})
	return out
}

func funcName(d *ast.FuncDecl) string {
	if d.Recv == nil || len(d.Recv.List) == 0 {
		return d.Name.Name
	}
	t := d.Recv.List[0].Type
	for {
		switch x := t.(type) {
		case *ast.StarExpr:
			t = x.X
			continue
		case *ast.IndexExpr:
			t = x.X
			continue
		case *ast.IndexListExpr:
			t = x.X
			continue
		case *ast.Ident:
			return x.Name + "." + d.Name.Name
		}
		return d.Name.Name
	}
}

func boolParams(ft *ast.FuncType, owner string, add func(string, string)) {
	for _, fl := range []*ast.FieldList{ft.Params, ft.Results} {
		if fl == nil {
			continue
		}
		for _, p := range fl.List {
			if !boolType(p.Type) {
				continue
			}
			for _, n := range p.Names {
				if prefixed(n.Name) {
					add(owner+"("+n.Name+")", "bool with is/has prefix")
				}
			}
		}
	}
}

// prefixed reports whether name is is or has followed by an upper case letter.
func prefixed(name string) bool {
	for _, p := range []string{"is", "Is", "has", "Has"} {
		if r := strings.TrimPrefix(name, p); r != name && r != "" && r[0] >= 'A' && r[0] <= 'Z' {
			return true
		}
	}
	return false
}

func boolType(e ast.Expr) bool {
	if s, ok := e.(*ast.StarExpr); ok {
		e = s.X
	}
	id, ok := e.(*ast.Ident)
	return ok && id.Name == "bool"
}

// boolValue reports whether the i'th name of a var spec is a bool, by its
// declared type or, without one, by a literal or boolean expression value.
func boolValue(s *ast.ValueSpec, i int) bool {
	if s.Type != nil {
		return boolType(s.Type)
	}
	if i >= len(s.Values) {
		return false
	}
	switch v := s.Values[i].(type) {
	case *ast.Ident:
		return v.Name == "true" || v.Name == "false"
	case *ast.UnaryExpr:
		return v.Op == token.NOT
	case *ast.BinaryExpr:
		switch v.Op {
		case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.LAND, token.LOR:
			return true
		}
	}
	return false
}

// scanIfaces reports every embedded interface, as "file:Owner embeds X".
// Type set terms such as ~int or int | string are not interfaces.
func scanIfaces(f *ast.File, _ *token.FileSet, file string) []finding {
	var out []finding
	owner := "anonymous"
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			owner = t.Name.Name
		case *ast.FuncDecl:
			owner = funcName(t)
		case *ast.InterfaceType:
			for _, m := range t.Methods.List {
				if len(m.Names) > 0 || !namedType(m.Type) {
					continue
				}
				out = append(out, finding{file + ":" + owner + " embeds " + types.ExprString(m.Type), "interface embeds an interface"})
			}
		}
		return true
	})
	return out
}

func namedType(e ast.Expr) bool {
	switch t := e.(type) {
	case *ast.Ident:
		return !predeclared[t.Name]
	case *ast.SelectorExpr, *ast.InterfaceType:
		return true
	case *ast.IndexExpr:
		return namedType(t.X)
	case *ast.IndexListExpr:
		return namedType(t.X)
	}
	return false
}

var predeclared = map[string]bool{
	"any": true, "comparable": true,
	"bool": true, "byte": true, "rune": true, "string": true, "uintptr": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"float32": true, "float64": true, "complex64": true, "complex128": true,
}

// scanInits reports each file outside package main that has a func init().
func scanInits(f *ast.File, _ *token.FileSet, file string) []finding {
	if f.Name.Name == "main" {
		return nil
	}
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.FuncDecl); ok && d.Recv == nil && d.Name.Name == "init" {
			return []finding{{file, "func init() outside package main"}}
		}
	}
	return nil
}
