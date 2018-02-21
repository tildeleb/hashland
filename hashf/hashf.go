// Copyright © 2014,2015 Lawrence E. Bakst. All rights reserved.

package hashf

import (
	"crypto/sha1"
	"fmt"
	farm "github.com/dgryski/go-farm"
	metro "github.com/dgryski/go-metro"
	"github.com/jzelinskie/whirlpool"
	"github.com/minio/blake2b-simd"
	"github.com/tildeleb/aeshash"
	"github.com/tildeleb/hashland/crapwow"
	"github.com/tildeleb/hashland/gomap"
	"github.com/tildeleb/hashland/jenkins"
	"github.com/tildeleb/hashland/keccak"
	"github.com/tildeleb/hashland/keccakpg"
	"github.com/tildeleb/hashland/mahash"
	"github.com/tildeleb/hashland/murmur3"
	"github.com/tildeleb/hashland/nhash"
	"github.com/tildeleb/hashland/nullhash"
	"github.com/tildeleb/hashland/sbox"
	"github.com/tildeleb/hashland/siphash"
	"github.com/tildeleb/hashland/siphashpg"
	"github.com/tildeleb/hashland/skein"
	"github.com/tildeleb/hashland/spooky"
	"hash"
	"hash/adler32"
	"time"
	"unsafe"
)

// interfaces
var nh hash.Hash64
var nhf64 nhash.HashF64
var a32 hash.Hash32
var k643 hash.Hash
var k644 hash.Hash
var k648 hash.Hash
var k160 hash.Hash
var k224 hash.Hash
var skein256 hash.Hash
var sha1160 hash.Hash
var m332 hash.Hash32
var m364 hash.Hash64
var m3128 murmur3.Hash128
var wp hash.Hash
var b2b hash.Hash

// functions
var j264 = jenkins.Hash264
var j364 = jenkins.Jenkins364
var null = nullhashf
var aesg = aeshash.Hash

var Hf2 string // wow this has to go
type ff func() time.Duration

func nullhashf(b []byte, seed uint64) uint64 {
	return 0
}

type DispEntry struct {
	fp   unsafe.Pointer
	hi   hash.Hash
	hi32 hash.Hash32
	hi64 hash.Hash64
	kind int
}

const (
	h64s     = 2
	b32s     = iota
	b32x2sx2 = iota
)

type HashFunction struct {
	Name   string
	Size   int // in bits
	Crypto bool
	desc   string
	de     *DispEntry
	//	dummy		*int		// de			DispEntry
}

var nullhashfp = nullhash.Nullhash
var fnull = nullhashf
var faes = aeshash.Hash
var fgomap = gomap.Hash64

