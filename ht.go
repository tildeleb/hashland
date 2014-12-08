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
	"crypto/sha1"
	"github.com/tildeleb/hashland/jenkins"
	"github.com/tildeleb/hashland/mahash"
	"github.com/tildeleb/hashland/spooky"
	"github.com/tildeleb/hashland/siphash"
	"github.com/tildeleb/hashland/keccak"
	"github.com/tildeleb/hashland/skein"
	//"github.com/tildeleb/hashland/threefish"
	"github.com/tildeleb/cuckoo/primes"
	_ "github.com/tildeleb/hrff"
)

var k160 hash.Hash
var skein256 hash.Hash
var sha1160 hash.Hash
var hashFunctions = []string{"MaHash8v64", "j332c", "j332b", "j232", "j264l", "j264h", "j264xor", "spooky32", "siphashal", "siphashah", "siphashbl", "siphashbh",
	"skein256xor", "skein256low", "skein256hi", "sha1160", "keccak160l",
	}
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
	Size uint32
	SizeLog2 uint32
	SizeMask uint32
}

type HashTable struct {
	Buckets [][]Bucket
	Stats
}

func hashf(k []byte) uint32 {
	var seeds []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	var fp = make([]byte, 32)
	switch hf2 {
	case "MaHash8v64":
		h64 := mahash.MaHash8v64(k)
		return uint32(h64)
	case "j332c":
		c, _ := jenkins.Jenkins364(k, len(k), 0, 0)
		return c
	case "j332b":
		_, b := jenkins.Jenkins364(k, len(k), 0, 0)
		return b
	case "j232":
		hash := jenkins.Hash232(k, 0)
		return hash
	case "j264l":
		hash := jenkins.Hash264(k, 0)
		return uint32(hash&0xFFFFFFFF)
	case "j264h":
		hash := jenkins.Hash264(k, 0)
		return uint32((hash>>32)&0xFFFFFFFF)
	case "j264xor":
		hash := jenkins.Hash264(k, 0)
		return uint32(hash&0xFFFFFFFF) ^ uint32((hash>>32)&0xFFFFFFFF)
	case "spooky32":
		return spooky.Hash32(k, 0)
	case "siphashal":
		a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
		return uint32(a&0xFFFFFFFF)
	case "siphashah":
		a, _ := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, false)
		return uint32((a>>32)&0xFFFFFFFF)
	case "siphashbl":
		_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
		return uint32(b&0xFFFFFFFF)
	case "siphashbh":
		_, b := siphash.Siphash(k, seeds, siphash.Crounds, siphash.Drounds, true)
		return uint32((b>>32)&0xFFFFFFFF)
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
	        return uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low)
		} else {
			return uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3])
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
	        return uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low)
		} else {
	    	return uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3])
	    }
	case "skein256low":
		skein256.Reset()
		skein256.Write(k)
		fp = fp[0:0]
		fp := skein256.Sum(fp)
		//fmt.Printf("skein256: fp=%v\n", fp)
    	return uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3])
	case "skein256hi":
		skein256.Reset()
		skein256.Write(k)
		fp = fp[0:0]
		fp := skein256.Sum(fp)
		//fmt.Printf("skein256: fp=%v\n", fp)
    	return uint32(fp[28])<<24 | uint32(fp[29])<<16 | uint32(fp[30])<<8 | uint32(fp[31])
	case "sha1160":
		sha1160.Reset()
		sha1160.Write(k)
		fp = fp[0:0]
		fp := sha1160.Sum(fp)
		if false {
	        low := fp[0] ^ fp[4] ^ fp[8] ^ fp[12] ^ fp[16]
	        med := fp[1] ^ fp[5] ^ fp[9] ^ fp[13] ^ fp[17]
	        hii := fp[2] ^ fp[6] ^ fp[10] ^ fp[14] ^ fp[18]
	        top := fp[3] ^ fp[7] ^ fp[11] ^ fp[15] ^ fp[19]
        	return uint32(top)<<24 | uint32(hii)<<16 | uint32(med)<<8 | uint32(low)
		} else {
			return uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[2])<<8 | uint32(fp[3])
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
	ht.SizeLog2 = NextLog2(uint32(ht.Lines)) + uint32(*extra)
	ht.Size = 1 << ht.SizeLog2
	if *prime {
		ht.Size = uint32(primes.NextPrime(int(ht.Size)))
	}
	ht.SizeMask = ht.Size - 1
	ht.Buckets = make([][]Bucket, ht.Size, ht.Size)
	return ht
}

