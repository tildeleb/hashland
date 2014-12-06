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
	_ "github.com/tildeleb/hrff"
)

var k160 hash.Hash
var skein256 hash.Hash
var sha1160 hash.Hash

type Bucket struct {
	Key []byte
}

type Stats struct {
	Cols int
	Probes int
	Dups int
}

type HashTable struct {
	Buckets []Bucket
	Lines int
	Size uint32
	SizeLog2 uint32
	SizeMask uint32
	Stats
}

var hashFunctions = []string{"MaHash8v64", "j332c", "j332b", "j232", "j264l", "j264h", "j264xor", "spooky32", "siphashal", "siphashah", "siphashbl", "siphashbh",
	"skein256xor", "sha1160", "keccak160l",
	}
var hf2 string
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
	    	return uint32(fp[0])<<24 | uint32(fp[1])<<16 | uint32(fp[1])<<8 | uint32(fp[2])
	    }
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
			panic("hashf")
	}
	return 0
}

func NewHashTable(lines int) *HashTable {
	ht := new(HashTable)
	ht.Lines = lines
	ht.SizeLog2 = NextLog2(uint32(ht.Lines)) + uint32(*extra)
	ht.Size = 1 << ht.SizeLog2
	ht.SizeMask = ht.Size - 1
	ht.Buckets = make([]Bucket, ht.Size, ht.Size)
	return ht
}

func (ht *HashTable) Add(k []byte) {
	hash := hashf(k) // jenkins.Hash232(k, 0)
	idx := hash&ht.SizeMask
	//fmt.Printf("index=%d\n", idx)
	cnt := 0
	for {
		if len(ht.Buckets[idx].Key) == 0 {
			//fmt.Printf("Add: %s\n", k)
			//ht.Buckets[idx].Key = make([]byte, len(k), len(k))
			//ht.Buckets[idx].Key = ht.Buckets[idx].Key[:]
			//copy(ht.Buckets[idx].Key, k)
			ht.Buckets[idx].Key = k
			//fmt.Printf("Add: %s\n", ht.Buckets[idx].Key)
			return
		}
		if cnt == 0 {
			ht.Cols++
			h := hashf(ht.Buckets[idx].Key)
			if h == hash {
				//fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", hash, idx, k)
				//fmt.Printf("hash=0x%08x, idx=%d, key=%q\n", h, idx, ht.Buckets[idx].Key)
				ht.Dups++
			}
		} else {
			ht.Probes++
		}
		idx++
		cnt++
		if idx > uint32(ht.Size) - 1 {
			idx = 0
		}
	}
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

func run(file string, hf2 string) {
	//var lines int
	var ht *HashTable
/*
	var countlines = func(line string) {
		lines++
	}
*/
	var addLine = func(line string) {
		ht.Add([]byte(line))
	}

	//fmt.Printf("run: file=%q\n", file)
	lines := ReadFile(file, nil)
	//fmt.Printf("run: lines=%d, hf2=%q\n", lines, hf2)
	ht = NewHashTable(lines)
	//fmt.Printf("ht=%v\n", ht)
	ReadFile(file, addLine)
	fmt.Printf("lines=%d, size=%d, cols=%d, per=%0.2f%%, probes=%d, dups=%d\n", ht.Lines, ht.Size, ht.Cols, float64(ht.Cols)/float64(ht.Size)*100.0, ht.Probes, ht.Dups) 
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

var file = flag.String("file", "", "words to read")
var hf = flag.String("hf", "j332c", "hash function")
var extra = flag.Int("e", 1, "extra bis in table size")
func main() {
	flag.Parse()
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
		run(*file, *hf)
	}
	for _, v := range flag.Args() {
		fmt.Printf("file=%q\n", v)
		for _, hf2 = range hashFunctions {
			fmt.Printf("\t%20q: ", hf2)
			run(v, hf2)
		}
		fmt.Printf("\n")
	}
}
