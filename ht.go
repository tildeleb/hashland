// Copyright © 2014 Lawrence E. Bakst. All rights reserved.
package main

// based on http://amsoftware.narod.ru/algo.html
// However I can't replicate the result

import (
	"flag"
	"fmt"
	"os"
	"io"
	"bufio"
	"hash"
	"hash/adler32"
	"github.com/tildeleb/hashland/nhash"
	"sort"
	"time"
	"crypto/sha1"
	"github.com/tildeleb/hashland/sbox"
	"github.com/tildeleb/hashland/crapwow"
	"github.com/tildeleb/hashland/jenkins"
	"github.com/tildeleb/hashland/mahash"
	"github.com/tildeleb/hashland/spooky"
	"github.com/tildeleb/hashland/siphash"
	"github.com/tildeleb/hashland/keccak"
	"github.com/tildeleb/hashland/skein"
	//"github.com/tildeleb/hashland/threefish"
	"github.com/tildeleb/cuckoo/primes"
	"github.com/tildeleb/hrff"
)

var a32 hash.Hash32
var k160 hash.Hash
var skein256 hash.Hash
var sha1160 hash.Hash

var hf2 string

type Bucket struct {
	Key []byte
}

type Stats struct {
	Inserts int
	Cols int
	Probes int
	Heads int
	Dups int
	Nbuckets int
	Entries int
	Q float64
	//
	Lines int
	Size uint64
	SizeLog2 uint64
	SizeMask uint64
}

type HashTable struct {
	Buckets [][]Bucket
	Stats
}

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
var hashFunctions = []string{"j264", "siphash128a", "siphash128b", "MaHash8v64", "spooky64", "spooky128h", "spooky128l", "spooky128xor", "sbox",
	"j332c", "j332b", "j232", "j264l", "j264h", "j264xor", "spooky32",
	"siphash64al", "siphash64ah", "siphash64bl", "siphash64bh",
	"skein256xor", "skein256low", "skein256hi", "sha1", "keccak160l", 
}

// crappy generic adapter that just slows us down
// will be removed
func hashf(k []byte) uint64 {
	var seeds []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	var fp = make([]byte, 32)
	_, ok := HashFunctions[hf2]
	if !ok {
		fmt.Printf("%q not found\n", hf2)
		panic("hashf")
	}
	switch hf2 {
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
			fmt.Printf("hf=%q\n", hf2)
			panic("hashf")
	}
	return 0
}

func NewHashTable(lines int) *HashTable {
	ht := new(HashTable)
	ht.Lines = lines
	ht.SizeLog2 = uint64(NextLog2(uint32(ht.Lines)) + uint32(*extra))
	ht.Size = 1 << ht.SizeLog2
	if *prime {
		ht.Size = uint64(primes.NextPrime(int(ht.Size)))
	}
	ht.SizeMask = ht.Size - 1
	ht.Buckets = make([][]Bucket, ht.Size, ht.Size)
	return ht
}