var HashFunctions = map[string]HashFunction{
	"perfecthash":   HashFunction{"perfecthash", 64, true, "perfecthash, 64 bit", nil},
	"nullhash":      HashFunction{"nullhash", 64, true, "nullhash, 64 bit", &DispEntry{fp: unsafe.Pointer(&fnull), kind: h64s}},
	"nullhashF64ns": HashFunction{"nullhashF64ns", 64, true, "nullhashF64ns, 64 bit, no seed", nil},
	"aeshash64":     HashFunction{"aeshash64", 64, true, "aeshash, 64 bit, accelerated", &DispEntry{fp: unsafe.Pointer(&faes), kind: h64s}},
	"siphash64":     HashFunction{"siphash64", 64, true, "siphash, 64 bit, accelerated", nil},
	"siphash64pg":   HashFunction{"siphash64pg", 64, true, "siphash, pure go, 64 bit, a bits", nil},
	/*
		"siphash64":		HashFunction{"siphash64", 		64,		true,	"siphash, 64 bit, a bits", nil},
		"siphash128a":		HashFunction{"siphasha", 		64,		true,	"siphash, 128 bit, a bits", nil},
		"siphash128b":		HashFunction{"siphashb", 		64,		true,	"siphash, 128 bit, b bits", nil},
		"siphash64al":		HashFunction{"siphash64al", 	32,		true,	"siphash, 64 bit, a bits, low", nil},
		"siphash64ah":		HashFunction{"siphash64ah", 	32,		true,	"siphash, 64 bit, a bits, high", nil},
		"siphash64bl":		HashFunction{"siphash64bl", 	32,		true,	"siphash, 128 bit, b bits, low", nil},
		"siphash64bh":		HashFunction{"siphash64bh", 	32,		true,	"siphash, 128 bit, b bits, high", nil},
	*/
	"MaHash8v64": HashFunction{"MaHash8v64", 64, false, "russian hash function", nil},

	// tribute to Robert Jenkins goes here
	"spooky32":     HashFunction{"spooky32", 32, false, "jenkins, spooky, 32 bit", nil},
	"spooky64":     HashFunction{"spooky64", 64, false, "jenkins, spooky, 64 bit", nil},
	"spooky128h":   HashFunction{"spooky128h", 64, false, "jenkins, spooky, 128 bit, high bits", nil},
	"spooky128l":   HashFunction{"spooky128l", 64, false, "jenkins, spooky, 128 bit, low bits", nil},
	"spooky128xor": HashFunction{"spooky128xor", 64, false, "jenkins, spooky, 128, high xor low bits", nil},
	"j364":         HashFunction{"j364", 64, false, "jenkins, lookup3. 64 bit, c low order bits, b high order bits", nil},
	"j264":         HashFunction{"j264", 64, false, "jenkins, lookup8. 64 bit", nil},
	"j332c":        HashFunction{"j332c", 32, false, "jenkins, lookup3, 32 bit, c bits", nil},
	"j332b":        HashFunction{"j332b", 32, false, "jenkins, lookup3, 32 bit, b bits", nil},
	"j232":         HashFunction{"j232", 32, false, "jenkins, lookup8, 32 bit", nil},
	"j264l":        HashFunction{"j264l", 32, false, "jenkins, lookup8, 64 bit, low bits", nil},
	"j264h":        HashFunction{"j264h", 32, false, "jenkins, lookup8, 64 bit, high bits", nil},
	"j264xor":      HashFunction{"j264xor", 32, false, "jenkins, lookup8, 64 bit, high xor low bits", nil},
	"sbox":         HashFunction{"sbox", 32, false, "sbox", nil},

	"gomap32": HashFunction{"gomap32", 32, false, "gomap32", nil},
	"gomap64": HashFunction{"gomap64", 64, false, "gomap64", &DispEntry{fp: unsafe.Pointer(&fgomap), kind: h64s}},

	"murmur332": HashFunction{"murmur332", 32, false, "murmur332", nil},
	"murmur364": HashFunction{"murmur364", 64, false, "murmur364", nil},

	"FarmHash32":  HashFunction{"FarmHash32", 32, false, "FarmHash32", nil},
	"FarmHash64":  HashFunction{"FarmHash64", 64, false, "FarmHash64", nil},
	"FarmHash128": HashFunction{"FarmHash128", 128, false, "FarmHash128", nil},

	"MetroHash64-1":  HashFunction{"MetroHash64-1", 64, false, "MetroHash64-1", nil},
	"MetroHash64-2":  HashFunction{"MetroHash64-2", 64, false, "MetroHash64-2", nil},
	"MetroHash128-1": HashFunction{"MetroHash128-1", 128, false, "MetroHash128-1", nil},
	"MetroHash128-2": HashFunction{"MetroHash128-2", 128, false, "MetroHash128-2", nil},

	"keccak224":   HashFunction{"keccak224", 64, true, "keccak, 224 bit to 64 bit", nil},
	"keccakpg643": HashFunction{"keccak643", 64, true, "keccak, 64 bit, 3 rounds", nil},
	"keccakpg644": HashFunction{"keccak644", 64, true, "keccak, 64 bit, 4 rounds", nil},
	"keccakpg648": HashFunction{"keccak648", 64, true, "keccak, 64 bit, 8 rounds", nil},
	"skein256":    HashFunction{"skein256", 64, true, "skein256, 64 bit , low 64 bits", nil},
	"sha1":        HashFunction{"sha1", 64, true, "sha1, 160 bit hash", nil},
	"keccak160":   HashFunction{"keccak160", 64, true, "keccak160l", nil},

	"skein256low": HashFunction{"skein256low", 32, true, "skein256low", nil},
	"skein256hi":  HashFunction{"skein256hi", 32, true, "skein256hi", nil},
	"skein256xor": HashFunction{"skein256xor", 32, true, "skein256xor", nil},

	"whirlpool": HashFunction{"whirlpool", 64, true, "whirlpool", nil},
	"blake2b":   HashFunction{"blake2b", 64, true, "blake2b", nil},

	// so bad we can't include them
	"CrapWow": HashFunction{"CrapWow", 32, false, "CrapWow", nil},
	"adler32": HashFunction{"adler32", 32, false, "adler32", nil},
}

