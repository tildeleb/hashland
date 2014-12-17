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
	_ "github.com/tildeleb/hashland/siphashpg"
	"github.com/tildeleb/hashland/siphash"
	"github.com/tildeleb/hashland/keccak"
	"github.com/tildeleb/hashland/skein"
	"github.com/tildeleb/hashland/aeshash"
	//"github.com/tildeleb/hashland/threefish"
)

var a32 hash.Hash32
var k643 hash.Hash
var k644 hash.Hash
var k648 hash.Hash
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
	"aeshash64":		HashFunction{"aeshash64", 		64,		true,	"aeshash, 64 bit, accelerated"},
	"siphash64":		HashFunction{"siphash64", 		64,		true,	"siphash, 64 bit, accelerated"},
/*
	"siphash64":		HashFunction{"siphash64", 		64,		true,	"siphash, 64 bit, a bits"},
	"siphash128a":		HashFunction{"siphasha", 		64,		true,	"siphash, 128 bit, a bits"},
	"siphash128b":		HashFunction{"siphashb", 		64,		true,	"siphash, 128 bit, b bits"},
	"siphash64al":		HashFunction{"siphash64al", 	32,		true,	"siphash, 64 bit, a bits, low"},
	"siphash64ah":		HashFunction{"siphash64ah", 	32,		true,	"siphash, 64 bit, a bits, high"},
	"siphash64bl":		HashFunction{"siphash64bl", 	32,		true,	"siphash, 128 bit, b bits, low"},
	"siphash64bh":		HashFunction{"siphash64bh", 	32,		true,	"siphash, 128 bit, b bits, high"},
*/
	"MaHash8v64":		HashFunction{"MaHash8v64", 		64,		false,	"russian hash function"},

	// tribute to Robert Jenkins goes here
	"spooky32":			HashFunction{"spooky32", 		32,		false,	"jenkins, spooky, 32 bit"},
	"spooky64":			HashFunction{"spooky64", 		64,		false,	"jenkins, spooky, 64 bit"},
	"spooky128h":		HashFunction{"spooky128h", 		64,		false,	"jenkins, spooky, 128 bit, high bits"},
	"spooky128l":		HashFunction{"spooky128l", 		64,		false,	"jenkins, spooky, 128 bit, low bits"},
	"spooky128xor":		HashFunction{"spooky128xor",	64,		false,	"jenkins, spooky, 128, high xor low bits"},
	"j364":				HashFunction{"j364", 			64,		false,	"jenkins, lookup3. 64 bit, c low order bits, b high order bits"},
	"j264":				HashFunction{"j264", 			64,		false,	"jenkins, lookup8. 64 bit"},
	"j332c":			HashFunction{"j332c", 			32,		false,	"jenkins, lookup3, 32 bit, c bits"},
	"j332b":			HashFunction{"j332b", 			32,		false,	"jenkins, lookup3, 32 bit, b bits"},
	"j232":				HashFunction{"j232", 			32,		false,	"jenkins, lookup8, 32 bit"},
	"j264l":			HashFunction{"j264l", 			32,		false,	"jenkins, lookup8, 64 bit, low bits"},
	"j264h":			HashFunction{"j264h", 			32,		false,	"jenkins, lookup8, 64 bit, high bits"},
	"j264xor":			HashFunction{"j264xor",			32,		false,	"jenkins, lookup8, 64 bit, high xor low bits"},
	"sbox":				HashFunction{"sbox", 			32,		false,	"sbox"},

	"keccak643":		HashFunction{"keccak643", 		64,		true,	"keccak, 64 bit, 3 rounds"},
	"keccak644":		HashFunction{"keccak644", 		64,		true,	"keccak, 64 bit, 4 rounds"},
	"keccak648":		HashFunction{"keccak648", 		64,		true,	"keccak, 64 bit, 8 rounds"},
	"skein256":			HashFunction{"skein256", 		64,		true,	"skein256, 64 bit , low 64 bits"},
	"sha1":				HashFunction{"sha1", 			64,		true,	"sha1, 160 bit hash, 64 bit, low 64 bits"},
	"keccak160":		HashFunction{"keccak160", 		64,		true,	"keccak160l"},

	"skein256low":		HashFunction{"skein256low", 	32,		true,	"skein256low"},
	"skein256hi":		HashFunction{"skein256hi", 		32,		true,	"skein256hi"},
	"skein256xor":		HashFunction{"skein256xor", 	32,		true,	"skein256xor"},


	"CrapWow":			HashFunction{"CrapWow", 		32,		false,	"CrapWow"},
	"adler32":			HashFunction{"adler32", 		32,		false,	"adler32"},
}

