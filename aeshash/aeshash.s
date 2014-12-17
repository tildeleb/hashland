// Copyright © 2014 Lawrence E. Bakst. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Go's hash function used by map on X64 hardware with AESNI
// liberated from go runtime/asm_amd64.s

#include "textflag.h"
#include "funcdata.h"

// func Hash(b []byte, seed uint64) uint64
TEXT ·Hash(SB),NOSPLIT,$0-40
	MOVQ	b_base+0(FP), AX	// ptr to bytes
	MOVQ	b_len+8(FP), CX		// length of slice
	MOVQ	seed+24(FP), X0		// seed to low 64 bits of xmm0
	CALL	·aeshashbody(SB)
	MOVQ	X0, ret+32(FP)
	RET

// func HashStr(s string, seed uint64) uint64
TEXT ·HashStr(SB),NOSPLIT,$0-32
	MOVQ	s_base+0(FP), AX	// ptr to string data
	MOVQ	s_len+8(FP), CX		// length of string
	MOVQ	seed+16(FP), X0		// seed to low 64 bits of xmm0
	CALL	·aeshashbody(SB)
	MOVQ	X0, ret+24(FP)
	RET

// AX: data
// CX: length
// X0: seed
// func aeshashbody()
TEXT ·aeshashbody(SB),NOSPLIT,$0-0
	PINSRQ	$1, CX, X0		// size to high 64 bits of xmm0
	MOVO	·aeskeysched+0(SB), X2
	MOVO	·aeskeysched+16(SB), X3
	CMPQ	CX, $16
	JB	aessmall
aesloop:
	CMPQ	CX, $16
	JBE	aesloopend
	MOVOU	(AX), X1
	AESENC	X2, X0
	AESENC	X1, X0
	SUBQ	$16, CX
	ADDQ	$16, AX
	JMP	aesloop
// 1-16 bytes remaining
aesloopend:
	// This load may overlap with the previous load above.
	// We'll hash some bytes twice, but that's ok.
	MOVOU	-16(AX)(CX*1), X1
	JMP	partial
// 0-15 bytes
aessmall:
	TESTQ	CX, CX
	JE	finalize	// 0 bytes

	CMPB	AX, $0xf0
	JA	highpartial

	// 16 bytes loaded at this address won't cross
	// a page boundary, so we can load it directly.
	MOVOU	(AX), X1
	ADDQ	CX, CX
	MOVQ	$masks<>(SB), BP
	PAND	(BP)(CX*8), X1
	JMP	partial
highpartial:
	// address ends in 1111xxxx.  Might be up against
	// a page boundary, so load ending at last byte.
	// Then shift bytes down using pshufb.
	MOVOU	-16(AX)(CX*1), X1
	ADDQ	CX, CX
	MOVQ	$shifts<>(SB), BP
	PSHUFB	(BP)(CX*8), X1
partial:
	// incorporate partial block into hash
	AESENC	X3, X0
	AESENC	X1, X0
finalize:	
	// finalize hash
	AESENC	X2, X0
	AESENC	X3, X0
	AESENC	X2, X0
aesret:
	RET


// put the seed s into the low 64 bits of xmm0
// put the data v into the high 64 bits of xmm0
// perform 3 AES rounds with 2 alternating round keys
// func Hash64(k uint64, seed uint64) uint64
TEXT ·Hash64(SB),NOSPLIT,$0-24
	MOVQ	seed+8(FP), X0	// seed
	MOVQ	k+0(FP), AX		// data
	PINSRQ	$1, AX, X0		// 64 bit data key to high order 64 bits of X0
	AESENC	·aeskeysched+0(SB), X0
	AESENC	·aeskeysched+16(SB), X0
	AESENC	·aeskeysched+0(SB), X0
	MOVQ	X0, ret+16(FP)
	RET

// func Hash32(k uint32, seed uint64) uint64
TEXT ·Hash32(SB),NOSPLIT,$0-24
	MOVQ	seed+8(FP), X0	// seed
	MOVQ	k+0(FP), AX		// 32 bit data key
	PINSRD	$2, AX, X0		// data to the low order 32 bits of the high order 64 bits
	PINSRD	$3, AX, X0		// data to the high order 32 bits of the high order 64 bits
	AESENC	·aeskeysched+0(SB), X0
	AESENC	·aeskeysched+16(SB), X0
	AESENC	·aeskeysched+0(SB), X0
	MOVQ	X0, ret+16(FP)
	RET


