package aeshash

import _ "unsafe"

var masks [32]uint64
var shifts [32]uint64

// used in asm_{386,amd64}.s
const hashRandomBytes = 32

// this is really 2 x 128 bit round keys
var aeskeysched[hashRandomBytes]byte

var aesdebug[hashRandomBytes]byte

func aeshashbody()

//func Hash(p unsafe.Pointer, s, h uintptr) uintptr
//func HashStr(p string, s, h uintptr) uintptr
func Hash(b []byte, seed uint64) uint64
func HashStr(s string, seed uint64) uint64
func Hash64(v uint64, s uint64) uint64
func Hash32(v uint32, s uint64) uint64

//func aeshash(p unsafe.Pointer, s, h uintptr) uintptr
//func aeshash32(p unsafe.Pointer, s, h uintptr) uintptr
//func aeshash64(p unsafe.Pointer, s, h uintptr) uintptr
//func aeshashstr(p unsafe.Pointer, s, h uintptr) uintptr

func init() {
	p := aeskeysched[:]
	p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7] = 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8
	p[8], p[9], p[10], p[11], p[12], p[13], p[14], p[15] = 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF, 0x10
	p[16], p[17], p[18], p[19], p[20], p[21], p[22], p[23] = 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18
	p[24], p[25], p[26], p[27], p[28], p[29], p[30], p[31] = 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0xFF
	p = aesdebug[:]
	p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7] = 0xFF, 0, 0, 0, 0, 0, 0, 0xFE
	p[8], p[9], p[10], p[11], p[12], p[13], p[14], p[15] = 0xFD, 0, 0, 0, 0, 0, 0, 0xFC
}