// "CrapWow" removed because it generates some many dup hashes with duplicated words it goes from O(1) to O(N)
// "adler32" removed for the same reasons
// 	"siphash64al", "siphash64ah", "siphash64bl", "siphash64bh",
// 	"skein256xor", "skein256low", "skein256hi", "sha1", "keccak160l", 
// 	"siphash64", "siphash128a", "siphash128b",
// 	"keccak644", "keccak648" "keccak160", 
var TestHashFunctions = []string{"aeshash64","j364", "j264",
	"siphash64",
	"MaHash8v64", "spooky64", "spooky128h", "spooky128l", "spooky128xor",
	"j332c", "j332b", "j232", "j264l", "j264h", "j264xor", "spooky32",  "sbox",
	"sha1", "keccak643", "skein256",
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

var seeds []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
// crappy generic adapter that just slows us down
// will be removed
func Hashf(k []byte, seed uint64) uint64 {
/*
	var sipSeedSet = func(seed uint64) {
		seeds[0], seeds[1], seeds[2], seeds[3], seeds[4], seeds[5], seeds[6], seeds[7] =
			byte(seed&0xFF), byte((seed>>8)&0xFF), byte((seed>>16)&0xFF), byte((seed>>24)&0xFF),
			byte((seed>>32)&0xFF), byte((seed>>40)&0xFF), byte((seed>>48)&0xFF), byte((seed>>56)&0xFF)
		seeds[8], seeds[9], seeds[10], seeds[11], seeds[12], seeds[13], seeds[14], seeds[15] = seeds[0], seeds[1], seeds[2], seeds[3], seeds[4], seeds[5], seeds[6], seeds[7]
	}
*/
/*
	_, ok := HashFunctions[Hf2]
	if !ok {
		fmt.Printf("%q not found\n", Hf2)
		panic("hashf")
	}
*/
	switch Hf2 {
	case "aeshash64":
		h := aeshash.Hash(k, seed)
		return h
	case "adler32":
		a32.Reset()
		a32.Write(k)
		h := a32.Sum32()
		//fmt.Printf("a32 hash=0x%08x\n", h)
		return uint64(h)
	case "sbox":
		h := sbox.Sbox(k, uint32(seed))
		return uint64(h)
	case "CrapWow":
		h := crapwow.CrapWow(k, uint32(seed))
		//fmt.Printf("key=%q, hash=0x%08x\n", string(k), hash)
		return uint64(h)
	case "MaHash8v64":
		h64 := mahash.MaHash8v64(k)
		return h64
	case "j364":
		c, b := jenkins.Jenkins364(k, len(k), uint32(seed), uint32(seed))
		return uint64(b)<<32 | uint64(c)
	case "j332c":
		c, _ := jenkins.Jenkins364(k, len(k), uint32(seed), uint32(seed))
		return uint64(c)
	case "j332b":
		_, b := jenkins.Jenkins364(k, len(k), uint32(seed), uint32(seed))
		return uint64(b)
	case "j232":
		h := jenkins.Hash232(k, uint32(seed))
		return uint64(h)
	case "j264":
		h := jenkins.Hash264(k, seed)
		return h
	case "j264l":
		h := jenkins.Hash264(k, seed)
		return uint64(h&0xFFFFFFFF)
	case "j264h":
		h := jenkins.Hash264(k, seed)
		return uint64((h>>32)&0xFFFFFFFF)
	case "j264xor":
		h := jenkins.Hash264(k, seed)
		return uint64(uint32(h&0xFFFFFFFF) ^ uint32((h>>32)&0xFFFFFFFF))
	case "spooky32":
		return uint64(spooky.Hash32(k, uint32(seed)))
	case "spooky64":
		return spooky.Hash64(k, seed)
	case "spooky128h":
		h, _ := spooky.Hash128(k, seed)
		return h
	case "spooky128l":
		_, l := spooky.Hash128(k, seed)
		return l
	case "spooky128xor":
		h, l := spooky.Hash128(k, seed)
		return h ^ l
	case "siphash64":
		h := siphash.Hash(0, 0, k)
		return h
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
	case "keccak643":
		fp := make([]byte, 8, 8)
		fp = fp[0:0]
		k643.Reset()
		k643.Write(k)
		fp = k643.Sum(fp)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		return uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8  | uint64(fp[7])<<0 
	case "keccak644":
		fp := make([]byte, 8, 8)
		fp = fp[0:0]
		k644.Reset()
		k644.Write(k)
		fp = k644.Sum(fp)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		return uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8  | uint64(fp[7])<<0 
	case "keccak648":
		fp := make([]byte, 8, 8)
		fp = fp[0:0]
		k648.Reset()
		k648.Write(k)
		fp = k648.Sum(fp)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		return uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8  | uint64(fp[7])<<0 
	case "keccak160":
		fp := make([]byte, 32)
		fp = fp[0:0]
		k160.Reset()
		k160.Write(k)
		fp = k160.Sum(fp)
		//fmt.Printf("keccak160xor: fp=%v\n", fp)
		if false {
	        low := fp[0] ^ fp[4] ^ fp[8] ^ fp[12] ^ fp[16]
	        med := fp[1] ^ fp[5] ^ fp[9] ^ fp[13] ^ fp[17]
	        hii := fp[2] ^ fp[6] ^ fp[10] ^ fp[14] ^ fp[18]
	        top := fp[3] ^ fp[7] ^ fp[11] ^ fp[15] ^ fp[19]
	        return uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
			return uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8  | uint64(fp[7])<<0 
		}
	case "skein256xor":
		fp := make([]byte, 32)
		fp = fp[0:0]
		skein256.Reset()
		skein256.Write(k)
		fp = skein256.Sum(fp)
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
	case "skein256":
		fp := make([]byte, 32)
		fp = fp[0:0]
		skein256.Reset()
		skein256.Write(k)
		fp = skein256.Sum(fp)
		//fmt.Printf("skein256: fp=%v\n", fp)
		return uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8  | uint64(fp[7])<<0 
	case "skein256hi":
		fp := make([]byte, 32)
		fp = fp[0:0]
		skein256.Reset()
		skein256.Write(k)
		fp = skein256.Sum(fp)
		//fmt.Printf("skein256: fp=%v\n", fp)
    	return uint64(uint32(fp[28])<<24 | uint32(fp[29])<<16 | uint32(fp[30])<<8 | uint32(fp[31]))
	case "sha1":
		fp := make([]byte, 32)
		fp = fp[0:0]
		sha1160.Reset()
		sha1160.Write(k)
		fp = sha1160.Sum(fp)
		if false {
	        low := fp[0] ^ fp[4] ^ fp[8] ^ fp[12] ^ fp[16]
	        med := fp[1] ^ fp[5] ^ fp[9] ^ fp[13] ^ fp[17]
	        hii := fp[2] ^ fp[6] ^ fp[10] ^ fp[14] ^ fp[18]
	        top := fp[3] ^ fp[7] ^ fp[11] ^ fp[15] ^ fp[19]
        	return uint64(uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low))
		} else {
			return uint64(fp[0])<<56 | uint64(fp[1])<<48 | uint64(fp[2])<<40 | uint64(fp[3])<<32 | uint64(fp[4])<<24 | uint64(fp[5])<<16 | uint64(fp[6])<<8  | uint64(fp[7])<<0 
		}
	default:
		fmt.Printf("hf=%q\n", Hf2)
		panic("hashf")
	}
	return 0
}

func init() {
	k643 = keccak.NewCustom(64, 3)
	k644 = keccak.NewCustom(64, 4)
	k648 = keccak.NewCustom(64, 8)
	k160 = keccak.New160()
	skein256 = skein.New256()
	//skein32 := skein.New(256, 32)
	sha1160 = sha1.New()
	a32 = adler32.New()
}
