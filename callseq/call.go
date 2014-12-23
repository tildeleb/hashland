package main 

import (
	"fmt"
	"time"
	"unsafe"
	"github.com/tildeleb/hashland/nhash"
	"github.com/tildeleb/hashland/jenkins"
	)

func nullhash(b []byte, seed uint64) uint64 {
	return 0
}

type digest struct {
}

func (d *digest) Size() int {
	return 8
}

func  (d *digest) BlockSize() int {
	return 8
}

func  (d *digest) NumSeedBytes() int {
	return 8
}

func  (d *digest) Hash64(b []byte, seeds ...uint64) uint64 {
	return 0
}

func  (d *digest) Hash64S(b []byte, seed uint64) uint64 {
	return 0
}

func New() nhash.HashF64 {
	var f digest
	return &f
}


func tdiff(begin, end time.Time) time.Duration {
    d := end.Sub(begin)
    return d
}

// initial results are even worse than I thought
// nullhash: 1.135146477s
// Hash64 with seed: 59.690928864s
// Hash64 no seed: 3.684813561s

const n = 1000000000
var b = make([]byte, 100, 100)
var seed uint64
var intf = New()

func b1() time.Duration {
	start := time.Now()
	for i:= 0; i < n; i++ {
		nullhash(b, seed)
	}
	stop := time.Now()
	return tdiff(start, stop)
}

func b2() time.Duration {
	start := time.Now()
	for i:= 0; i < n; i++ {
		intf.Hash64(b, seed)
	}
	stop := time.Now()
	return tdiff(start, stop)
}

func b3() time.Duration {
	start := time.Now()
	for i := 0; i < n; i++ {
		intf.Hash64(b)
	}
	stop := time.Now()
	return tdiff(start, stop)
}

type hf32 func(b []byte, seed uint32) uint32
type hf322 func(b []byte, l int, seeda, seedb uint32) (uint32, uint32)
type hf64 func(b []byte, seed uint64) uint64
type hf128e func(b []byte, seeda, seedb uint64) (uint64, uint64)
type ff func() time.Duration

type DispEntry struct {
	fp		unsafe.Pointer
	kind	int

}

var f2 = jenkins.Hash264
var f3 = jenkins.Jenkins364

var dispTable = []DispEntry{DispEntry{unsafe.Pointer(&f2), 2}, DispEntry{unsafe.Pointer(&f3), 3}}


// 		c, b := jenkins.Jenkins364(k, len(k), uint32(seed), uint32(seed))

func hash(de *DispEntry, b []byte, seed uint64) (ret uint64) {
	if de.kind == 2 {
		pf := (*hf64)(de.fp)
		ret = (*pf)(b, seed)
	} else if de.kind == 3 {
		pf := (*hf322)(de.fp)
		c, b  := (*pf)(b, len(b), uint32(seed), uint32(seed>>32))
		ret = uint64(b)<<32 | uint64(c)
	} else {
		panic("hash")
	}
	return
}


func main() {
	//var u uintptr

	key := make([]byte, 8, 8)
	i := uint64(0)
	key[0], key[1], key[2], key[3], key[4], key[5], key[6], key[7] = byte(i), byte(i>>8), byte(i>>16), byte(i>>24), byte(i>>32), byte(i>>40), byte(i>>48), byte(i>>56)

	l := 0 // len(key)
	seedc := uint32(0)
	seedb := uint32(0)
	key = key[0:0]

	hash1 := hash(&dispTable[0], key, 0)
	hash2 := hash(&dispTable[1], key, 0)
	fmt.Printf("hash1=%#x\n", hash1)
	fmt.Printf("hash2=%#x\n", hash2)

	addr := (ff)(b3)
	ib3 := b3
	up := unsafe.Pointer(&ib3)
	code := **(**uintptr)(up)

	ifp := jenkins.Jenkins364
	up2 := unsafe.Pointer(&ifp)

	pf := (*hf322)(up2)


	c, b := jenkins.Jenkins364(key, l, seedc, seedb)
	fmt.Printf("hash: c=%#X, b=%#X\n", c, b)
	c, b = (*pf)(key, l, seedc, seedb)
	fmt.Printf("hash: c=%#X, b=%#X\n", c, b)
	s := ""
	k := ([]byte)(s)
	fmt.Printf("len(k)=%d\n", len(k))
	c, b = jenkins.Jenkins364(k, len(k), seedc, seedb)
	fmt.Printf("hash: c=%#X, b=%#X\n", c, b)

	fmt.Printf("b3:\t%T %v %#X %p\n", b3, b3, b3, b3)
	fmt.Printf("ib3:\t%T %v %#X %p\n", addr, addr, addr, addr)
	fmt.Printf("ib3:\t%T %v %#X %p\n", ib3, ib3, ib3, ib3)
	fmt.Printf("up:\t%T %v %#X %p\n", up, up, up, up)
	fmt.Printf("code:\t%T %v %#X %p\n", code, code, code, code)
	return
	fmt.Printf("nullhash: ")
	t1 := b1()
	fmt.Printf("%v\n", t1)
	fmt.Printf("Hash64 with seed: ")
	t2 := b2()
	fmt.Printf("%v\n", t2)
	fmt.Printf("Hash64 no seed: ")
	t3 := b3()
	fmt.Printf("%v\n", t3)
}

