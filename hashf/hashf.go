// Copyright © 2014 Lawrence E. Bakst. All rights reserved.

package hashf

import (
	"fmt"
	"hash"
	"hash/adler32"
	"crypto/sha1"
	"github.com/tildeleb/hashland/nhash"
	"github.com/tildeleb/hashland/sbox"
	"github.com/tildeleb/hashland/crapwow"
	"github.com/tildeleb/hashland/jenkins"
	"github.com/tildeleb/hashland/mahash"
	"github.com/tildeleb/hashland/spooky"
	"github.com/tildeleb/hashland/siphash"
	"github.com/tildeleb/hashland/keccak"
	"github.com/tildeleb/hashland/skein"
	//"github.com/tildeleb/hashland/threefish"
)

var a32 hash.Hash32
var k160 hash.Hash
var skein256 hash.Hash
var sha1160 hash.Hash

var Hf2 string // wow this has to go

type HashFunction struct {
	Name		string
	Size		int // in bits
	Crypto		bool
	desc		string
}


var HashFunctions = map[string]HashFunction{
	"siphash64":		HashFunction{"siphash64", 		64,		true,	"siphash, 64 bit, a bits"},
	"siphash128a":		HashFunction{"siphasha", 		64,		true,	"siphash, 128 bit, a bits"},
	"siphash128b":		HashFunction{"siphashb", 		64,		true,	"siphash, 128 bit, b bits"},
	"siphash64al":		HashFunction{"siphash64al", 	32,		true,	"siphash, 64 bit, a bits, low"},
	"siphash64ah":		HashFunction{"siphash64ah", 	32,		true,	"siphash, 64 bit, a bits, high"},
	"siphash64bl":		HashFunction{"siphash64bl", 	32,		true,	"siphash, 128 bit, b bits, low"},
	"siphash64bh":		HashFunction{"siphash64bh", 	32,		true,	"siphash, 128 bit, b bits, high"},

	"MaHash8v64":		HashFunction{"MaHash8v64", 		64,		false,	"russian hash function"},

	// tribute to Robert Jenkins goes here
	"spooky32":			HashFunction{"spooky32", 		32,		false,	"jenkins, spooky, 32 bit"},
	"spooky64":			HashFunction{"spooky64", 		64,		false,	"jenkins, spooky, 64 bit"},
	"spooky128h":		HashFunction{"spooky128h", 		64,		false,	"jenkins, spooky, 128 bit, high bits"},
	"spooky128l":		HashFunction{"spooky128l", 		64,		false,	"jenkins, spooky, 128 bit, low bits"},
	"spooky128xor":		HashFunction{"spooky128xor",	64,		false,	"jenkins, spooky, 128, high xor low bits"},
	"j264":				HashFunction{"j264", 			64,		false,	"jenkins, lookup8. 64 bit"},
	"j332c":			HashFunction{"j332c", 			32,		false,	"jenkins, lookup3, 32 bit, c bits"},
	"j332b":			HashFunction{"j332b", 			32,		false,	"jenkins, lookup3, 32 bit, b bits"},
	"j232":				HashFunction{"j232", 			32,		false,	"jenkins, lookup8, 32 bit"},
	"j264l":			HashFunction{"j264l", 			32,		false,	"jenkins, lookup8, 64 bit, low bits"},
	"j264h":			HashFunction{"j264h", 			32,		false,	"jenkins, lookup8, 64 bit, high bits"},
	"j264xor":			HashFunction{"j264xor",			32,		false,	"jenkins, lookup8, 64 bit, high xor low bits"},

	"sbox":				HashFunction{"sbox", 			32,		false,	"sbox"},
	"skein256low":		HashFunction{"skein256low", 	32,		true,	"skein256low"},
	"skein256hi":		HashFunction{"skein256hi", 		32,		true,	"skein256hi"},
	"skein256xor":		HashFunction{"skein256xor", 	32,		true,	"skein256xor"},
	"sha1":				HashFunction{"sha1", 			32,		true,	"sha1"},
	"keccak160l":		HashFunction{"keccak160l", 		32,		true,	"keccak160l"},

	"CrapWow":			HashFunction{"CrapWow", 		32,		false,	"CrapWow"},
	"adler32":			HashFunction{"adler32", 		32,		false,	"adler32"},
}

