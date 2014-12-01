// Copyright © 2014 Lawrence E. Bakst. All rights reserved.
package jenkins3_test

//import "flag"
import "fmt"
//import "math"
//import "math/rand"
//import "runtime"
import "github.com/tildeleb/hashes/jenkins3"
import "testing"

func TestBasic(t *testing.T) {

	q := "This is the time for all good men to come to the aid of their country..."
	//qq := []byte{"xThis is the time for all good men to come to the aid of their country..."}
	//qqq := []byte{"xxThis is the time for all good men to come to the aid of their country..."}
	//qqqq[] := []byte{"xxxThis is the time for all good men to come to the aid of their country..."}

	u := stu(q)
	h1 := jenkins3.HashWordsLen(u, 13)
	fmt.Printf("%08x, %0x8, %08x\n", h1)

	b, c := uint32(0), uint32(0)
	c, b = jenkins3.HashString("", c, b)
	//fmt.Printf("%08x, %08x\n", c, b)
	if c != 0xdeadbeef || b != 0xdeadbeef {
		t.Logf("c=0x%x != 0xdeadbeef || b=0x%x != 0xdeadbeef\n", c, b)
		t.FailNow()
	}

	b, c = 0xdeadbeef, 0
	c, b = jenkins3.HashString("", c, b)
	//fmt.Printf("%08x, %08x\n", c, b)	// bd5b7dde deadbeef
	if c != 0xbd5b7dde || b != 0xdeadbeef {
		t.Logf("c=0x%x != 0xbd5b7dde || b=0x%x != 0xdeadbeef\n", c, b)
		t.FailNow()
	}

  	b, c = 0xdeadbeef, 0xdeadbeef
	c, b = jenkins3.HashString("", c, b)
	//fmt.Printf("%08x, %08x\n", c, b)	// 9c093ccd bd5b7dde
	if c != 0x9c093ccd || b != 0xbd5b7dde {
		t.Logf("c=0x%x != 0x9c093ccd || b=0x%x != 0xbd5b7dde\n", c, b)
		t.FailNow()
	}

	b, c = 0, 0
	c, b = jenkins3.HashString("Four score and seven years ago", c, b)
	//fmt.Printf("%08x, %08x\n", c, b)	// 17770551 ce7226e6
	if c != 0x17770551 || b != 0xce7226e6 {
		t.Logf("c=0x%x != 0x17770551 || b=0x%x != 0xce7226e6\n", c, b)
		t.FailNow()
	}

	b, c = 1, 0
	c, b = jenkins3.HashString("Four score and seven years ago", c, b)
	//fmt.Printf("%08x, %08x\n", c, b)	// e3607cae bd371de4
	if c != 0xe3607cae || b != 0xbd371de4 {
		t.Logf("c=0x%x != 0xe3607cae || b=0x%x != 0xbd371de4\n", c, b)
		t.FailNow()
	}

	b, c = 0, 1
	c, b = jenkins3.HashString("Four score and seven years ago", c, b)
	//fmt.Printf("%08x, %08x\n", c, b)	// cd628161 6cbea4b3
	if c != 0xcd628161 || b != 0x6cbea4b3 {
		t.Logf("c=0x%x != 0xcd628161 || b=0x%x != 0x6cbea4b3\n", c, b)
		t.FailNow()
	}

}

func BenchmarkJenkins(b *testing.B) {
	//tmp := make([]byte, 4, 4)
	us := make([]uint32, 1)
	b.SetBytes(int64(b.N * 4))
	for i := 1; i <= b.N; i++ {
		us[0] = uint32(i)
		//tmp[0], tmp[1], tmp[2], tmp[3] = byte(key&0xFF), byte((key>>8)&0xFF), byte((key>>16)&0xFF), byte((key>>24)&0xFF)
		jenkins3.HashWords(us, 0)
	}
}

/*
func main() {
	q := "This is the time for all good men to come to the aid of their country..."
	//qq := []byte{"xThis is the time for all good men to come to the aid of their country..."}
	//qqq := []byte{"xxThis is the time for all good men to come to the aid of their country..."}
	//qqqq[] := []byte{"xxxThis is the time for all good men to come to the aid of their country..."}

	u := stu(q)
	h1 := hashword(u, (len(q)-1)/4, 13)
	h2 := hashword(u, (len(q)-5)/4, 13)
	h3 := hashword(u, (len(q)-9)/4, 13)
	fmt.Printf("%08x, %0x8, %08x\n", h1, h2, h3)


}
*/