func (ht *HashTable) Insert(ka []byte) {
	k := make([]byte, len(ka), len(ka))
	k = k[:]
	amt := copy(k, ka)
	if amt != len(ka) {
		panic("Add")
	}
	ht.Inserts++
	idx := uint64(0)
	h := hashf(k) // jenkins.Hash232(k, 0)
	if *prime {
		idx = h % ht.Size
	} else {
		idx = h & ht.SizeMask
	}
	//fmt.Printf("index=%d\n", idx)
	cnt := 0
	pass := 0

	//fmt.Printf("Add: %x\n", k)
	//ht.Buckets[idx].Key = k
	//len(ht.Buckets[idx].Key) == 0
	for {
		if ht.Buckets[idx] == nil {
			// no entry or chain at this location, make it
			ht.Buckets[idx] = append(ht.Buckets[idx], Bucket{Key: k})
			//fmt.Printf("Add: idx=%d, len=%d, hash=0x%08x, key=%q\n", idx, len(ht.Buckets[idx]), h, ht.Buckets[idx][0].Key)
			ht.Probes++
			ht.Heads++
			return
		}
		if *oa {
			if cnt == 0 {
				ht.Cols++
			} else {
				ht.Probes++
			}

			// check for a duplicate key
			bh := hashf(ht.Buckets[idx][0].Key)
			if bh == h {
				if *pd {
					fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", h, idx, k)
					fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", bh, idx, ht.Buckets[idx][0].Key)
				}
				ht.Dups++
			}
			idx++
			cnt++
			if idx > ht.Size - 1 {
				pass++
				if pass > 1 {
					panic("Add: pass")
				}
				idx = 0
			}
		} else {
			// first scan slice for dups
			for j := range ht.Buckets[idx] {
				bh := hashf(ht.Buckets[idx][j].Key)
				//fmt.Printf("idx=%d, j=%d/%d, bh=0x%08x, h=0x%08x, key=%q\n", idx, j, len(ht.Buckets[idx]), bh, h, ht.Buckets[idx][j].Key)
				if bh == h {
					if *pd {
						fmt.Printf("idx=%d, j=%d/%d, bh=0x%08x, h=0x%08x, key=%q, bkey=%q\n", idx, j, len(ht.Buckets[idx]), bh, h, k, ht.Buckets[idx][j].Key)
						//fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", h, idx, k)
						//fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", bh, idx, ht.Buckets[idx][0].Key)
					}
					ht.Dups++
				}
			}
			// add element
			ht.Buckets[idx] = append(ht.Buckets[idx], Bucket{Key: k})
			ht.Probes++
			break
		}
	}
}

// The theoretical metric from "Red Dragon Book"
func (ht *HashTable) HashQuality() float64 {
	n := float64(0.0)
	buckets := 0
	entries := 0
	for _, v := range ht.Buckets {
		if v != nil {
			buckets++
			count := float64(len(v))
			entries += len(v)
			n += count * (count + 1.0)
		}
	}
	n *= float64(ht.Size)
	d := float64(ht.Inserts) * (float64(ht.Inserts) + 2.0 * float64(ht.Size) - 1.0) 	// (n / 2m) * (n + 2m - 1)
	//fmt.Printf("buckets=%d, entries=%d, inserts=%d, size=%d, n=%f, d=%f, n/d=%f\n", buckets, entries, ht.Inserts, ht.Size, n, d, n/d)
	ht.Nbuckets = buckets
	ht.Entries = entries
	ht.Q = n / d
	return n / d
}

// Henry Warren, "Hacker's Delight", ch. 5.3
func NextLog2(x uint32) uint32 {
	if x <= 1 {
		return x
	} 
	x--
	n := uint32(0)
	y := uint32(0)
	y = x >>16
	if y != 0 {
		n += 16
		x = y
	}
	y = x >> 8
	if y != 0 {
		n += 8
		x = y
	}
	y = x >> 4;
	if y != 0 {
		n +=  4
		x = y
	}
	y = x >> 2
	if y != 0 {
		n +=  2
		x = y
	}
	y = x >> 1
	if y != 0 {
		return n + 2
	}
	return n + x
}

func ReadFile(file string, cb func(line string)) int {
	var lines int
	//fmt.Printf("ReadFile: file=%q\n", file)
	f, err := os.Open(file)
    if err != nil {
        panic("ReadFile: opening file")
    }
    defer f.Close()

    rl := bufio.NewReader(f)
    //rs := csv.NewReader(f)
    // rs.Comma = '\t'      // Use tab-separated values

    for {
		//r, err := rs.Read()
        s, err := rl.ReadString(10) // 0x0A separator = newline
        if err == io.EOF {
               // fmt.Printf("ReadFile: EOF\n")
                return lines
        } else if err != nil {
                panic("reading file")
        }
		if s[len(s)-1] == '\n' {
			s = s[:len(s)-1]
		}
		if s[len(s)-1] == '\r' {
			s = s[:len(s)-1]
		}
		if s[len(s)-1] == ' ' {
			s = s[:len(s)-1]
		}
		//fmt.Printf("%q\n", s)
		if cb != nil {
			cb(s)
		}
        lines++
    }
}

