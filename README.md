HashLand
========

**Trust no code or hash functions here yet**

**Work in progress, more hashes to come soon**

**Some licensing information missing; will be rectified soon**


Introduction
------------
HashLand contains the following functionality.

1. A (currently barely) curated set of pure Go hash functions including various jenkins hash functions and his latest SpookyV2, Murmur3, Cityhash, sbox, MaHash8v64, CrapWow, Siphash, keccak, skein and more.

2. AES based hash functions extracted Go's runtime and used by the map implementation on Intel X86 architecture machines that support AES-NI. (not working yet)

3. Tests which (mostly) use file based dictionaries to gather statistics about the above hash functions.

4. An extraction with a little generalization of the SMHasher functions from the Go runtime. These functions were already ported from the SMHasher C code by the Go Authors.

5. The ability to benchmark hash functions.

6. A package, "nhash" with a proposed new set of Go interfaces for hash functions that complement the existing core Go streaming hash interface. The core of the proposal is:

	`Hash32(b []byte, seeds ...uint32) uint32`  
	`Hash64(b []byte, seeds ...uint64) uint64`  
	`Hash128(b []byte, seeds ...uint64) (uint64, uint64)`  
	`Hash(in []byte, out []byte, seeds ...uint64) []byte`  
	I am considering changing some or all of the variadic arguments to `…byte` but I suspect doing so will make calling these functions less convenient. I would love feedback on that.