// "CrapWow" removed because it generates so many dup hashes with duplicated words it goes from O(1) to O(N)
// "adler32" removed for the same reasons
// 	"siphash64al", "siphash64ah", "siphash64bl", "siphash64bh",
// 	"skein256xor", "skein256low", "skein256hi", "sha1", "keccak160l",
// 	"siphash64", "siphash128a", "siphash128b",
// 	"keccak644", "keccak648" "keccak160",
var TestHashFunctions = []string{"nullhash", //"perfecthash",
	"aeshash64", "gomap64", "j364", "j264", "murmur364",
	"siphash64",
	"siphash64pg",
	"MaHash8v64", "spooky64", "spooky128h", "spooky128l", "spooky128xor",
	"murmur332", "j332c", "j332b", "j232", "j264l", "j264h", "j264xor", "spooky32", "sbox", "gomap32",
	"FarmHash32",
	"FarmHash64",
	"FarmHash128-low", "FarmHash128-high", "FarmHash128-xor",
	"MetroHash128-2l", "MetroHash128-2h", "MetroHash128-2xor",
	"MetroHash64-1", "MetroHash64-2",
	"MetroHash128-1l", "MetroHash128-1h", "MetroHash128-1xor",
	"MetroHash128-2l", "MetroHash128-2h", "MetroHash128-2xor",
	"sha1", "keccakpg643", "keccak224", "skein256", "whirlpool", "blake2b",
}

type hf32 func(b []byte, seed uint32) uint32
type hf322 func(b []byte, l int, seeda, seedb uint32) (uint32, uint32)
type hf64 func(b []byte, seed uint64) uint64
type hf128e func(b []byte, seeda, seedb uint64) (uint64, uint64)

func hashspatch(de *DispEntry, b []byte, seed uint64) (ret uint64) {
	if de.kind == 2 {
		//pf := (*hf64)(de.fp)
		ret = (*(*hf64)(de.fp))(b, seed)
	} else if de.kind == 3 {
		pf := (*hf322)(de.fp)
		c, b := (*pf)(b, len(b), uint32(seed), uint32(seed>>32))
		ret = uint64(b)<<32 | uint64(c)
	} else if de.kind == 4 {
		//fmt.Printf("len(b)=%d\n", len(b))
		de.hi64.Reset()
		de.hi64.Write(b)
		ret = de.hi64.Sum64()
	} else {
		panic("hash")
	}
	return
}

func Halloc(hfs string) (hf32 nhash.HashF32) {
	switch hfs {
	case "sbox":
		hf32 = sbox.New(0)
	case "CrapWow":
		hf32 = crapwow.New(0)
	case "j332c":
		hf32 = jenkins.New332c(0)
	case "j232":
		hf32 = jenkins.New232(0)
	}
	return
}