func TestA(file string, hf2 string) (ht *HashTable) {
	//var lines int
/*
	var countlines = func(line string) {
		lines++
	}
*/
	var addLine = func(line string) {
		ht.Insert([]byte(line))
	}

	//fmt.Printf("\t%20q: ", hf2)
	//fmt.Printf("run: file=%q\n", file)
	lines := ReadFile(file, nil)
	//fmt.Printf("run: lines=%d, hf2=%q\n", lines, hf2)
	ht = NewHashTable(lines)
	//fmt.Printf("ht=%v\n", ht)
	ReadFile(file, addLine)
	return
}

func TestB(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line += "\n"
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestC(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line += line + "\n\n\n\n"
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestD(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line = "ABCDE" + line
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestE(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line = line + line
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestF(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line = line + line + line + line
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func reverse(s string) string {
	if len(s) == 0 {
		return ""
	}
	return reverse(s[1:]) + string(s[0])
}

func TestG(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line2 := reverse(line)
		//fmt.Printf("line=%q, line2=%q", line, line2)
		ht.Insert([]byte(line2))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestH(file string, hf2 string) (ht *HashTable) {
	var cnt int
	var counter = func(word string) {
		cnt++
	}
	var addWord = func(word string) {
		ht.Insert([]byte(word))
	}
	//test := []string{"abcdefgh", "efghijkl", "ijklmnop", "mnopqrst", "qrstuvwx", "uvwxyz01"} // 262144 words

	genWords(letters, counter)
	ht = NewHashTable(cnt)
	genWords(letters, addWord)
	return
}

func TestI(file string, hf2 string) (ht *HashTable) {
	//fmt.Printf("ni=%d\n", *ni)
	bs := make([]byte, 4, 4)
	ht = NewHashTable(*ni)
	for i := 0; i < *ni; i++ {
		bs[0], bs[1], bs[2], bs[3] = byte(i), byte(i>>8), byte(i>>16), byte(i>>24)
		ht.Insert(bs)
		//fmt.Printf("i=%d, 0x%08x, h=0x%08x\n", i, i, h)
	}
	return
}

func TestJ(file string, hf2 string) (ht *HashTable) {
	length := 900
	keys := length * 8
	key := make([]byte, length, length)
	key = key[:]
	ht = NewHashTable(keys)
	for k := range key {
		for i := uint(0); i < 8; i++ {
			key[k] = 1 << i
			//fmt.Printf("k=%d, i=%d, key=%v\n", k, i, key)
			ht.Insert(key)
			key[k] = 0
		}
	}
	return
}

// [ABCDEFGH][EFGHIJKL][IJKLMNOP][MNOPQRST][QRSTUVWX][UVWXYZ01]
// given a slice of strings, generate all the combinations in order
func genWords(perms []string, f func(word string)) {
	var indices = make([]int, len(perms), len(perms))
	var idx int
	var inc = func() bool {
		// increment counter with carry
		for idx = 0 ;; {
			indices[idx]++
			if indices[idx] >= len(perms[idx]) {
				indices[idx] = 0
				idx++
				if (idx >= len(perms)) {
					return true
				}
				continue
			} else {
				break
			}
		}
		return false
	}
	var letter = func(idx int, s string) string {
		return string(s[indices[idx]])
	}
	var word func(p []string) string
	word = func(p []string) string {
		if len(p) == 0 {
			return ""
		}
		l := len(p)
		idx := len(perms) - l
		tmp := letter(idx, p[0]) + word(p[1:])
		return tmp
	}
	// generate a word, hand it out, bump counter, repeat
	for {
		aword := word(perms)
		f(aword)
		if inc() {
			return
		}
	}
}

func tdiff(begin, end time.Time) time.Duration {
    d := end.Sub(begin)
    return d
}

func benchmark32s(n int) {
	//var hashes = make(Uint32Slice, n)
	//var u = make([]uint32, 1, 1)
	bs := make([]byte, 4, 4)
	var pn = hrff.Int64{int64(n), ""}
	var ps = hrff.Int64{int64(n*4), "B"}
	fmt.Printf("benchmark32s: gen n=%d, n=%h, size=%h\n", n, pn, ps)
	start := time.Now()
	for i := 0; i < n; i++ {
		bs[0], bs[1], bs[2], bs[3] = byte(i)&0xFF, (byte(i)>>8)&0xFF, (byte(i)>>16)&0xFF, (byte(i)>>24)&0xFF
		_ = jenkins.Hash232(bs, 0)
		//hashes[i] = h
		//fmt.Printf("i=%d, 0x%08x, h=0x%08x\n", i, i, h)
	}
	stop := time.Now()
	d := tdiff(start, stop)
	hsec := hrff.Float64{(float64(n) / d.Seconds()), "hashes/sec"}
	bsec := hrff.Float64{(float64(n) * 4 / d.Seconds()), "B/sec"}
	fmt.Printf("benchmark32s: %h\n", hsec)
	fmt.Printf("benchmark32s: %h\n", bsec)
	return

	fmt.Printf("benchmark32s: sort n=%d\n", n)
	//hashes.Sort()
/*
	for i := 0; i < n; i++ {
		fmt.Printf("i=%d, 0x%08x, h=0x%08x\n", i, i, hashes[i])
	}
*/
	fmt.Printf("benchmark32s: dup check n=%d\n", n)
	//dups, mrun := checkForDups32(hashes)
	//fmt.Printf("benchmark32: dups=%d, mrun=%d\n", dups, mrun)
}

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Uint32Slice) Sort() { sort.Sort(p) }

func checkForDups32(u Uint32Slice) (dups, mrun int) {
	i := 0
	run := 0
	for k, v := range u {
		if k == 0 || i == k {
			continue
		}
		if u[i] == v {
			run++
			dups++
			continue
		} else {
			if run > mrun {
				mrun = run
			}
			run = 0
			i = k
		}
	}
	return
}

func benchmark32g(h nhash.HashF32, n int) {
	var hashes = make(Uint32Slice, n)
	//var u = make([]uint32, 1, 1)
	bs := make([]byte, 4, 4)
	var pn = hrff.Int64{int64(n), ""}
	var ps = hrff.Int64{int64(n*4), "B"}
	fmt.Printf("benchmark32g: gen n=%d, n=%h, size=%h\n", n, pn, ps)
	start := time.Now()
	for i := 0; i < n; i++ {
		bs[0], bs[1], bs[2], bs[3] = byte(i), byte(i>>8), byte(i>>16), byte(i>>24)
		hashes[i] = h.Hash32(bs)
		//hashes[i] = h
		//fmt.Printf("i=%d, 0x%08x, h=0x%08x\n", i, i, h)
	}
	stop := time.Now()
	d := tdiff(start, stop)
	hsec := hrff.Float64{(float64(n) / d.Seconds()), "hashes/sec"}
	bsec := hrff.Float64{(float64(n) * 4 / d.Seconds()), "B/sec"}
	fmt.Printf("benchmark32g: %h\n", hsec)
	fmt.Printf("benchmark32g: %h\n", bsec)
	//return

	fmt.Printf("benchmark32g: sort n=%d\n", n)
	hashes.Sort()
/*
	for i := 0; i < n; i++ {
		fmt.Printf("i=%d, 0x%08x, h=0x%08x\n", i, i, hashes[i])
	}
*/
	fmt.Printf("benchmark32g: dup check n=%d\n", n)
	dups, mrun := checkForDups32(hashes)
	fmt.Printf("benchmark32: dups=%d, mrun=%d\n", dups, mrun)
}

func halloc(hfs string) (hf32 nhash.HashF32) {
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

var benchmarks = []string{"j332c", "j232", "sbox", "CrapWow"}
//var benchmarks = []string{"j332c"}

func benchmark(hashes []string, n int) {
	for _, v := range hashes {
		hf32 := halloc(v)
		fmt.Printf("benchmark32g: %q\n", v)
		benchmark32g(hf32, n)
		fmt.Printf("\n")
	}
}

type Test struct {
	name		string
	flag		**bool
	ptf			func(file string, hashf string) (ht *HashTable)
	desc		string
}

var Tests = []Test{
	{"TestA", &A, TestA, "insert keys"},
	{"TestB", &B, TestB, "add newline to key"},	
	{"TestC", &C, TestC, "add 4 newlines to key"},
	{"TestD", &D, TestD, "prepend ABCDE to key"},
	{"TestE", &E, TestE, "add 1 duplicate key"},		
	{"TestF", &F, TestF, "add 3 duplicate keys"},
	{"TestG", &G, TestF, "reverse letter order in key"},
	{"TestH", &H, TestH, "words from letter combinations in wc"},
	{"TestI", &I, TestI, "integers from 0 to ni-1 (does not read file)"},
	{"TestJ", &J, TestJ, "one bit keys (does not read file)"},
}

func runTestsWithFileAndHashes(file string, hf []string) {
	var test Test
	var print = func(s *HashTable) {
		q := s.HashQuality()
		if *oa {
			if test.name != "TestI" && test.name != "TestJ" && (s.Lines != s.Inserts || s.Lines != s.Heads || s.Lines != s.Nbuckets || s.Lines != s.Entries) {
				panic("runTestsWithFileAndHashes")
			}
			fmt.Printf("inserts=%d, size=%d, cols=%d, probes=%d, cpi=%0.2f%%, ppi=%04.2f, dups=%d\n",
				s.Inserts, s.Size, s.Cols, s.Probes, float64(s.Cols)/float64(s.Size)*100.0, float64(s.Probes)/float64(s.Inserts), s.Dups)
		} else {
			if test.name != "TestI" && test.name != "TestJ" && (s.Lines != s.Inserts || s.Lines != s.Probes || s.Lines != s.Entries) {
				fmt.Printf("lines=%d, inserts=%d, probes=%d, entries=%d\n", s.Lines, s.Inserts, s.Probes, s.Entries)
				panic("runTestsWithFileAndHashes")
			}
			fmt.Printf("inserts=%d, size=%d, buckets=%d, dups=%d, q=%0.2f\n",
				s.Inserts, s.Size, s.Nbuckets, s.Dups, q)
		}
	}
	if file != "" {
		fmt.Printf("file=%q\n", file)
	}
	for _, test = range Tests {
		if **test.flag {
			fmt.Printf("%s - %s\n", test.name, test.desc)
			for _, hf2 = range hf {
				hi := HashFunctions[hf2]
				if *c && !hi.Crypto {
					continue
				}
				if *h32 && hi.Size != 32 {
					continue
				}
				if *h64 && hi.Size != 64 {
					continue
				}
				fmt.Printf("\t%20q: ", hf2)
				ht := test.ptf(file, hf2)
				print(ht)
			}
		}
	}
	if *b {
		benchmark32s(*n)
		benchmark(benchmarks, *n)
	}
}

var file = flag.String("file", "", "words to read")
var hf = flag.String("hf", "all", "hash function")
var extra = flag.Int("e", 1, "extra bis in table size")
var prime = flag.Bool("p", false, "table size is primes and use mod")
var all = flag.Bool("a", false, "run all tests")
var pd = flag.Bool("pd", false, "print duplicate hashes")
var oa = flag.Bool("oa", false, "open addressing (no buckets)")

var c = flag.Bool("c", false, "only test crypto hash functions")
var h32 = flag.Bool("h32", false, "only test 32 bit has functions")
var h64 = flag.Bool("h64", false, "only test 64 bit has functions")

//var wc = flags.String("wc", "abcdefgh, efghijkl, ijklmnop, mnopqrst, qrstuvwx, uvwxyz01", "letter combinations for word") // 262144 words)
var ni = flag.Int("ni", 200000, "number of integer keys")
var n = flag.Int("n", 100000000, "number of hashes for benchmark")
var b = flag.Bool("b", false, "run benchmarks")
var A = flag.Bool("A", false, "test A")
var B = flag.Bool("B", false, "test B")
var C = flag.Bool("C", false, "test C")
var D = flag.Bool("D", false, "test D")
var E = flag.Bool("E", false, "test E")
var F = flag.Bool("F", false, "test F")
var G = flag.Bool("G", false, "test G")
var H = flag.Bool("H", false, "test H")
var I = flag.Bool("I", false, "test I")
var J = flag.Bool("J", false, "test J")

var letters = []string{"abcdefgh", "efghijkl", "ijklmnop", "mnopqrst", "qrstuvwx", "uvwxyz01"} // 262144 words
var TestPointers = []**bool{&A, &B, &C, &D, &E, &F, &G, &H, &I, &J}


func allTestsOn() {
	*A, *B, *C, *D, *E, *F, *G, *H , *I, *J = true, true, true, true, true, true, true, true, true, true
}

func allTestsOff() {
	*A, *B, *C, *D, *E, *F, *G, *H , *I, *J = false, false, false, false, false, false, false, false, false, false
}


func main() {
/*
	var cnt int
	var f = func(word string) {
		cnt++
		//fmt.Printf("%q\n", word)
	}

	test := []string{"ab", "cd"}
	test := []string{"abcdefgh", "efghijkl", "ijklmnop", "mnopqrst", "qrstuvwx", "uvwxyz01"} // 262144 words
	genWords(test, f)
	fmt.Printf("cnt=%d\n", cnt)
	return
*/
	flag.Parse()
	if *all {
		allTestsOn()
	}
	//fmt.Printf("%d lines read\n", lines)

	// read file and count lines
	// create table
	// read file and insert
	// stats

	k160 = keccak.New160()
	skein256 = skein.New256()
	//skein32 := skein.New(256, 32)
	sha1160 = sha1.New()
	a32 = adler32.New()
	switch {
	case *file != "":
		if *hf == "all" {
			runTestsWithFileAndHashes(*file, hashFunctions)
		} else {
			hf2 = *hf
			runTestsWithFileAndHashes(*file, []string{*hf})
		}
	case len(flag.Args()) != 0:
		for _, v := range flag.Args() {
			if *hf == "all" {
				runTestsWithFileAndHashes(v, hashFunctions)
			} else {
				hf2 = *hf
				runTestsWithFileAndHashes(v, []string{*hf})
			}
		}
	case len(flag.Args()) == 0:
		fmt.Printf("default\n")
		// no files specified run the only two tests we can with the specified hash functions
		allTestsOff()
		*I, *J = true, true
		if *hf == "all" {
			runTestsWithFileAndHashes("", hashFunctions)
		} else {
			hf2 = *hf
			runTestsWithFileAndHashes("", []string{*hf})
		}
	}
	if *b {
		benchmark32s(*n)
		fmt.Printf("\n")
		benchmark(benchmarks, *n)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s: [flags] [dictionary-files]\n", os.Args[0])
    	flag.PrintDefaults()
	}
}

/*
		case *test = A:
			fmt.Printf("TestA (simple hash check)\n")
			for _, hf2 = range hf {
				s = TestA(file, hf2)
				print(s)
			}
			*A = false
		case *B:
			fmt.Printf("TestB (add newline)\n")
			for _, hf2 = range hf {
				s = TestB(file, hf2)
				print(s)
			}
			*B = false
		case *C:
			fmt.Printf("TestC (add 4 newlines)\n")
			for _, hf2 = range hf {
				s = TestC(file, hf2)
				print(s)
			}
			*C = false
		case *D:
			fmt.Printf("TestD (prepend ABCDE)\n")
			for _, hf2 = range hf {
				s = TestD(file, hf2)
				print(s)
			}
			*D = false
		case *E:
			fmt.Printf("TestE (add 1 dup)\n")
			for _, hf2 = range hf {
				s = TestE(file, hf2)
				print(s)
			}
			*E = false
		case *F:
			fmt.Printf("TestF (add 3 dups)\n")
			for _, hf2 = range hf {
				s = TestF(file, hf2)
				print(s)
			}
			*F = false
		case *G:
			fmt.Printf("TestG (reverse word)\n")
			for _, hf2 = range hf {
				s = TestG(file, hf2)
				print(s)
			}
			*G = false
			fmt.Printf("\n")
		case *H:
			fmt.Printf("TestH (words from letter combinations)\n")
			for _, hf2 = range hf {
				s = TestH(file, hf2)
				print(s)
			}
			*H = false
		case *I:
			fmt.Printf("TestI (integers from 0 to n-1)\n")
			for _, hf2 = range hf {
				s = TestI(file, hf2)
				print(s)
			}
			*I = false
		case *J:
			fmt.Printf("TestI (one bit keys)\n")
			for _, hf2 = range hf {
				s = TestI(file, hf2)
				print(s)
			}
			*J = false
*/