// simple mask to get rid of data in the high part of the register.
// var masks [32]uint64
DATA masks<>+0x00(SB)/8, $0x0000000000000000
DATA masks<>+0x08(SB)/8, $0x0000000000000000
DATA masks<>+0x10(SB)/8, $0x00000000000000ff
DATA masks<>+0x18(SB)/8, $0x0000000000000000
DATA masks<>+0x20(SB)/8, $0x000000000000ffff
DATA masks<>+0x28(SB)/8, $0x0000000000000000
DATA masks<>+0x30(SB)/8, $0x0000000000ffffff
DATA masks<>+0x38(SB)/8, $0x0000000000000000
DATA masks<>+0x40(SB)/8, $0x00000000ffffffff
DATA masks<>+0x48(SB)/8, $0x0000000000000000
DATA masks<>+0x50(SB)/8, $0x000000ffffffffff
DATA masks<>+0x58(SB)/8, $0x0000000000000000
DATA masks<>+0x60(SB)/8, $0x0000ffffffffffff
DATA masks<>+0x68(SB)/8, $0x0000000000000000
DATA masks<>+0x70(SB)/8, $0x00ffffffffffffff
DATA masks<>+0x78(SB)/8, $0x0000000000000000
DATA masks<>+0x80(SB)/8, $0xffffffffffffffff
DATA masks<>+0x88(SB)/8, $0x0000000000000000
DATA masks<>+0x90(SB)/8, $0xffffffffffffffff
DATA masks<>+0x98(SB)/8, $0x00000000000000ff
DATA masks<>+0xa0(SB)/8, $0xffffffffffffffff
DATA masks<>+0xa8(SB)/8, $0x000000000000ffff
DATA masks<>+0xb0(SB)/8, $0xffffffffffffffff
DATA masks<>+0xb8(SB)/8, $0x0000000000ffffff
DATA masks<>+0xc0(SB)/8, $0xffffffffffffffff
DATA masks<>+0xc8(SB)/8, $0x00000000ffffffff
DATA masks<>+0xd0(SB)/8, $0xffffffffffffffff
DATA masks<>+0xd8(SB)/8, $0x000000ffffffffff
DATA masks<>+0xe0(SB)/8, $0xffffffffffffffff
DATA masks<>+0xe8(SB)/8, $0x0000ffffffffffff
DATA masks<>+0xf0(SB)/8, $0xffffffffffffffff
DATA masks<>+0xf8(SB)/8, $0x00ffffffffffffff
GLOBL masks<>(SB), RODATA, $256

// these are arguments to pshufb.  They move data down from
// the high bytes of the register to the low bytes of the register.
// index is how many bytes to move.
// var shifts [32]uint64
DATA shifts<>+0x00(SB)/8, $0x0000000000000000
DATA shifts<>+0x08(SB)/8, $0x0000000000000000
DATA shifts<>+0x10(SB)/8, $0xffffffffffffff0f
DATA shifts<>+0x18(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x20(SB)/8, $0xffffffffffff0f0e
DATA shifts<>+0x28(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x30(SB)/8, $0xffffffffff0f0e0d
DATA shifts<>+0x38(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x40(SB)/8, $0xffffffff0f0e0d0c
DATA shifts<>+0x48(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x50(SB)/8, $0xffffff0f0e0d0c0b
DATA shifts<>+0x58(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x60(SB)/8, $0xffff0f0e0d0c0b0a
DATA shifts<>+0x68(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x70(SB)/8, $0xff0f0e0d0c0b0a09
DATA shifts<>+0x78(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x80(SB)/8, $0x0f0e0d0c0b0a0908
DATA shifts<>+0x88(SB)/8, $0xffffffffffffffff
DATA shifts<>+0x90(SB)/8, $0x0e0d0c0b0a090807
DATA shifts<>+0x98(SB)/8, $0xffffffffffffff0f
DATA shifts<>+0xa0(SB)/8, $0x0d0c0b0a09080706
DATA shifts<>+0xa8(SB)/8, $0xffffffffffff0f0e
DATA shifts<>+0xb0(SB)/8, $0x0c0b0a0908070605
DATA shifts<>+0xb8(SB)/8, $0xffffffffff0f0e0d
DATA shifts<>+0xc0(SB)/8, $0x0b0a090807060504
DATA shifts<>+0xc8(SB)/8, $0xffffffff0f0e0d0c
DATA shifts<>+0xd0(SB)/8, $0x0a09080706050403
DATA shifts<>+0xd8(SB)/8, $0xffffff0f0e0d0c0b
DATA shifts<>+0xe0(SB)/8, $0x0908070605040302
DATA shifts<>+0xe8(SB)/8, $0xffff0f0e0d0c0b0a
DATA shifts<>+0xf0(SB)/8, $0x0807060504030201
DATA shifts<>+0xf8(SB)/8, $0xff0f0e0d0c0b0a09
GLOBL shifts<>(SB), RODATA, $256

