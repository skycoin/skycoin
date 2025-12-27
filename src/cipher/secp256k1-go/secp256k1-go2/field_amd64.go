//go:build !purego

package secp256k1go

// hasAsm indicates assembly implementations are available
const hasAsm = true

// fieldMulAsm computes r = a * b mod p in assembly
// This is a direct port of the Go Field.Mul implementation
//
//go:noescape
func fieldMulAsm(r, a, b *Field)

// fieldSqrAsm computes r = a * a mod p in assembly
// This is a direct port of the Go Field.Sqr implementation
//
//go:noescape
func fieldSqrAsm(r, a *Field)
