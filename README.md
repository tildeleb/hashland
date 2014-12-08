HashLand
========

**Trust no code or hash functions here yet**

**Work in progress, more hashes to come soon**

**Some licensing information missing; will be rectified soon**

HashLand provides pure Go versions of murmur, cityhash, jenkins, spooky, sbox, crapwow, siphash, keccak, skein and more.

The "hash land" application evaluates the quality of each hash function with real world testing.

Introduction
------------
As part of another project I decided to write some hash functions in Go and test them out. There are no "bad" hash functions here. Must of these do a very good job of hashing keys with good distribution. However, my tests and dictionaries are still very basic. I've been focused on getting the hash functions written.

Don't use the non crypto hash functions if you have uncontrolled input or an adversary. If you do use at least SipHash.

Next Steps
----------
1. A few more hash functions
2. Make sure licensing and author information is accurate
3. Performance optimization
4. Better tests
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
	"spooky32":		Jenkins Spooky 32 bit flavor

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
	Usage of hashland:
	hashland: [flags] [dictionary-files]
	  -A=false: test A
	  -B=false: test B
	  -C=false: test C
	  -D=false: test D
	  -E=false: test E
	  -F=false: test F
	  -G=false: test G
	  -H=false: test H
	  -I=false: test I
	  -a=false: run all tests
	  -e=1: extra bis in table size
	  -file="": words to read
	  -hf="all": hash function
	  -oa=false: open addressing (no buckets)
	  -p=false: table size is primes and use mod
	  -pd=false: print duplicate hashes

Benchmarks
----------

	leb@hula:~/gotest/src/github.com/tildeleb/hashland % ./hashland -A -hf="all" -oa db/dic_common_words.txt                    
	file="db/dic_common_words.txt"
	TestA (simple hash check)
		              "sbox": lines=500, size=1024, cols=119, probes=586, cpi=11.62%, ppi=1.17, dups=0
		           "CrapWow": lines=500, size=1024, cols=111, probes=585, cpi=10.84%, ppi=1.17, dups=0
		        "MaHash8v64": lines=500, size=1024, cols=121, probes=614, cpi=11.82%, ppi=1.23, dups=0
		             "j332c": lines=500, size=1024, cols=110, probes=603, cpi=10.74%, ppi=1.21, dups=0
		             "j332b": lines=500, size=1024, cols=115, probes=610, cpi=11.23%, ppi=1.22, dups=0
		              "j232": lines=500, size=1024, cols=115, probes=637, cpi=11.23%, ppi=1.27, dups=0
		             "j264l": lines=500, size=1024, cols=116, probes=591, cpi=11.33%, ppi=1.18, dups=0
		             "j264h": lines=500, size=1024, cols=122, probes=626, cpi=11.91%, ppi=1.25, dups=0
		           "j264xor": lines=500, size=1024, cols=109, probes=624, cpi=10.64%, ppi=1.25, dups=0
		          "spooky32": lines=500, size=1024, cols=117, probes=594, cpi=11.43%, ppi=1.19, dups=0
		         "siphashal": lines=500, size=1024, cols=115, probes=578, cpi=11.23%, ppi=1.16, dups=0
		         "siphashah": lines=500, size=1024, cols=106, probes=596, cpi=10.35%, ppi=1.19, dups=0
		         "siphashbl": lines=500, size=1024, cols=113, probes=586, cpi=11.04%, ppi=1.17, dups=0
		         "siphashbh": lines=500, size=1024, cols=126, probes=593, cpi=12.30%, ppi=1.19, dups=0
		       "skein256xor": lines=500, size=1024, cols=126, probes=615, cpi=12.30%, ppi=1.23, dups=0
		       "skein256low": lines=500, size=1024, cols=124, probes=610, cpi=12.11%, ppi=1.22, dups=0
		        "skein256hi": lines=500, size=1024, cols=115, probes=588, cpi=11.23%, ppi=1.18, dups=0
		           "sha1160": lines=500, size=1024, cols=129, probes=623, cpi=12.60%, ppi=1.25, dups=0
		        "keccak160l": lines=500, size=1024, cols=108, probes=583, cpi=10.55%, ppi=1.17, dups=0
	
	leb@hula:~/gotest/src/github.com/tildeleb/hashland % ./hashland -A -hf="all" db/dic_common_words.txt 
	file="db/dic_common_words.txt"
	TestA (simple hash check)
		              "sbox": lines=500, size=1024, buckets=397, dups=0, q=1.01
		           "CrapWow": lines=500, size=1024, buckets=403, dups=0, q=0.97
		        "MaHash8v64": lines=500, size=1024, buckets=398, dups=0, q=1.00
		             "j332c": lines=500, size=1024, buckets=402, dups=0, q=0.98
		             "j332b": lines=500, size=1024, buckets=398, dups=0, q=0.99
		              "j232": lines=500, size=1024, buckets=402, dups=0, q=0.99
		             "j264l": lines=500, size=1024, buckets=400, dups=0, q=1.00
		             "j264h": lines=500, size=1024, buckets=395, dups=0, q=1.00
		           "j264xor": lines=500, size=1024, buckets=402, dups=0, q=0.99
		          "spooky32": lines=500, size=1024, buckets=394, dups=0, q=1.00
		         "siphashal": lines=500, size=1024, buckets=398, dups=0, q=0.99
		         "siphashah": lines=500, size=1024, buckets=408, dups=0, q=0.97
		         "siphashbl": lines=500, size=1024, buckets=406, dups=0, q=0.97
		         "siphashbh": lines=500, size=1024, buckets=393, dups=0, q=1.00
		       "skein256xor": lines=500, size=1024, buckets=396, dups=0, q=1.01
		       "skein256low": lines=500, size=1024, buckets=394, dups=0, q=1.01
		        "skein256hi": lines=500, size=1024, buckets=404, dups=0, q=0.97
		           "sha1160": lines=500, size=1024, buckets=387, dups=0, q=1.01
		        "keccak160l": lines=500, size=1024, buckets=408, dups=0, q=0.97
	
	leb@hula:~/gotest/src/github.com/tildeleb/hashland % 