// fingerprints are stored in these globals to take the allocation out of the loop
var fp8 = make([]byte, 8, 8)
var fp20 = make([]byte, 20, 20)
var fp28 = make([]byte, 28, 28)
var fp32 = make([]byte, 32, 32)
var fp64 = make([]byte, 64, 64)

var seeds []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

// crappy generic adapter that just slows us down
// will be removed
var sipSeedSet = func(seed uint64) {
	seeds[0], seeds[1], seeds[2], seeds[3], seeds[4], seeds[5], seeds[6], seeds[7] =
		byte(seed&0xFF), byte((seed>>8)&0xFF), byte((seed>>16)&0xFF), byte((seed>>24)&0xFF),
		byte((seed>>32)&0xFF), byte((seed>>40)&0xFF), byte((seed>>48)&0xFF), byte((seed>>56)&0xFF)
	seeds[8], seeds[9], seeds[10], seeds[11], seeds[12], seeds[13], seeds[14], seeds[15] = seeds[0], seeds[1], seeds[2], seeds[3], seeds[4], seeds[5], seeds[6], seeds[7]
}

var dis int = 0

/*
func Hashf(k []byte, seed uint64) uint64 {
	if dis == 0 {
		h := aeshash.Hash(k, seed)
		return h
	}
	return 0
}
*/

