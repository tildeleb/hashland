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
	"github.com/tildeleb/hrff"
	. "github.com/tildeleb/hashland/hashf" // cleaved
	. "github.com/tildeleb/hashland/hashtable" // cleaved
	"github.com/tildeleb/hashland/nhash"
	"github.com/tildeleb/hashland/jenkins" // remove
	"sort"
	"time"
)

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
	ht = NewHashTable(lines, *extra, *pd, *oa, *prime)
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
	ht = NewHashTable(lines, *extra, *pd, *oa, *prime)
	ReadFile(file, addLine)
	return
}

func TestC(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line += line + "\n\n\n\n"
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines, *extra, *pd, *oa, *prime)
	ReadFile(file, addLine)
	return
}

func TestD(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line = "ABCDE" + line
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines, *extra, *pd, *oa, *prime)
	ReadFile(file, addLine)
	return
}

func TestE(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line = line + line
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines, *extra, *pd, *oa, *prime)
	ReadFile(file, addLine)
	return
}

func TestF(file string, hf2 string) (ht *HashTable) {
	var addLine = func(line string) {
		line = line + line + line + line
		ht.Insert([]byte(line))
	}
	lines := ReadFile(file, nil)
	ht = NewHashTable(lines, *extra, *pd, *oa, *prime)
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
	ht = NewHashTable(lines, *extra, *pd, *oa, *prime)
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
	ht = NewHashTable(cnt, *extra, *pd, *oa, *prime)
	genWords(letters, addWord)
	return
}

func TestI(file string, hf2 string) (ht *HashTable) {
	//fmt.Printf("ni=%d\n", *ni)
	bs := make([]byte, 4, 4)
	ht = NewHashTable(*ni, *extra, *pd, *oa, *prime)
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
	ht = NewHashTable(keys, *extra, *pd, *oa, *prime)
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

var keySizes = []int{4, 8, 16, 32, 64, 512, 1024}
func benchmark32s(n int) {
	//var hashes = make(Uint32Slice, n)
	const nbytes = 1024
	bs := make([]byte, nbytes, nbytes)
	bs = bs[:]
	for _, ksiz := range keySizes {
		if ksiz == 512 {
			n = n / 10
		}
		bs = bs[:ksiz]
		fmt.Printf("ksiz=%d, len(bs)=%d\n", ksiz, len(bs))
		pn := hrff.Int64{int64(n), ""}
		ps := hrff.Int64{int64(n*ksiz), "B"}
		fmt.Printf("benchmark32s: gen n=%d, n=%h, keySize=%d,  size=%h\n", n, pn, ksiz, ps)
		start := time.Now()
		for i := 0; i < n; i++ {
			bs[0], bs[1], bs[2], bs[3] = byte(i)&0xFF, (byte(i)>>8)&0xFF, (byte(i)>>16)&0xFF, (byte(i)>>24)&0xFF
			_ = jenkins.Hash232(bs, 0)
			//_, _ = jenkins.Jenkins364(bs, 0, 0, 0)
			//hashes[i] = h
			//fmt.Printf("i=%d, 0x%08x, h=0x%08x\n", i, i, h)
		}
		stop := time.Now()
		d := tdiff(start, stop)
		hsec := hrff.Float64{(float64(n) / d.Seconds()), "hashes/sec"}
		bsec := hrff.Float64{(float64(n) * float64(ksiz) / d.Seconds()), "B/sec"}
		fmt.Printf("benchmark32s: %h\n", hsec)
		fmt.Printf("benchmark32s: %h\n\n", bsec)
	}
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
	var hashes Uint32Slice
	const nbytes = 1024

	bs := make([]byte, nbytes, nbytes)
	bs = bs[:]
	for _, ksiz := range keySizes {
		if ksiz == 512 {
			n = n / 10
		}
		bs = bs[:ksiz]
		hashes = make(Uint32Slice, n, n)
		hashes = hashes[:]
		fmt.Printf("ksiz=%d, len(bs)=%d\n", ksiz, len(bs))
		pn := hrff.Int64{int64(n), ""}
		ps := hrff.Int64{int64(n*ksiz), "B"}
		fmt.Printf("benchmark32g: gen n=%d, n=%h, keySize=%d,  size=%h\n", n, pn, ksiz, ps)
		start := time.Now()
		for i := 0; i < n; i++ {
			bs[0], bs[1], bs[2], bs[3] = byte(i)&0xFF, (byte(i)>>8)&0xFF, (byte(i)>>16)&0xFF, (byte(i)>>24)&0xFF
			//_ = jenkins.Hash232(bs, 0)
			_, _ = jenkins.Jenkins364(bs, 0, 0, 0)
			//hashes[i] = h
			//fmt.Printf("i=%d, 0x%08x, h=0x%08x\n", i, i, h)
		}
		stop := time.Now()
		d := tdiff(start, stop)
		hsec := hrff.Float64{(float64(n) / d.Seconds()), "hashes/sec"}
		bsec := hrff.Float64{(float64(n) * float64(ksiz) / d.Seconds()), "B/sec"}
		fmt.Printf("benchmark32g: %h\n", hsec)
		fmt.Printf("benchmark32g: %h\n\n", bsec)
	}

	if *cd {
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
}

var benchmarks = []string{"j332c", "j232", "sbox", "CrapWow"}
//var benchmarks = []string{"j332c"}

func benchmark(hashes []string, n int) {
	for _, v := range hashes {
		hf32 := Halloc(v)
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
			for _, Hf2 = range hf {
				hi := HashFunctions[Hf2]
				if *c && !hi.Crypto {
					continue
				}
				if *h32 && hi.Size != 32 {
					continue
				}
				if *h64 && hi.Size != 64 {
					continue
				}
				fmt.Printf("\t%20q: ", Hf2)
				ht := test.ptf(file, Hf2)
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

var b = flag.Bool("b", false, "run benchmarks")
var cd = flag.Bool("cd", false, "check for duplicate hashs when running benchmarks")

//var wc = flags.String("wc", "abcdefgh, efghijkl, ijklmnop, mnopqrst, qrstuvwx, uvwxyz01", "letter combinations for word") // 262144 words)
var ni = flag.Int("ni", 200000, "number of integer keys")
var n = flag.Int("n", 100000000, "number of hashes for benchmark")
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

	switch {
	case *file != "":
		if *hf == "all" {
			runTestsWithFileAndHashes(*file, TestHashFunctions)
		} else {
			Hf2 = *hf
			runTestsWithFileAndHashes(*file, []string{*hf})
		}
	case len(flag.Args()) != 0:
		for _, v := range flag.Args() {
			if *hf == "all" {
				runTestsWithFileAndHashes(v, TestHashFunctions)
			} else {
				Hf2 = *hf
				runTestsWithFileAndHashes(v, []string{*hf})
			}
		}
	case len(flag.Args()) == 0 && !*b:
		// no files specified run the only two tests we can with the specified hash functions
		allTestsOff()
		*I, *J = true, false
		if *hf == "all" {
			runTestsWithFileAndHashes("", TestHashFunctions)
		} else {
			Hf2 = *hf
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