func (ht *HashTable) Add(k []byte) {
	ht.Inserts++
	idx := uint32(0)
	hash := hashf(k) // jenkins.Hash232(k, 0)
	if *prime {
		idx = hash % ht.Size
	} else {
		idx = hash & ht.SizeMask
	}
	//fmt.Printf("index=%d\n", idx)
	cnt := 0
	pass := 0

	//fmt.Printf("Add: %s\n", k)
	//ht.Buckets[idx].Key = k
	//len(ht.Buckets[idx].Key) == 0
	for {
		if ht.Buckets[idx] == nil {
			// no entry or chain at this location, make it
			ht.Buckets[idx] = append(ht.Buckets[idx], Bucket{Key: k})
			//fmt.Printf("Add: len=%d, %s\n", len(ht.Buckets[idx]), ht.Buckets[idx][0].Key)
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
			h := hashf(ht.Buckets[idx][0].Key)
			if h == hash {
				if *pd {
					fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", hash, idx, k)
					fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", h, idx, ht.Buckets[idx][0].Key)
				}
				ht.Dups++
			}
			idx++
			cnt++
			if idx > uint32(ht.Size) - 1 {
				pass++
				if pass > 1 {
					panic("Add: pass")
				}
				idx = 0
			}
		} else {
			// first scan slice for dups
			for j := range ht.Buckets[idx] {
				h := hashf(ht.Buckets[idx][j].Key)
				if h == hash {
					if *pd {
						fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", hash, idx, k)
						fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", h, idx, ht.Buckets[idx][0].Key)
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
		ht.Add([]byte(line))
	}

	fmt.Printf("\t%20q: ", hf2)
	//fmt.Printf("run: file=%q\n", file)
	lines := ReadFile(file, nil)
	//fmt.Printf("run: lines=%d, hf2=%q\n", lines, hf2)
	ht = NewHashTable(lines)
	//fmt.Printf("ht=%v\n", ht)
	ReadFile(file, addLine)
	return
}

func TestB(file string, hf2 string) (ht *HashTable) {
	//var lines int
	var addLine = func(line string) {
		line += "\n"
		ht.Add([]byte(line))
	}
	fmt.Printf("\t%20q: ", hf2)
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestC(file string, hf2 string) (ht *HashTable) {
	//var lines int
	var addLine = func(line string) {
		line += line + "\n\n\n\n"
		ht.Add([]byte(line))
	}
	fmt.Printf("\t%20q: ", hf2)
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestD(file string, hf2 string) (ht *HashTable) {
	//var lines int
	var addLine = func(line string) {
		line = "ABCDE" + line
		ht.Add([]byte(line))
	}
	fmt.Printf("\t%20q: ", hf2)
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestE(file string, hf2 string) (ht *HashTable) {
	//var lines int
	var addLine = func(line string) {
		line = line + line
		ht.Add([]byte(line))
	}
	fmt.Printf("\t%20q: ", hf2)
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines)
	ReadFile(file, addLine)
	return
}

func TestF(file string, hf2 string) (ht *HashTable) {
	//var lines int
	var addLine = func(line string) {
		line = line + line + line + line
		ht.Add([]byte(line))
	}
	fmt.Printf("\t%20q: ", hf2)
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
	//var lines int
	var addLine = func(line string) {
		line2 := reverse(line)
		//fmt.Printf("line=%q, line2=%q", line, line2)
		ht.Add([]byte(line2))
	}
	fmt.Printf("\t%20q: ", hf2)
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
		ht.Add([]byte(word))
	}
	//test := []string{"abcdefgh", "efghijkl", "ijklmnop", "mnopqrst", "qrstuvwx", "uvwxyz01"} // 262144 words
	test := []string{"abcdefgh", "efghijkl", "ijklmnop", "mnopqrst", "qrstuvwx", "uvwxyz01"} // 262144 words
	genWords(test, counter)
	fmt.Printf("\t%20q: ", hf2)
	ht = NewHashTable(cnt)
	genWords(test, addWord)
	return
}

func TestI(file string, hf2 string) (ht *HashTable) {
	length := 900
	keys := length * 8
	key := make([]byte, length, length)
	key = key[:]
	fmt.Printf("\t%20q: ", hf2)
	ht = NewHashTable(keys)
	for k := range key {
		for i := uint(0); i < 8; i++ {
			key[k] = 1 << i
			//fmt.Printf("k=%d, i=%d, key=%v\n", k, i, key)
			ht.Add(key)
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

func runTestsWithFileAndHashes(file string, hf []string) {
	var s *HashTable
	var print = func(s *HashTable) {
		q := s.HashQuality()
		fmt.Printf("lines=%d, inserts=%d, size=%d, cols=%d, probes=%d, heads=%d, buckets=%d, entries=%d, cpi=%0.2f%%, ppi=%04.2f, dups=%d, q=%0.2f\n",
			s.Lines, s.Inserts, s.Size, s.Cols, s.Probes, s.Heads, s.Nbuckets, s.Entries, float64(s.Cols)/float64(s.Size)*100.0, float64(s.Probes)/float64(s.Inserts), s.Dups, q)
	}
	fmt.Printf("file=%q\n", file)
	for {
		switch {
		case *A:
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
			fmt.Printf("TestI (one bit keys)\n")
			for _, hf2 = range hf {
				s = TestI(file, hf2)
				print(s)
			}
			*I = false
		default:
			return
		}
		fmt.Printf("\n")
	}
}

var file = flag.String("file", "", "words to read")
var hf = flag.String("hf", "all", "hash function")
var extra = flag.Int("e", 1, "extra bis in table size")
var prime = flag.Bool("p", false, "table size is primes and use mod")
var all = flag.Bool("a", false, "run all tests")
var pd = flag.Bool("pd", false, "print duplicate hashes")
var oa = flag.Bool("oa", false, "open addressing (no buckets)")
var A = flag.Bool("A", false, "test A")
var B = flag.Bool("B", false, "test B")
var C = flag.Bool("C", false, "test C")
var D = flag.Bool("D", false, "test D")
var E = flag.Bool("E", false, "test E")
var F = flag.Bool("F", false, "test F")
var G = flag.Bool("G", false, "test G")
var H = flag.Bool("H", false, "test H")
var I = flag.Bool("I", false, "test I")

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
		*A, *B, *C, *D, *E, *F, *G, *H , *I = true, true, true, true, true, true, true, true, true
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
	if *file != "" {
		if *hf == "all" {
			runTestsWithFileAndHashes(*file, hashFunctions)
		} else {
			hf2 = *hf
			runTestsWithFileAndHashes(*file, []string{*hf})
		}
	}
	for _, v := range flag.Args() {
		if *hf == "all" {
			runTestsWithFileAndHashes(v, hashFunctions)
		} else {
			hf2 = *hf
			runTestsWithFileAndHashes(v, []string{*hf})
		}
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s: [flags] [dictionary-files]\n", os.Args[0])
    	flag.PrintDefaults()
	}
}