func Hashf(k []byte, seed uint64) (h uint64) {
	/*
		_, ok := HashFunctions[Hf2]
		if !ok {
			fmt.Printf("%q not found\n", Hf2)
			panic("hashf")
		}
	*/
	switch Hf2 {
	case "perfecthash":
		fmt.Printf("k=%v\n", k)
		//h = uint64(k[0])<<56 | uint64(k[1])<<48 | uint64(k[2])<<40 | uint64(k[3])<<32 | uint64(k[4])<<24 | uint64(k[5])<<16 | uint64(k[6])<<8 | uint64(k[7])<<0
		//h = uint64(k[7])<<56 | uint64(k[6])<<48 | uint64(k[5])<<40 | uint64(k[4])<<32 | uint64(k[3])<<24 | uint64(k[2])<<16 | uint64(k[1])<<8 | uint64(k[0])<<0
		h = uint64(k[3])<<24 | uint64(k[2])<<16 | uint64(k[1])<<8 | uint64(k[0])<<0
		//fmt.Printf("h=%d\n", h)
	case "nullhash":
		nh.Reset()
		nh.Write(k)
		h = nh.Sum64()
	case "nullhashF64ns":
		h = nhf64.Hash64S(k, 0)
	case "gomap64":
		h = gomap.Hash64(k, uint64(seed))
	case "gomap32":
		h = uint64(gomap.Hash32(k, uint32(seed)))
	case "aeshash64":
		//fmt.Printf("k=%v\n", k)
		h = aeshash.Hash(k, seed)
	case "aeshash32I":
		//fmt.Printf("k=%v\n", k)
		h = aeshash.Hash64(uint64(k[0])<<24|uint64(k[1])<<16|uint64(k[2])<<8|uint64(k[3])<<0, seed)
	case "aeshash64I":
		h = aeshash.Hash64(uint64(k[0])<<56|uint64(k[1])<<48|uint64(k[2])<<40|uint64(k[3])<<32|uint64(k[4])<<24|uint64(k[5])<<16|uint64(k[6])<<8|uint64(k[7])<<0, seed)
	case "adler32":
		a32.Reset()
		a32.Write(k)
		h = uint64(a32.Sum32())
		//fmt.Printf("a32 hash=0x%08x\n", h)
	case "sbox":
		h = uint64(sbox.Sbox(k, uint32(seed)))
	case "CrapWow":
		h = uint64(crapwow.CrapWow(k, uint32(seed)))
		//fmt.Printf("key=%q, hash=0x%08x\n", string(k), hash)
	case "MaHash8v64":
		h = mahash.MaHash8v64(k)
	case "j364":
		c, b := jenkins.Jenkins364(k, len(k), uint32(seed), uint32(seed))
		h = uint64(b)<<32 | uint64(c)
	case "j332c":
		c, _ := jenkins.Jenkins364(k, len(k), uint32(seed), uint32(seed))
		h = uint64(c)
	case "j332b":
		_, b := jenkins.Jenkins364(k, len(k), uint32(seed), uint32(seed))
		h = uint64(b)
	case "j232":
		h = uint64(jenkins.Hash232(k, uint32(seed)))
	case "j264":
		h = jenkins.Hash264(k, seed)
	case "j264l":
		h = uint64(jenkins.Hash264(k, seed) & 0xFFFFFFFF)
	case "j264h":
		t := jenkins.Hash264(k, seed)
		h = uint64((t >> 32) & 0xFFFFFFFF)
	case "j264xor":
		t := jenkins.Hash264(k, seed)
		h = uint64(uint32(t&0xFFFFFFFF) ^ uint32((t>>32)&0xFFFFFFFF))
	case "spooky32":
		h = uint64(spooky.Hash32(k, uint32(seed)))
	case "spooky64":
		h = spooky.Hash64(k, seed)
	case "spooky128h":
		h, _ = spooky.Hash128(k, seed)
	case "spooky128l":
		_, h = spooky.Hash128(k, seed)
	case "spooky128xor":
		t, l := spooky.Hash128(k, seed)
		h = t ^ l
	case "murmur332":
		m332.Reset()
		m332.Write(k)
		t := m332.Sum32()
		h = uint64(t)
	case "murmur364":
		m364.Reset()
		m364.Write(k)
		h = m364.Sum64()
	case "siphash64":
		h = siphash.Hash(0, 0, k)
	case "siphash64pg":
		sipSeedSet(seed)
		h, _ = siphashpg.Siphash(k, seeds, siphashpg.Crounds, siphashpg.Drounds, false)
	case "FarmHash32":
		t := farm.Fingerprint32(k)
		h = uint64(t & 0xFFFFFFFF)
	case "FarmHash64":
		h = farm.Fingerprint64(k)
	case "FarmHash128-high":
		h, _ = farm.Fingerprint128(k)
	case "FarmHash128-low":
		_, h = farm.Fingerprint128(k)
	case "FarmHash128-xor":
		t, l := farm.Fingerprint128(k)
		h = t ^ l
	case "MetroHash64-1":
		h = metro.Hash64_1(k, uint32(seed))
	case "MetroHash64-2":
		h = metro.Hash64_2(k, uint32(seed))
	case "MetroHash128-1h":
		h, _ = metro.Hash128_1(k, uint32(seed))
	case "MetroHash128-1l":
		_, h = metro.Hash128_1(k, uint32(seed))
	case "MetroHash128-1xor":
		t, l := metro.Hash128_1(k, uint32(seed))
		h = t ^ l
	case "MetroHash128-2h":
		h, _ = metro.Hash128_2(k, uint32(seed))
	case "MetroHash128-2l":
		_, h = metro.Hash128_2(k, uint32(seed))
	case "MetroHash128-2xor":
		t, l := metro.Hash128_2(k, uint32(seed))
		h = t ^ l

		/*
			case "siphash64":
				sipSeedSet(seed)
				a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
				return a
			case "siphash128a":
				sipSeedSet(seed)
				a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
				return a
			case "siphash128b":
				sipSeedSet(seed)
				_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
				return b
			case "siphash64al":
				sipSeedSet(seed)
				a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
				return uint64(a&0xFFFFFFFF)
			case "siphash64ah":
				sipSeedSet(seed)
				a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
				return uint64((a>>32)&0xFFFFFFFF)
			case "siphash64bl":
				sipSeedSet(seed)
				_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
				return uint64(b&0xFFFFFFFF)
			case "siphash64bh":
				sipSeedSet(seed)
				_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
				return uint64((b>>32)&0xFFFFFFFF)
		*/
	case "keccak224":
		fp28 = fp28[0:0]
		k224.Reset()
		//fmt.Printf("len(k)=%d\n", k)
		k224.Write(k)
		fp28 = k224.Sum(fp28) // crashes 		fp8 = k224.Sum(fp28)
		//fmt.Printf("len(k)=%d, k=%#X, fp=%#x\n", len(k), k, fp28)
		//fmt.Printf("len(fp28)=%d\n", fp28)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		h = uint64(fp28[0])<<56 | uint64(fp28[1])<<48 | uint64(fp28[2])<<40 | uint64(fp28[3])<<32 | uint64(fp28[4])<<24 | uint64(fp28[5])<<16 | uint64(fp28[6])<<8 | uint64(fp28[7])<<0
	case "keccakpg643":
		fp8 = fp8[0:0]
		k643.Reset()
		k643.Write(k)
		fp8 = k643.Sum(fp8)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		h = uint64(fp8[0])<<56 | uint64(fp8[1])<<48 | uint64(fp8[2])<<40 | uint64(fp8[3])<<32 | uint64(fp8[4])<<24 | uint64(fp8[5])<<16 | uint64(fp8[6])<<8 | uint64(fp8[7])<<0
	case "keccakpg644":
		fp := make([]byte, 8, 8)
		fp = fp[0:0]
		k644.Reset()
		k644.Write(k)
		fp = k644.Sum(fp)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		h = uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8 | uint64(fp[7])<<0
	case "keccakpg648":
		fp := make([]byte, 8, 8)
		fp = fp[0:0]
		k648.Reset()
		k648.Write(k)
		fp = k648.Sum(fp)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		h = uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8 | uint64(fp[7])<<0
	case "keccakpg160":
		fp32 = fp32[0:0]
		k160.Reset()
		k160.Write(k)
		fp32 = k160.Sum(fp32)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		if false {
			low := fp32[0] ^ fp32[4] ^ fp32[8] ^ fp32[12] ^ fp32[16]
			med := fp32[1] ^ fp32[5] ^ fp32[9] ^ fp32[13] ^ fp32[17]
			hii := fp32[2] ^ fp32[6] ^ fp32[10] ^ fp32[14] ^ fp32[18]
			top := fp32[3] ^ fp32[7] ^ fp32[11] ^ fp32[15] ^ fp32[19]
			h = uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
			h = uint64(fp32[0])<<56 | uint64(fp32[1])<<48 | uint64(fp32[2])<<40 | uint64(fp32[3])<<32 |
				uint64(fp32[4])<<24 | uint64(fp32[5])<<16 | uint64(fp32[6])<<8 | uint64(fp32[7])<<0
		}
	case "skein256xor":
		fp32 = fp32[0:0]
		skein256.Reset()
		skein256.Write(k)
		fp32 = skein256.Sum(fp32)
		//fmt.Printf("skein256: fp=%v\n", fp)
		if true {
			low := fp32[0] ^ fp32[4] ^ fp32[8] ^ fp32[12] ^ fp32[16]
			med := fp32[1] ^ fp32[5] ^ fp32[9] ^ fp32[13] ^ fp32[17]
			hii := fp32[2] ^ fp32[6] ^ fp32[10] ^ fp32[14] ^ fp32[18]
			top := fp32[3] ^ fp32[7] ^ fp32[11] ^ fp32[15] ^ fp32[19]
			h = uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
			h = uint64(fp32[0])<<56 | uint64(fp32[1])<<48 | uint64(fp32[2])<<40 | uint64(fp32[3])<<32 |
				uint64(fp32[4])<<24 | uint64(fp32[5])<<16 | uint64(fp32[6])<<8 | uint64(fp32[7])<<0
		}
	case "skein256":
		fp32 = fp32[0:0]
		skein256.Reset()
		skein256.Write(k)
		fp32 = skein256.Sum(fp32)
		//fmt.Printf("skein256: fp=%v\n", fp)
		h = uint64(fp32[0])<<56 | uint64(fp32[1])<<48 | uint64(fp32[2])<<40 | uint64(fp32[3])<<32 |
			uint64(fp32[4])<<24 | uint64(fp32[5])<<16 | uint64(fp32[6])<<8 | uint64(fp32[7])<<0
	case "skein256hi":
		fp32 = fp32[0:0]
		skein256.Reset()
		skein256.Write(k)
		fp32 = skein256.Sum(fp32)
		//fmt.Printf("skein256: fp=%v\n", fp)
		h = uint64(fp32[0])<<56 | uint64(fp32[1])<<48 | uint64(fp32[2])<<40 | uint64(fp32[3])<<32 |
			uint64(fp32[4])<<24 | uint64(fp32[5])<<16 | uint64(fp32[6])<<8 | uint64(fp32[7])<<0
	case "whirlpool":
		fp64 = fp64[0:0]
		wp.Reset()
		wp.Write(k)
		fp64 = wp.Sum(fp64)
		//fmt.Printf("whirlpool: fp=%v\n", fp)
		h = uint64(fp64[0])<<56 | uint64(fp64[1])<<48 | uint64(fp64[2])<<40 | uint64(fp64[3])<<32 |
			uint64(fp64[4])<<24 | uint64(fp64[5])<<16 | uint64(fp64[6])<<8 | uint64(fp64[7])<<0
		//fmt.Printf("whirlpool: h=%#016x\n", h)
	case "blake2b":
		fp32 = fp32[0:0]
		b2b.Reset()
		b2b.Write(k)
		fp32 = b2b.Sum(fp32)
		//fmt.Printf("blake2b: fp=%v\n", fp)
		h = uint64(fp32[0])<<56 | uint64(fp32[1])<<48 | uint64(fp32[2])<<40 | uint64(fp32[3])<<32 |
			uint64(fp32[4])<<24 | uint64(fp32[5])<<16 | uint64(fp32[6])<<8 | uint64(fp32[7])<<0
		//fmt.Printf("blake2b: h=%#016x\n", h)
	case "sha1":
		//fp := make([]byte, 20)
		//fp = fp[0:0]
		sha1160.Reset()
		sha1160.Write(k)
		fp20 = fp20[0:0]
		fp20 = sha1160.Sum(fp20)
		if false {
			low := fp20[0] ^ fp20[4] ^ fp20[8] ^ fp20[12] ^ fp20[16]
			med := fp20[1] ^ fp20[5] ^ fp20[9] ^ fp20[13] ^ fp20[17]
			hii := fp20[2] ^ fp20[6] ^ fp20[10] ^ fp20[14] ^ fp20[18]
			top := fp20[3] ^ fp20[7] ^ fp20[11] ^ fp20[15] ^ fp20[19]
			h = uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
			h = uint64(fp20[0])<<56 | uint64(fp20[1])<<48 | uint64(fp20[2])<<40 | uint64(fp20[3])<<32 | uint64(fp20[4])<<24 | uint64(fp20[5])<<16 | uint64(fp20[6])<<8 | uint64(fp20[7])<<0
		}
	default:
		fmt.Printf("hf=%q\n", Hf2)
		panic("hashf")
	}
	return
}

func init() {
	nh = nullhash.New()
	nhf64 = nullhash.NewF64()
	m332 = murmur3.New32()
	m364 = murmur3.New64()
	k643 = keccakpg.NewCustom(64, 3)
	k644 = keccakpg.NewCustom(64, 4)
	k648 = keccakpg.NewCustom(64, 8)
	k160 = keccakpg.New160()
	k224 = keccak.New224()
	skein256 = skein.New256()
	//skein32 := skein.New(256, 32)
	sha1160 = sha1.New()
	a32 = adler32.New()
	wp = whirlpool.New()
	b2b = blake2b.New256()
	//HashFunctions["keccak224"].hf = k224
}
