//go:build purego || !amd64

package secp256k1go

// hasAsm indicates assembly implementations are NOT available
const hasAsm = false

// fieldMulAsm is a stub - not used on non-amd64 platforms
func fieldMulAsm(r, a, b *Field) {
	panic("fieldMulAsm should not be called on non-amd64 platforms")
}

// fieldSqrAsm is a stub - not used on non-amd64 platforms
func fieldSqrAsm(r, a *Field) {
	panic("fieldSqrAsm should not be called on non-amd64 platforms")
}
