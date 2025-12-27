//go:build !purego && amd64

#include "textflag.h"

// Helper macros for field operations
// MULADDQ multiplies two 32-bit limbs and adds to accumulator
#define MULADDQ(aoff, boff, acc) \
	MOVL aoff(SI), AX \
	MOVL boff(BX), CX \
	IMULQ CX, AX \
	ADDQ AX, acc

// SQRADDQ2 computes 2*a[i]*a[j] and adds to accumulator (for squaring)
#define SQRADDQ2(aoff, boff, acc) \
	MOVL aoff(SI), AX \
	MOVL boff(SI), CX \
	IMULQ CX, AX \
	ADDQ AX, acc \
	ADDQ AX, acc

// SQRADDQ computes a[i]*a[i] and adds to accumulator (diagonal terms)
#define SQRADDQ(aoff, acc) \
	MOVL aoff(SI), AX \
	IMULQ AX, AX \
	ADDQ AX, acc

// Field multiplication: r = a * b mod p
// func fieldMulAsm(r, a, b *Field)
//
// Field is [10]uint32 (40 bytes total)
// Each limb is 26 bits stored in uint32
//
TEXT ·fieldMulAsm(SB), NOSPLIT, $160-24
	MOVQ r+0(FP), DI     // DI = r (result)
	MOVQ a+8(FP), SI     // SI = a (operand 1)
	MOVQ b+16(FP), BX    // BX = b (operand 2)
	
	// Register allocation:
	// R8  = c (carry/accumulator)
	// R9  = temp for storing masked limb
	// R10-R15 = scratch
	// AX, CX, DX = used by IMULQ
	
	// ===== Phase 1: Compute cross-products t0-t19 =====
	
	// t0 = a[0] * b[0]
	MOVL (SI), AX
	MOVL (BX), CX
	IMULQ CX, AX
	MOVQ AX, R8           // c = a[0] * b[0]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9   // t0 = c & 0x3FFFFFF
	MOVQ R9, 0(SP)
	SHRQ $26, R8          // c >>= 26
	
	// t1 = c + a[0]*b[1] + a[1]*b[0]
	MULADDQ(0, 4, R8)     // c += a[0] * b[1]
	MULADDQ(4, 0, R8)     // c += a[1] * b[0]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 8(SP)
	SHRQ $26, R8
	
	// t2 = c + a[0]*b[2] + a[1]*b[1] + a[2]*b[0]
	MULADDQ(0, 8, R8)
	MULADDQ(4, 4, R8)
	MULADDQ(8, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 16(SP)
	SHRQ $26, R8
	
	// t3 = c + a[0]*b[3] + a[1]*b[2] + a[2]*b[1] + a[3]*b[0]
	MULADDQ(0, 12, R8)
	MULADDQ(4, 8, R8)
	MULADDQ(8, 4, R8)
	MULADDQ(12, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 24(SP)
	SHRQ $26, R8
	
	// t4 = c + a[0]*b[4] + a[1]*b[3] + a[2]*b[2] + a[3]*b[1] + a[4]*b[0]
	MULADDQ(0, 16, R8)
	MULADDQ(4, 12, R8)
	MULADDQ(8, 8, R8)
	MULADDQ(12, 4, R8)
	MULADDQ(16, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 32(SP)
	SHRQ $26, R8
	
	// t5 = c + a[0]*b[5] + ... + a[5]*b[0]
	MULADDQ(0, 20, R8)
	MULADDQ(4, 16, R8)
	MULADDQ(8, 12, R8)
	MULADDQ(12, 8, R8)
	MULADDQ(16, 4, R8)
	MULADDQ(20, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 40(SP)
	SHRQ $26, R8
	
	// t6 = c + a[0]*b[6] + ... + a[6]*b[0]
	MULADDQ(0, 24, R8)
	MULADDQ(4, 20, R8)
	MULADDQ(8, 16, R8)
	MULADDQ(12, 12, R8)
	MULADDQ(16, 8, R8)
	MULADDQ(20, 4, R8)
	MULADDQ(24, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 48(SP)
	SHRQ $26, R8
	
	// t7 = c + a[0]*b[7] + ... + a[7]*b[0]
	MULADDQ(0, 28, R8)
	MULADDQ(4, 24, R8)
	MULADDQ(8, 20, R8)
	MULADDQ(12, 16, R8)
	MULADDQ(16, 12, R8)
	MULADDQ(20, 8, R8)
	MULADDQ(24, 4, R8)
	MULADDQ(28, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 56(SP)
	SHRQ $26, R8
	
	// t8 = c + a[0]*b[8] + ... + a[8]*b[0]
	MULADDQ(0, 32, R8)
	MULADDQ(4, 28, R8)
	MULADDQ(8, 24, R8)
	MULADDQ(12, 20, R8)
	MULADDQ(16, 16, R8)
	MULADDQ(20, 12, R8)
	MULADDQ(24, 8, R8)
	MULADDQ(28, 4, R8)
	MULADDQ(32, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 64(SP)
	SHRQ $26, R8
	
	// t9 = c + a[0]*b[9] + ... + a[9]*b[0]
	MULADDQ(0, 36, R8)
	MULADDQ(4, 32, R8)
	MULADDQ(8, 28, R8)
	MULADDQ(12, 24, R8)
	MULADDQ(16, 20, R8)
	MULADDQ(20, 16, R8)
	MULADDQ(24, 12, R8)
	MULADDQ(28, 8, R8)
	MULADDQ(32, 4, R8)
	MULADDQ(36, 0, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 72(SP)
	SHRQ $26, R8
	
	// t10 = c + a[1]*b[9] + ... + a[9]*b[1]
	MULADDQ(4, 36, R8)
	MULADDQ(8, 32, R8)
	MULADDQ(12, 28, R8)
	MULADDQ(16, 24, R8)
	MULADDQ(20, 20, R8)
	MULADDQ(24, 16, R8)
	MULADDQ(28, 12, R8)
	MULADDQ(32, 8, R8)
	MULADDQ(36, 4, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 80(SP)
	SHRQ $26, R8
	
	// t11 = c + a[2]*b[9] + ... + a[9]*b[2]
	MULADDQ(8, 36, R8)
	MULADDQ(12, 32, R8)
	MULADDQ(16, 28, R8)
	MULADDQ(20, 24, R8)
	MULADDQ(24, 20, R8)
	MULADDQ(28, 16, R8)
	MULADDQ(32, 12, R8)
	MULADDQ(36, 8, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 88(SP)
	SHRQ $26, R8
	
	// t12 = c + a[3]*b[9] + ... + a[9]*b[3]
	MULADDQ(12, 36, R8)
	MULADDQ(16, 32, R8)
	MULADDQ(20, 28, R8)
	MULADDQ(24, 24, R8)
	MULADDQ(28, 20, R8)
	MULADDQ(32, 16, R8)
	MULADDQ(36, 12, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 96(SP)
	SHRQ $26, R8
	
	// t13 = c + a[4]*b[9] + ... + a[9]*b[4]
	MULADDQ(16, 36, R8)
	MULADDQ(20, 32, R8)
	MULADDQ(24, 28, R8)
	MULADDQ(28, 24, R8)
	MULADDQ(32, 20, R8)
	MULADDQ(36, 16, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 104(SP)
	SHRQ $26, R8
	
	// t14 = c + a[5]*b[9] + ... + a[9]*b[5]
	MULADDQ(20, 36, R8)
	MULADDQ(24, 32, R8)
	MULADDQ(28, 28, R8)
	MULADDQ(32, 24, R8)
	MULADDQ(36, 20, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 112(SP)
	SHRQ $26, R8
	
	// t15 = c + a[6]*b[9] + ... + a[9]*b[6]
	MULADDQ(24, 36, R8)
	MULADDQ(28, 32, R8)
	MULADDQ(32, 28, R8)
	MULADDQ(36, 24, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 120(SP)
	SHRQ $26, R8
	
	// t16 = c + a[7]*b[9] + a[8]*b[8] + a[9]*b[7]
	MULADDQ(28, 36, R8)
	MULADDQ(32, 32, R8)
	MULADDQ(36, 28, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 128(SP)
	SHRQ $26, R8
	
	// t17 = c + a[8]*b[9] + a[9]*b[8]
	MULADDQ(32, 36, R8)
	MULADDQ(36, 32, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 136(SP)
	SHRQ $26, R8
	
	// t18 = c + a[9]*b[9]
	MULADDQ(36, 36, R8)
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 144(SP)
	SHRQ $26, R8
	
	// t19 = c
	MOVQ R8, 152(SP)
	
	// ===== Phase 2: Modular Reduction =====
	// Reduce using secp256k1 modulus: p = 2^256 - 2^32 - 977
	// Constants: 0x3D10, 0x400, 0x1000003D10, 0x3D1, 0x40
	
	// c = t0 + t10*0x3D10
	MOVQ 0(SP), R8        // R8 = t0
	MOVQ 80(SP), AX       // AX = t10
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 0(SP)        // t0 = c & 0x3FFFFFF
	SHRQ $26, R8
	
	// c = c + t1 + t10*0x400 + t11*0x3D10
	MOVQ 8(SP), R9        // R9 = t1
	ADDQ R9, R8
	MOVQ 80(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 88(SP), AX       // AX = t11
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 8(SP)        // t1 = c & 0x3FFFFFF
	SHRQ $26, R8
	
	// c = c + t2 + t11*0x400 + t12*0x3D10
	MOVQ 16(SP), R9
	ADDQ R9, R8
	MOVQ 88(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 96(SP), AX       // AX = t12
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 16(SP)       // t2 = c & 0x3FFFFFF
	SHRQ $26, R8
	
	// r.n[3] = (c + t3 + t12*0x400 + t13*0x3D10) & 0x3FFFFFF
	MOVQ 24(SP), R9
	ADDQ R9, R8
	MOVQ 96(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 104(SP), AX      // AX = t13
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDL $0x3FFFFFF, R9   // 32-bit AND
	MOVL R9, 12(DI)       // r.n[3]
	SHRQ $26, R8
	
	// r.n[4] = (c + t4 + t13*0x400 + t14*0x3D10) & 0x3FFFFFF
	MOVQ 32(SP), R9
	ADDQ R9, R8
	MOVQ 104(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 112(SP), AX      // AX = t14
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDL $0x3FFFFFF, R9
	MOVL R9, 16(DI)       // r.n[4]
	SHRQ $26, R8
	
	// r.n[5] = (c + t5 + t14*0x400 + t15*0x3D10) & 0x3FFFFFF
	MOVQ 40(SP), R9
	ADDQ R9, R8
	MOVQ 112(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 120(SP), AX      // AX = t15
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDL $0x3FFFFFF, R9
	MOVL R9, 20(DI)       // r.n[5]
	SHRQ $26, R8
	
	// r.n[6] = (c + t6 + t15*0x400 + t16*0x3D10) & 0x3FFFFFF
	MOVQ 48(SP), R9
	ADDQ R9, R8
	MOVQ 120(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 128(SP), AX      // AX = t16
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDL $0x3FFFFFF, R9
	MOVL R9, 24(DI)       // r.n[6]
	SHRQ $26, R8
	
	// r.n[7] = (c + t7 + t16*0x400 + t17*0x3D10) & 0x3FFFFFF
	MOVQ 56(SP), R9
	ADDQ R9, R8
	MOVQ 128(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 136(SP), AX      // AX = t17
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDL $0x3FFFFFF, R9
	MOVL R9, 28(DI)       // r.n[7]
	SHRQ $26, R8
	
	// r.n[8] = (c + t8 + t17*0x400 + t18*0x3D10) & 0x3FFFFFF
	MOVQ 64(SP), R9
	ADDQ R9, R8
	MOVQ 136(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 144(SP), AX      // AX = t18
	IMULQ $0x3D10, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDL $0x3FFFFFF, R9
	MOVL R9, 32(DI)       // r.n[8]
	SHRQ $26, R8
	
	// r.n[9] = (c + t9 + t18*0x400 + t19*0x1000003D10) & 0x03FFFFF
	MOVQ 72(SP), R9
	ADDQ R9, R8
	MOVQ 144(SP), AX
	IMULQ $0x400, AX
	ADDQ AX, R8
	MOVQ 152(SP), AX      // AX = t19
	MOVQ $0x1000003D10, CX
	IMULQ CX, AX
	ADDQ AX, R8
	MOVQ R8, R9
	ANDL $0x03FFFFF, R9
	MOVL R9, 36(DI)       // r.n[9]
	SHRQ $22, R8          // Note: shift by 22, not 26
	
	// Final carry propagation
	// Save c before we modify it
	MOVQ R8, R15          // R15 = c (saved)
	
	// d = t0 + c*0x3D1
	MOVQ 0(SP), R10       // R10 = t0
	MOVQ R8, AX           // AX = c
	IMULQ $0x3D1, AX      // AX = c * 0x3D1
	ADDQ AX, R10          // R10 = t0 + c*0x3D1
	MOVQ R10, R9
	ANDL $0x3FFFFFF, R9
	MOVL R9, 0(DI)        // r.n[0]
	SHRQ $26, R10
	
	// d = d + t1 + c*0x40
	MOVQ 8(SP), R11
	ADDQ R10, R11
	MOVQ R15, AX          // AX = c (original value)
	IMULQ $0x40, AX       // AX = c * 0x40
	ADDQ AX, R11
	MOVQ R11, R9
	ANDL $0x3FFFFFF, R9
	MOVL R9, 4(DI)        // r.n[1]
	SHRQ $26, R11
	
	// r.n[2] = t2 + d
	MOVQ 16(SP), R12
	ADDQ R11, R12
	MOVL R12, 8(DI)       // r.n[2]
	
	RET

// Field squaring: r = a * a mod p
// func fieldSqrAsm(r, a *Field)
//
// Optimized squaring: computes cross-products once
// Algorithm: for i<j, compute 2*a[i]*a[j] (not both a[i]*a[j] and a[j]*a[i])
//            for i=j, compute a[i]*a[i] (diagonal terms)
//
TEXT ·fieldSqrAsm(SB), NOSPLIT, $160-16
	MOVQ r+0(FP), DI     // DI = r (result)
	MOVQ a+8(FP), SI     // SI = a (operand)
	
	// Register allocation:
	// R8  = c (carry/accumulator)
	// R9  = temp for storing masked limb
	// R10-R15 = scratch
	// AX, CX, DX = used by IMULQ
	
	// ===== Phase 1: Compute cross-products t0-t19 =====
	// Using optimized squaring: 2*a[i]*a[j] for i<j, a[i]*a[i] for i=j
	
	// t0 = a[0] * a[0]
	MOVL (SI), AX
	IMULQ AX, AX
	MOVQ AX, R8           // c = a[0] * a[0]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9   // t0 = c & 0x3FFFFFF
	MOVQ R9, 0(SP)
	SHRQ $26, R8          // c >>= 26
	
	// t1 = c + 2*a[0]*a[1]
	MOVL (SI), AX
	MOVL 4(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // Double it
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 8(SP)
	SHRQ $26, R8
	
	// t2 = c + 2*a[0]*a[2] + a[1]*a[1]
	MOVL (SI), AX
	MOVL 8(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[2]
	MOVL 4(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // + a[1]*a[1]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 16(SP)
	SHRQ $26, R8
	
	// t3 = c + 2*a[0]*a[3] + 2*a[1]*a[2]
	MOVL (SI), AX
	MOVL 12(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[3]
	MOVL 4(SI), AX
	MOVL 8(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[2]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 24(SP)
	SHRQ $26, R8
	
	// t4 = c + 2*a[0]*a[4] + 2*a[1]*a[3] + a[2]*a[2]
	MOVL (SI), AX
	MOVL 16(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[4]
	MOVL 4(SI), AX
	MOVL 12(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[3]
	MOVL 8(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[2]*a[2]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 32(SP)
	SHRQ $26, R8
	
	// t5 = c + 2*a[0]*a[5] + 2*a[1]*a[4] + 2*a[2]*a[3]
	MOVL (SI), AX
	MOVL 20(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[5]
	MOVL 4(SI), AX
	MOVL 16(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[4]
	MOVL 8(SI), AX
	MOVL 12(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[2]*a[3]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 40(SP)
	SHRQ $26, R8
	
	// t6 = c + 2*a[0]*a[6] + 2*a[1]*a[5] + 2*a[2]*a[4] + a[3]*a[3]
	MOVL (SI), AX
	MOVL 24(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[6]
	MOVL 4(SI), AX
	MOVL 20(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[5]
	MOVL 8(SI), AX
	MOVL 16(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[2]*a[4]
	MOVL 12(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[3]*a[3]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 48(SP)
	SHRQ $26, R8
	
	// t7 = c + 2*a[0]*a[7] + 2*a[1]*a[6] + 2*a[2]*a[5] + 2*a[3]*a[4]
	MOVL (SI), AX
	MOVL 28(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[7]
	MOVL 4(SI), AX
	MOVL 24(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[6]
	MOVL 8(SI), AX
	MOVL 20(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[2]*a[5]
	MOVL 12(SI), AX
	MOVL 16(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[3]*a[4]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 56(SP)
	SHRQ $26, R8
	
	// t8 = c + 2*a[0]*a[8] + 2*a[1]*a[7] + 2*a[2]*a[6] + 2*a[3]*a[5] + a[4]*a[4]
	MOVL (SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[8]
	MOVL 4(SI), AX
	MOVL 28(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[7]
	MOVL 8(SI), AX
	MOVL 24(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[2]*a[6]
	MOVL 12(SI), AX
	MOVL 20(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[3]*a[5]
	MOVL 16(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[4]*a[4]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 64(SP)
	SHRQ $26, R8
	
	// t9 = c + 2*a[0]*a[9] + 2*a[1]*a[8] + 2*a[2]*a[7] + 2*a[3]*a[6] + 2*a[4]*a[5]
	MOVL (SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[0]*a[9]
	MOVL 4(SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[8]
	MOVL 8(SI), AX
	MOVL 28(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[2]*a[7]
	MOVL 12(SI), AX
	MOVL 24(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[3]*a[6]
	MOVL 16(SI), AX
	MOVL 20(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[4]*a[5]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 72(SP)
	SHRQ $26, R8
	
	// t10 = c + 2*a[1]*a[9] + 2*a[2]*a[8] + 2*a[3]*a[7] + 2*a[4]*a[6] + a[5]*a[5]
	MOVL 4(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[1]*a[9]
	MOVL 8(SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[2]*a[8]
	MOVL 12(SI), AX
	MOVL 28(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[3]*a[7]
	MOVL 16(SI), AX
	MOVL 24(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[4]*a[6]
	MOVL 20(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[5]*a[5]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 80(SP)
	SHRQ $26, R8
	
	// t11 = c + 2*a[2]*a[9] + 2*a[3]*a[8] + 2*a[4]*a[7] + 2*a[5]*a[6]
	MOVL 8(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[2]*a[9]
	MOVL 12(SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[3]*a[8]
	MOVL 16(SI), AX
	MOVL 28(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[4]*a[7]
	MOVL 20(SI), AX
	MOVL 24(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[5]*a[6]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 88(SP)
	SHRQ $26, R8
	
	// t12 = c + 2*a[3]*a[9] + 2*a[4]*a[8] + 2*a[5]*a[7] + a[6]*a[6]
	MOVL 12(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[3]*a[9]
	MOVL 16(SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[4]*a[8]
	MOVL 20(SI), AX
	MOVL 28(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[5]*a[7]
	MOVL 24(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[6]*a[6]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 96(SP)
	SHRQ $26, R8
	
	// t13 = c + 2*a[4]*a[9] + 2*a[5]*a[8] + 2*a[6]*a[7]
	MOVL 16(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[4]*a[9]
	MOVL 20(SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[5]*a[8]
	MOVL 24(SI), AX
	MOVL 28(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[6]*a[7]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 104(SP)
	SHRQ $26, R8
	
	// t14 = c + 2*a[5]*a[9] + 2*a[6]*a[8] + a[7]*a[7]
	MOVL 20(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[5]*a[9]
	MOVL 24(SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[6]*a[8]
	MOVL 28(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[7]*a[7]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 112(SP)
	SHRQ $26, R8
	
	// t15 = c + 2*a[6]*a[9] + 2*a[7]*a[8]
	MOVL 24(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[6]*a[9]
	MOVL 28(SI), AX
	MOVL 32(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[7]*a[8]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 120(SP)
	SHRQ $26, R8
	
	// t16 = c + 2*a[7]*a[9] + a[8]*a[8]
	MOVL 28(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[7]*a[9]
	MOVL 32(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[8]*a[8]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 128(SP)
	SHRQ $26, R8
	
	// t17 = c + 2*a[8]*a[9]
	MOVL 32(SI), AX
	MOVL 36(SI), CX
	IMULQ CX, AX
	ADDQ AX, R8
	ADDQ AX, R8           // 2*a[8]*a[9]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 136(SP)
	SHRQ $26, R8
	
	// t18 = c + a[9]*a[9]
	MOVL 36(SI), AX
	IMULQ AX, AX
	ADDQ AX, R8           // a[9]*a[9]
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 144(SP)
	SHRQ $26, R8
	
	// t19 = c
	MOVQ R8, 152(SP)
	
	// ===== Phase 2: Modular reduction =====
	// Same as fieldMulAsm
	
	// c = t0 + t10*0x3D10
	MOVQ 0(SP), R8        // t0
	MOVQ 80(SP), R10      // t10
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// t0 = c & 0x3FFFFFF
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 0(SP)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t1 + t10*0x400 + t11*0x3D10
	MOVQ 8(SP), R10       // t1
	ADDQ R10, R8
	MOVQ 80(SP), R10      // t10
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 88(SP), R10      // t11
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// t1 = c & 0x3FFFFFF
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 8(SP)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t2 + t11*0x400 + t12*0x3D10
	MOVQ 16(SP), R10      // t2
	ADDQ R10, R8
	MOVQ 88(SP), R10      // t11
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 96(SP), R10      // t12
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// t2 = c & 0x3FFFFFF
	MOVQ R8, R9
	ANDQ $0x3FFFFFF, R9
	MOVQ R9, 16(SP)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t3 + t12*0x400 + t13*0x3D10
	MOVQ 24(SP), R10      // t3
	ADDQ R10, R8
	MOVQ 96(SP), R10      // t12
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 104(SP), R10     // t13
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// r.n[3] = c & 0x3FFFFFF
	MOVQ R8, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, 12(DI)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t4 + t13*0x400 + t14*0x3D10
	MOVQ 32(SP), R10      // t4
	ADDQ R10, R8
	MOVQ 104(SP), R10     // t13
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 112(SP), R10     // t14
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// r.n[4] = c & 0x3FFFFFF
	MOVQ R8, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, 16(DI)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t5 + t14*0x400 + t15*0x3D10
	MOVQ 40(SP), R10      // t5
	ADDQ R10, R8
	MOVQ 112(SP), R10     // t14
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 120(SP), R10     // t15
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// r.n[5] = c & 0x3FFFFFF
	MOVQ R8, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, 20(DI)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t6 + t15*0x400 + t16*0x3D10
	MOVQ 48(SP), R10      // t6
	ADDQ R10, R8
	MOVQ 120(SP), R10     // t15
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 128(SP), R10     // t16
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// r.n[6] = c & 0x3FFFFFF
	MOVQ R8, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, 24(DI)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t7 + t16*0x400 + t17*0x3D10
	MOVQ 56(SP), R10      // t7
	ADDQ R10, R8
	MOVQ 128(SP), R10     // t16
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 136(SP), R10     // t17
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// r.n[7] = c & 0x3FFFFFF
	MOVQ R8, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, 28(DI)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t8 + t17*0x400 + t18*0x3D10
	MOVQ 64(SP), R10      // t8
	ADDQ R10, R8
	MOVQ 136(SP), R10     // t17
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 144(SP), R10     // t18
	IMULQ $0x3D10, R10
	ADDQ R10, R8
	
	// r.n[8] = c & 0x3FFFFFF
	MOVQ R8, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, 32(DI)
	
	// c = c >> 26
	SHRQ $26, R8
	
	// c = c + t9 + t18*0x400 + t19*0x1000003D10
	MOVQ 72(SP), R10      // t9
	ADDQ R10, R8
	MOVQ 144(SP), R10     // t18
	IMULQ $0x400, R10
	ADDQ R10, R8
	MOVQ 152(SP), R10     // t19
	MOVQ R10, R15         // Save t19 to R15
	MOVQ $0x1000003D10, CX
	IMULQ CX, R10
	ADDQ R10, R8
	
	// r.n[9] = c & 0x03FFFFF
	MOVQ R8, R12
	ANDQ $0x03FFFFF, R12
	MOVL R12, 36(DI)
	
	// c = c >> 22
	SHRQ $22, R8
	
	// Final reduction pass
	// d = t0 + c*0x3D1
	MOVQ 0(SP), R11       // t0
	MOVQ R8, R10
	IMULQ $0x3D1, R10
	ADDQ R10, R11
	
	// r.n[0] = d & 0x3FFFFFF
	MOVQ R11, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, (DI)
	
	// d = d >> 26
	SHRQ $26, R11
	
	// d = d + t1 + c*0x40
	MOVQ 8(SP), R12       // t1
	ADDQ R12, R11
	MOVQ R8, R10
	IMULQ $0x40, R10
	ADDQ R10, R11
	
	// r.n[1] = d & 0x3FFFFFF
	MOVQ R11, R12
	ANDQ $0x3FFFFFF, R12
	MOVL R12, 4(DI)
	
	// d = d >> 26
	SHRQ $26, R11
	
	// r.n[2] = t2 + d
	MOVQ 16(SP), R12
	ADDQ R11, R12
	MOVL R12, 8(DI)
	
	RET
