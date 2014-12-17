package main

import (
	"fmt"
	_ "unsafe"
	"github.com/tildeleb/hashland/aeshash"
)

var strs = []string{"abcd", "efgh", "blow", "deadbeef"}
var ui32 = []uint32{1, 2, 3, 4}
var ui64= []uint64{1, 2, 3, 4, 5, 6, 7, 8}
var hash uint64

func main() {
	for _, v := range strs {
		for seed := uint32(0); seed < 2; seed++ {
			b := make([]byte, 8)
			b = b[:]
			//copy(b, v)
			s := uint64(seed)
			//h := uint64(0)
			//pb := &i8[k]
		    //up := unsafe.Pointer(v)
			//fmt.Printf("s=%d, h=%d\n", s, h)
			copy(b, v)
			b = b[0:len(v)]
			//hash := aeshash.HashStr(v, s)
			hash := aeshash.Hash(b, s)
			fmt.Printf("key=%v, seed=%d, hash=0x%016x\n", v, seed, hash)
		}
	}
}