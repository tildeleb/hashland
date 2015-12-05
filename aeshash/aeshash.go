package aeshash

import _ "unsafe"
import "github.com/tildeleb/hashland/nhash"

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

// Make sure interfaces are correctly implemented. Stolen from another implementation.
// I did something similar in another package to verify the interface but didn't know you could elide the variable in a var.
// What a cute wart it is.
var (
	//_ hash.Hash   = new(Digest)
	_ nhash.Hash64 = new(StateAES)
	_  nhash.HashStream = new(StateAES)
)

type StateAES struct {
	hash	uint64
	seed	uint64
	clen	int
	tail	[]byte
}

func NewAES(seed uint64) nhash.Hash64 {
	s := new(StateAES)
	s.seed = seed
	s.Reset()
	return s
}

// Return the size of the resulting hash.
func (d *StateAES) Size() int { return 8 }

// Return the blocksize of the hash which in this case is 1 byte.
func (d *StateAES) BlockSize() int { return 1 }

// Return the maximum number of seed bypes required. In this case 2 x 32
func (d *StateAES) NumSeedBytes() int {
	return 8
}

// Return the number of bits the hash function outputs.
func (d *StateAES) HashSizeInBits() int {
	return 64
}

// Reset the hash state.
func (d *StateAES) Reset() {
	d.hash = 0
	d.clen = 0
	d.tail = nil
}

// Accept a byte stream p used for calculating the hash. For now this call is lazy and the actual hash calculations take place in Sum() and Sum32().
func (d *StateAES) Write(p []byte) (nn int, err error) {
	l := len(p)
	d.clen += l
	d.tail = append(d.tail, p...)
	return l, nil
}

func (d *StateAES) Write64(h uint64) (err error) {
	d.clen += 8
	d.tail = append(d.tail, byte(h>>56), byte(h>>48), byte(h>>40), byte(h>>32), byte(h>>24), byte(h>>16), byte(h>>8), byte(h))
	return nil
}

// Return the current hash as a byte slice.
func (d *StateAES) Sum(b []byte) []byte {
	d.hash = Hash(d.tail, d.seed)
	h := d.hash
	return append(b, byte(h>>56), byte(h>>48), byte(h>>40), byte(h>>32), byte(h>>24), byte(h>>16), byte(h>>8), byte(h))
}

// Return the current hash as a 64 bit unsigned type.
func (d *StateAES) Sum64() uint64 {
	d.hash = Hash(d.tail, d.seed)
	return d.hash
}

func (d *StateAES) Hash64(b []byte, seeds ...uint64) uint64 {
	switch len(seeds) {
	case 1:
		d.seed = seeds[0]
	}
	d.hash = Hash(b, d.seed)
	//fmt.Printf("pc=0x%08x, pb=0x%08x\n", d.pc, d.pb)
	return d.hash
}