// "CrapWow" removed because it generates some many dup hashes with duplicated words it goes from O(1) to O(N)
// "adler32" removed for the same reasons
var TestHashFunctions = []string{"j264", "siphash128a", "siphash128b", "MaHash8v64", "spooky64", "spooky128h", "spooky128l", "spooky128xor", "sbox",
	"j332c", "j332b", "j232", "j264l", "j264h", "j264xor", "spooky32",
	"siphash64al", "siphash64ah", "siphash64bl", "siphash64bh",
	"skein256xor", "skein256low", "skein256hi", "sha1", "keccak160l", 
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

// crappy generic adapter that just slows us down
// will be removed
func Hashf(k []byte) uint64 {
	var seeds []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	var fp = make([]byte, 32)
	_, ok := HashFunctions[Hf2]
	if !ok {
		fmt.Printf("%q not found\n", Hf2)
		panic("hashf")
	}
	switch Hf2 {
	case "adler32":
		a32.Reset()
		a32.Write(k)
		h := a32.Sum32()
		//fmt.Printf("a32 hash=0x%08x\n", h)
		return uint64(h)
	case "sbox":
		h := sbox.Sbox(k, 0)
		return uint64(h)
	case "CrapWow":
		h := crapwow.CrapWow(k, 0)
		//fmt.Printf("key=%q, hash=0x%08x\n", string(k), hash)
		return uint64(h)
	case "MaHash8v64":
		h64 := mahash.MaHash8v64(k)
		return h64
	case "j332c":
		c, _ := jenkins.Jenkins364(k, len(k), 0, 0)
		return uint64(c)
	case "j332b":
		_, b := jenkins.Jenkins364(k, len(k), 0, 0)
		return uint64(b)
	case "j232":
		h := jenkins.Hash232(k, 0)
		return uint64(h)
	case "j264":
		h := jenkins.Hash264(k, 0)
		return h
	case "j264l":
		h := jenkins.Hash264(k, 0)
		return uint64(h&0xFFFFFFFF)
	case "j264h":
		h := jenkins.Hash264(k, 0)
		return uint64((h>>32)&0xFFFFFFFF)
	case "j264xor":
		h := jenkins.Hash264(k, 0)
		return uint64(uint32(h&0xFFFFFFFF) ^ uint32((h>>32)&0xFFFFFFFF))
	case "spooky32":
		return uint64(spooky.Hash32(k, 0))
	case "spooky64":
		return spooky.Hash64(k, 0)
	case "spooky128h":
		h, _ := spooky.Hash128(k, 0)
		return h
	case "spooky128l":
		_, l := spooky.Hash128(k, 0)
		return l
	case "spooky128xor":
		h, l := spooky.Hash128(k, 0)
		return h ^ l
	case "siphash64":
		a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
		return a
	case "siphash128a":
		a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
		return a
	case "siphash128b":
		_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
		return b
	case "siphash64al":
		a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
		return uint64(a&0xFFFFFFFF)
	case "siphash64ah":
		a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
		return uint64((a>>32)&0xFFFFFFFF)
	case "siphash64bl":
		_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
		return uint64(b&0xFFFFFFFF)
	case "siphash64bh":
		_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
		return uint64((b>>32)&0xFFFFFFFF)
	case "keccak160l":
		k160.Reset()
		k160.Write(k)
		fp = fp[0:0]
		fp := k160.Sum(fp)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		if false {
	        low := fp[0] ^ fp[4] ^ fp[8] ^ fp[12] ^ fp[16]
	        med := fp[1] ^ fp[5] ^ fp[9] ^ fp[13] ^ fp[17]
	        hii := fp[2] ^ fp[6] ^ fp[10] ^ fp[14] ^ fp[18]
	        top := fp[3] ^ fp[7] ^ fp[11] ^ fp[15] ^ fp[19]
	        return uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
			return uint64(uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3]))
		}
	case "skein256xor":
		skein256.Reset()
		skein256.Write(k)
		fp = fp[0:0]
		fp := skein256.Sum(fp)
		//fmt.Printf("skein256: fp=%v\n", fp)
		if true {
	        low := fp[0] ^ fp[4] ^ fp[8] ^ fp[12] ^ fp[16]
	        med := fp[1] ^ fp[5] ^ fp[9] ^ fp[13] ^ fp[17]
	        hii := fp[2] ^ fp[6] ^ fp[10] ^ fp[14] ^ fp[18]
	        top := fp[3] ^ fp[7] ^ fp[11] ^ fp[15] ^ fp[19]
	        return uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
	    	return uint64(uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3]))
	    }
	case "skein256low":
		skein256.Reset()
		skein256.Write(k)
		fp = fp[0:0]
		fp := skein256.Sum(fp)
		//fmt.Printf("skein256: fp=%v\n", fp)
    	return uint64(uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3]))
	case "skein256hi":
		skein256.Reset()
		skein256.Write(k)
		fp = fp[0:0]
		fp := skein256.Sum(fp)
		//fmt.Printf("skein256: fp=%v\n", fp)
    	return uint64(uint32(fp[28])<<24 | uint32(fp[29])<<16 | uint32(fp[30])<<8 | uint32(fp[31]))
	case "sha1":
		sha1160.Reset()
		sha1160.Write(k)
		fp = fp[0:0]
		fp := sha1160.Sum(fp)
		if false {
	        low := fp[0] ^ fp[4] ^ fp[8] ^ fp[12] ^ fp[16]
	        med := fp[1] ^ fp[5] ^ fp[9] ^ fp[13] ^ fp[17]
	        hii := fp[2] ^ fp[6] ^ fp[10] ^ fp[14] ^ fp[18]
	        top := fp[3] ^ fp[7] ^ fp[11] ^ fp[15] ^ fp[19]
        	return uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
			return uint64(uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3]))
		}
		default:
			fmt.Printf("hf=%q\n", Hf2)
			panic("hashf")
	}
	return 0
}

func init() {
	k160 = keccak.New160()
	skein256 = skein.New256()
	//skein32 := skein.New(256, 32)
	sha1160 = sha1.New()
	a32 = adler32.New()
}
