package main 

import (
	"fmt"
	"time"
	"github.com/tildeleb/hashland/nhash"
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

func main() {
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