Background
----------
In the process of writing a [cuckoo hash table](https://github.com/tildeleb/cuckoo) I wondered which hash functions would be the best ones to use. I was frustrated that the core Go libraries didn't contain a useful set of hash functions for constructing hash tables. Also, I had always wanted to experiment with hash functions. So I spend some time doing building HashLand to figure it all out. I wrote a bunch of hash functions in pure Go. I spent some time with the gc inliner seeing how much optimization could be done in the context of Go. I forked some hash functions from other repositories, and tested them out. I ended up putting in a bit more effort than I had planned

Quality of Hash Functions
-------------------------
There are no "bad" hash functions here. Most of these do a very good job of hashing keys with good distribution and few duplicate hashes. However, my tests and dictionaries are still very basic. I've been focused on getting the hash functions written and some testing infrastructure up and running. blah blah blah.

Performance
-----------
Some of these are woefully lacking in performance. Many will be difficult to improve with the gc based compilers in pure Go. If you want speed and need crypto quality use ...

Warning
-------
*Don't use the non crypto hash functions if you have uncontrolled inputs (i.e. a web based API or web facing data inputs or an adversary. If you do, use at least SipHash or one of the other crypto hash functions.*

Roadmap
----------
1. A few more hash functions
2. Make sure licensing and author information is accurate
3. Performance optimization
4. A few more tests
5. Better dictionaries
6. Better stats

Non Crypto Hash Functions
-------------------------
	"sbox":			simple hash function         
	"CrapWow":		another simple hash function
	"MaHash8v64":	russian hash function
	"j332c":		Jenkins lookup3 c bits hash
	"j332b":		Jenkins lookup3 b bits hash
	"j232":			Jenkins lookup8 32 bit
	"j264l": 		Jenkins lookup8 64 bit (low bits)
	"j264h": 		Jenkins lookup8 64 bit (high bits)
	"j264xor":		Jenkins lookup8 64 bit (high xor low bits)
	"spooky32":	Jenkins Spooky, 32 bit
	"spooky64":	Jenkins Spooky, 64 bit
	"spooky128h":	Jenkins Spooky, 128 bit, high half
	"spooky128l:	Jenkins Spooky, 128 bit, low half
	"spooky128xor:	Jenkins Spooky, 128 bit, low xor half

Crypto Hash Functions
---------------------
	"siphashal": 
	"siphashah": 
	"siphashbl": 
	"siphashbh": 
	"skein256xor": 
	"skein256low": 
	"skein256hi": 
	"sha1160": 
	"keccak160l"

Usage
-----
	Usage of ./hashland:
	./hashland: [flags] [dictionary-files]
	  -A=false: test A
	  -B=false: test B
	  -C=false: test C
	  -D=false: test D
	  -E=false: test E
	  -F=false: test F
	  -G=false: test G
	  -H=false: test H
	  -I=false: test I
	  -J=false: test J
	  -a=false: run all tests
	  -b=false: run benchmarks
	  -c=false: only test crypto hash functions
	  -cd=false: check for duplicate hashs when running benchmarks
	  -e=1: extra bis in table size
	  -file="": words to read
	  -h32=false: only test 32 bit has functions
	  -h64=false: only test 64 bit has functions
	  -hf="all": hash function
	  -n=100000000: number of hashes for benchmark
	  -ni=200000: number of integer keys
	  -oa=false: open addressing (no buckets)
	  -p=false: table size is primes and use mod
	  -pd=false: print duplicate hashes
	  -sm=false: run SMHasher
	  -v=false: verbose

SMHasher
--------
Still some work to do on this, particularly with `-v`.  

	leb@hula: % hashland -sm -hf=j264 -v
	"TestSmhasherSanity": 118.968289ms
	"TestSmhasherSeed": 30.632449ms
	"TestSmhasherText": 5.923103295s
	"TestSmhasherWindowed": 58.097608577s
	"TestSmhasherAvalanche": 		z=100000, n=16
		z=100000, n=32
		z=100000, n=64
		z=100000, n=128
		z=100000, n=256
		z=100000, n=1600
		z=100000, n=32
		z=100000, n=64
	1m6.748766357s
	"TestSmhasherPermutation": 
		n=8, s=[0 1 2 3 4 5 6 7]
		n=8, s=[0 536870912 1073741824 1610612736 2147483648 2684354560 3221225472 3758096384]
		n=20, s=[0 1]
		n=20, s=[0 2147483648]
		n=6, s=[0 1 2 3 4 5 6 7 536870912 1073741824 1610612736 2147483648 2684354560 3221225472 3758096384]20.717874625s
	"TestSmhasherSparse": 9.937053904s
	"TestSmhasherCyclic": 4.438443435s
	"TestSmhasherSmallKeys": 16.713253851s
	"TestSmhasherZeros": 6.013460606s
	"TestSmhasherAppendedZeros": 126.096µs

Benchmarks (currently broken)
-----------------------------
	leb@hula:~/gotest/src/github.com/tildeleb/hashland % hashland -b -hf=j364 -v
	
	ksiz=4, len(bs)=4
	benchmark32g: gen n=100000000, n=100 M, keySize=4,  size=400 MB
	benchmark32g: 38 Mhashes/sec
	benchmark32g: 151 MB/sec
	
	ksiz=8, len(bs)=8
	benchmark32g: gen n=100000000, n=100 M, keySize=8,  size=800 MB
	benchmark32g: 36 Mhashes/sec
	benchmark32g: 288 MB/sec
	
	ksiz=16, len(bs)=16
	benchmark32g: gen n=100000000, n=100 M, keySize=16,  size=1 GB
	benchmark32g: 22 Mhashes/sec
	benchmark32g: 345 MB/sec
	
	ksiz=32, len(bs)=32
	benchmark32g: gen n=100000000, n=100 M, keySize=32,  size=3 GB
	benchmark32g: 15 Mhashes/sec
	benchmark32g: 469 MB/sec
	
	ksiz=64, len(bs)=64
	benchmark32g: gen n=100000000, n=100 M, keySize=64,  size=6 GB
	benchmark32g: 8 Mhashes/sec
	benchmark32g: 508 MB/sec
	
	ksiz=512, len(bs)=512
	benchmark32g: gen n=10000000, n=10 M, keySize=512,  size=5 GB
	benchmark32g: 1 Mhashes/sec
	benchmark32g: 577 MB/sec
	
	ksiz=1024, len(bs)=1024
	benchmark32g: gen n=10000000, n=10 M, keySize=1024,  size=10 GB
	benchmark32g: 570 khashes/sec
	benchmark32g: 584 MB/sec
	
	leb@hula:~/gotest/src/github.com/tildeleb/hashland %

Tests
-----

	leb@hula:hashland % hashland -A -hf="all" -oa db/words-vak.txt 
	file="db/words-vak.txt"
	TestA - insert keys
		              "j364": inserts=326796, size=1048576, cols=51054, probes=349873, cpi=4.87%, ppi=1.07, dups=0
		              "j264": inserts=326796, size=1048576, cols=50784, probes=349673, cpi=4.84%, ppi=1.07, dups=0
		       "siphash128a": inserts=326796, size=1048576, cols=50847, probes=349580, cpi=4.85%, ppi=1.07, dups=0
		       "siphash128b": inserts=326796, size=1048576, cols=50826, probes=350076, cpi=4.85%, ppi=1.07, dups=0
		        "MaHash8v64": inserts=326796, size=1048576, cols=51022, probes=349666, cpi=4.87%, ppi=1.07, dups=0
		          "spooky64": inserts=326796, size=1048576, cols=50673, probes=349547, cpi=4.83%, ppi=1.07, dups=0
		        "spooky128h": inserts=326796, size=1048576, cols=50673, probes=349547, cpi=4.83%, ppi=1.07, dups=0
		        "spooky128l": inserts=326796, size=1048576, cols=50776, probes=349864, cpi=4.84%, ppi=1.07, dups=0
		      "spooky128xor": inserts=326796, size=1048576, cols=50906, probes=349734, cpi=4.85%, ppi=1.07, dups=0
		              "sbox": inserts=326796, size=1048576, cols=50755, probes=349435, cpi=4.84%, ppi=1.07, dups=26
		             "j332c": inserts=326796, size=1048576, cols=51054, probes=349873, cpi=4.87%, ppi=1.07, dups=9
		             "j332b": inserts=326796, size=1048576, cols=50564, probes=349852, cpi=4.82%, ppi=1.07, dups=10
		              "j232": inserts=326796, size=1048576, cols=50577, probes=349523, cpi=4.82%, ppi=1.07, dups=20
		             "j264l": inserts=326796, size=1048576, cols=50784, probes=349673, cpi=4.84%, ppi=1.07, dups=16
		             "j264h": inserts=326796, size=1048576, cols=50997, probes=350064, cpi=4.86%, ppi=1.07, dups=19
		           "j264xor": inserts=326796, size=1048576, cols=50437, probes=349331, cpi=4.81%, ppi=1.07, dups=13
		          "spooky32": inserts=326796, size=1048576, cols=50673, probes=349547, cpi=4.83%, ppi=1.07, dups=12
		       "siphash64al": inserts=326796, size=1048576, cols=50794, probes=349621, cpi=4.84%, ppi=1.07, dups=13
		       "siphash64ah": inserts=326796, size=1048576, cols=51062, probes=349332, cpi=4.87%, ppi=1.07, dups=7
		       "siphash64bl": inserts=326796, size=1048576, cols=50826, probes=350076, cpi=4.85%, ppi=1.07, dups=21
		       "siphash64bh": inserts=326796, size=1048576, cols=50763, probes=349610, cpi=4.84%, ppi=1.07, dups=12
		       "skein256xor": inserts=326796, size=1048576, cols=51001, probes=349511, cpi=4.86%, ppi=1.07, dups=9
		       "skein256low": inserts=326796, size=1048576, cols=50826, probes=349510, cpi=4.85%, ppi=1.07, dups=15
		        "skein256hi": inserts=326796, size=1048576, cols=50734, probes=350099, cpi=4.84%, ppi=1.07, dups=13
		              "sha1": inserts=326796, size=1048576, cols=50976, probes=349704, cpi=4.86%, ppi=1.07, dups=14
		        "keccak160l": inserts=326796, size=1048576, cols=50998, probes=349799, cpi=4.86%, ppi=1.07, dups=13
	leb@hula:hashland % hashland -A -hf="all" db/words-vak.txt 
	file="db/words-vak.txt"
	TestA - insert keys
		              "j364": inserts=326796, size=1048576, buckets=280641, dups=0, q=1.00
		              "j264": inserts=326796, size=1048576, buckets=280818, dups=0, q=1.00
		       "siphash128a": inserts=326796, size=1048576, buckets=280774, dups=0, q=1.00
		       "siphash128b": inserts=326796, size=1048576, buckets=280759, dups=0, q=1.00
		        "MaHash8v64": inserts=326796, size=1048576, buckets=280705, dups=0, q=1.00
		          "spooky64": inserts=326796, size=1048576, buckets=280983, dups=0, q=1.00
		        "spooky128h": inserts=326796, size=1048576, buckets=280983, dups=0, q=1.00
		        "spooky128l": inserts=326796, size=1048576, buckets=280886, dups=0, q=1.00
		      "spooky128xor": inserts=326796, size=1048576, buckets=280787, dups=0, q=1.00
		              "sbox": inserts=326796, size=1048576, buckets=280956, dups=26, q=1.00
		             "j332c": inserts=326796, size=1048576, buckets=280641, dups=9, q=1.00
		             "j332b": inserts=326796, size=1048576, buckets=281107, dups=10, q=1.00
		              "j232": inserts=326796, size=1048576, buckets=281053, dups=20, q=1.00
		             "j264l": inserts=326796, size=1048576, buckets=280818, dups=16, q=1.00
		             "j264h": inserts=326796, size=1048576, buckets=280703, dups=19, q=1.00
		           "j264xor": inserts=326796, size=1048576, buckets=281126, dups=13, q=1.00
		          "spooky32": inserts=326796, size=1048576, buckets=280983, dups=12, q=1.00
		       "siphash64al": inserts=326796, size=1048576, buckets=280859, dups=13, q=1.00
		       "siphash64ah": inserts=326796, size=1048576, buckets=280691, dups=7, q=1.00
		       "siphash64bl": inserts=326796, size=1048576, buckets=280759, dups=21, q=1.00
		       "siphash64bh": inserts=326796, size=1048576, buckets=280887, dups=12, q=1.00
		       "skein256xor": inserts=326796, size=1048576, buckets=280616, dups=9, q=1.00
		       "skein256low": inserts=326796, size=1048576, buckets=280803, dups=15, q=1.00
		        "skein256hi": inserts=326796, size=1048576, buckets=280949, dups=13, q=1.00
		              "sha1": inserts=326796, size=1048576, buckets=280710, dups=14, q=1.00
		        "keccak160l": inserts=326796, size=1048576, buckets=280733, dups=13, q=1.00
	leb@hula:~/gotest/src/github.com/tildeleb/hashland % 




