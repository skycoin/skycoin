package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"sort"
	"testing"
)

func keys(t *testing.T, scan func(*ast.File, *token.FileSet, string) []finding, file, src string) []string {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, fd := range scan(f, fset, file) {
		out = append(out, fd.key)
	}
	sort.Strings(out)
	return out
}

func TestScanNames(t *testing.T) {
	src := `package x

type Manager struct {
	isOpen, done bool
	HasPeer      *bool
	Island       bool
	isCount      int
}

type Info struct{}
type RouteInfo struct{}

func (m *Manager) Do(isRetry bool) {}
func (m *Manager) Get() {}
func (m *Manager) IsOpen() bool { return m.isOpen }
func (m *Manager) Error() string { return "" }
func Handler() { _ = func(hasX bool) {} }
func Get[T any]() {}

var isReady = true
var hasDone = 1 > 0
var isCount = 3
var isTyped bool

type Peer interface{ Dial(isLocal bool) }
`
	got := keys(t, scanNames, "p/x.go", src)
	want := []string{
		"p/x.go:Get", "p/x.go:Handler", "p/x.go:Handler.func(hasX)", "p/x.go:Info",
		"p/x.go:Manager", "p/x.go:Manager.Do", "p/x.go:Manager.Do(isRetry)", "p/x.go:Manager.Get",
		"p/x.go:Manager.HasPeer", "p/x.go:Manager.isOpen", "p/x.go:Peer.Dial(isLocal)",
		"p/x.go:hasDone", "p/x.go:isReady", "p/x.go:isTyped",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %q\nwant %q", got, want)
	}
}

func TestScanNamesPackage(t *testing.T) {
	if got := keys(t, scanNames, "pkg/util/a.go", "package util\n"); !reflect.DeepEqual(got, []string{"pkg/util:package util"}) {
		t.Fatalf("got %q", got)
	}
	if got := keys(t, scanNames, "pkg/common/a.go", "package shared\n"); !reflect.DeepEqual(got, []string{"pkg/common:package shared"}) {
		t.Fatalf("got %q", got)
	}
	if got := keys(t, scanNames, "pkg/netutil/a.go", "package netutil\n"); got != nil {
		t.Fatalf("got %q", got)
	}
}

func TestScanIfaces(t *testing.T) {
	src := `package x

import "io"

type A interface{ Foo() }
type B interface {
	A
	io.Reader
	Bar() io.Writer
}
type C interface{ ~int | string }
type D interface{ comparable }

func f() { var _ interface{ error } }
`
	got := keys(t, scanIfaces, "x.go", src)
	want := []string{"x.go:B embeds A", "x.go:B embeds io.Reader", "x.go:f embeds error"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestScanInits(t *testing.T) {
	if got := keys(t, scanInits, "a.go", "package a\nfunc init() {}\n"); !reflect.DeepEqual(got, []string{"a.go"}) {
		t.Fatalf("got %q", got)
	}
	if got := keys(t, scanInits, "m.go", "package main\nfunc init() {}\n"); got != nil {
		t.Fatalf("got %q", got)
	}
	if got := keys(t, scanInits, "b.go", "package b\ntype T struct{}\nfunc (T) init() {}\n"); got != nil {
		t.Fatalf("got %q", got)
	}
}
