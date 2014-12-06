hashes
======

Go versions of murmur, cityhash, siphash, jenkins, spooky and more.
Code to evaluate the quality of each hash function.

** Work in progress, more hashes to come soon **

** Some licensing information missing; will be rectified soon **

	leb@hula:~/gotest/src/github.com/tildeleb/hashland % ./hashland db/*
	file="db/dic_Shakespeare.txt"
		        "MaHash8v64": lines=3228, size=8192, cols=643, per=7.85%, probes=450, dups=0
		             "j332c": lines=3228, size=8192, cols=616, per=7.52%, probes=455, dups=0
		             "j332b": lines=3228, size=8192, cols=672, per=8.20%, probes=473, dups=0
		              "j232": lines=3228, size=8192, cols=648, per=7.91%, probes=480, dups=0
		             "j264l": lines=3228, size=8192, cols=650, per=7.93%, probes=485, dups=0
		             "j264h": lines=3228, size=8192, cols=634, per=7.74%, probes=405, dups=0
		           "j264xor": lines=3228, size=8192, cols=648, per=7.91%, probes=403, dups=0
		          "spooky32": lines=3228, size=8192, cols=625, per=7.63%, probes=375, dups=0
		         "siphashal": lines=3228, size=8192, cols=634, per=7.74%, probes=391, dups=0
		         "siphashah": lines=3228, size=8192, cols=648, per=7.91%, probes=393, dups=0
		         "siphashbl": lines=3228, size=8192, cols=666, per=8.13%, probes=447, dups=0
		         "siphashbh": lines=3228, size=8192, cols=668, per=8.15%, probes=455, dups=0
		       "skein256xor": lines=3228, size=8192, cols=640, per=7.81%, probes=391, dups=0
		           "sha1160": lines=3228, size=8192, cols=611, per=7.46%, probes=380, dups=0
		        "keccak160l": lines=3228, size=8192, cols=580, per=7.08%, probes=369, dups=0
	
	file="db/dic_common_words.txt"
		        "MaHash8v64": lines=500, size=1024, cols=121, per=11.82%, probes=114, dups=0
		             "j332c": lines=500, size=1024, cols=110, per=10.74%, probes=103, dups=0
		             "j332b": lines=500, size=1024, cols=115, per=11.23%, probes=110, dups=0
		              "j232": lines=500, size=1024, cols=115, per=11.23%, probes=137, dups=0
		             "j264l": lines=500, size=1024, cols=116, per=11.33%, probes=91, dups=0
		             "j264h": lines=500, size=1024, cols=122, per=11.91%, probes=126, dups=0
		           "j264xor": lines=500, size=1024, cols=109, per=10.64%, probes=124, dups=0
		          "spooky32": lines=500, size=1024, cols=117, per=11.43%, probes=94, dups=0
		         "siphashal": lines=500, size=1024, cols=115, per=11.23%, probes=78, dups=0
		         "siphashah": lines=500, size=1024, cols=106, per=10.35%, probes=96, dups=0
		         "siphashbl": lines=500, size=1024, cols=113, per=11.04%, probes=86, dups=0
		         "siphashbh": lines=500, size=1024, cols=126, per=12.30%, probes=93, dups=0
		       "skein256xor": lines=500, size=1024, cols=126, per=12.30%, probes=115, dups=0
		           "sha1160": lines=500, size=1024, cols=129, per=12.60%, probes=123, dups=0
		        "keccak160l": lines=500, size=1024, cols=108, per=10.55%, probes=83, dups=0
	
	file="db/dic_fr.txt"
		        "MaHash8v64": lines=13408, size=32768, cols=2790, per=8.51%, probes=1866, dups=0
		             "j332c": lines=13408, size=32768, cols=2733, per=8.34%, probes=1971, dups=0
		             "j332b": lines=13408, size=32768, cols=2749, per=8.39%, probes=1903, dups=0
		              "j232": lines=13408, size=32768, cols=2784, per=8.50%, probes=2085, dups=0
		             "j264l": lines=13408, size=32768, cols=2719, per=8.30%, probes=1712, dups=0
		             "j264h": lines=13408, size=32768, cols=2721, per=8.30%, probes=1822, dups=0
		           "j264xor": lines=13408, size=32768, cols=2700, per=8.24%, probes=1947, dups=0
		          "spooky32": lines=13408, size=32768, cols=2778, per=8.48%, probes=2000, dups=0
		         "siphashal": lines=13408, size=32768, cols=2759, per=8.42%, probes=1917, dups=0
		         "siphashah": lines=13408, size=32768, cols=2690, per=8.21%, probes=1773, dups=0
		         "siphashbl": lines=13408, size=32768, cols=2778, per=8.48%, probes=1871, dups=0
		         "siphashbh": lines=13408, size=32768, cols=2739, per=8.36%, probes=1944, dups=0
		       "skein256xor": lines=13408, size=32768, cols=2705, per=8.26%, probes=1863, dups=0
		           "sha1160": lines=13408, size=32768, cols=2770, per=8.45%, probes=2048, dups=0
		        "keccak160l": lines=13408, size=32768, cols=2742, per=8.37%, probes=1873, dups=0
	
	file="db/dic_ip.txt"
		        "MaHash8v64": lines=3925, size=8192, cols=932, per=11.38%, probes=911, dups=5
		             "j332c": lines=3925, size=8192, cols=931, per=11.36%, probes=929, dups=9
		             "j332b": lines=3925, size=8192, cols=986, per=12.04%, probes=1027, dups=7
		              "j232": lines=3925, size=8192, cols=923, per=11.27%, probes=712, dups=6
		             "j264l": lines=3925, size=8192, cols=939, per=11.46%, probes=849, dups=7
		             "j264h": lines=3925, size=8192, cols=944, per=11.52%, probes=807, dups=8
		           "j264xor": lines=3925, size=8192, cols=998, per=12.18%, probes=907, dups=5
		          "spooky32": lines=3925, size=8192, cols=923, per=11.27%, probes=876, dups=5
		         "siphashal": lines=3925, size=8192, cols=941, per=11.49%, probes=960, dups=7
		         "siphashah": lines=3925, size=8192, cols=900, per=10.99%, probes=767, dups=6
		         "siphashbl": lines=3925, size=8192, cols=946, per=11.55%, probes=890, dups=8
		         "siphashbh": lines=3925, size=8192, cols=904, per=11.04%, probes=926, dups=9
		       "skein256xor": lines=3925, size=8192, cols=912, per=11.13%, probes=754, dups=8
		           "sha1160": lines=3925, size=8192, cols=946, per=11.55%, probes=752, dups=8
		        "keccak160l": lines=3925, size=8192, cols=894, per=10.91%, probes=839, dups=9
	
	file="db/dic_numbers.txt"
		        "MaHash8v64": lines=500, size=1024, cols=123, per=12.01%, probes=156, dups=0
		             "j332c": lines=500, size=1024, cols=108, per=10.55%, probes=97, dups=0
		             "j332b": lines=500, size=1024, cols=113, per=11.04%, probes=129, dups=0
		              "j232": lines=500, size=1024, cols=117, per=11.43%, probes=88, dups=0
		             "j264l": lines=500, size=1024, cols=130, per=12.70%, probes=148, dups=0
		             "j264h": lines=500, size=1024, cols=135, per=13.18%, probes=135, dups=0
		           "j264xor": lines=500, size=1024, cols=132, per=12.89%, probes=104, dups=0
		          "spooky32": lines=500, size=1024, cols=148, per=14.45%, probes=131, dups=0
		         "siphashal": lines=500, size=1024, cols=115, per=11.23%, probes=88, dups=0
		         "siphashah": lines=500, size=1024, cols=119, per=11.62%, probes=96, dups=0
		         "siphashbl": lines=500, size=1024, cols=120, per=11.72%, probes=136, dups=0
		         "siphashbh": lines=500, size=1024, cols=117, per=11.43%, probes=99, dups=0
		       "skein256xor": lines=500, size=1024, cols=129, per=12.60%, probes=118, dups=0
		           "sha1160": lines=500, size=1024, cols=120, per=11.72%, probes=110, dups=0
		        "keccak160l": lines=500, size=1024, cols=132, per=12.89%, probes=115, dups=0
	
	file="db/dic_postfix.txt"
		        "MaHash8v64": lines=500, size=1024, cols=120, per=11.72%, probes=106, dups=0
		             "j332c": lines=500, size=1024, cols=110, per=10.74%, probes=121, dups=0
		             "j332b": lines=500, size=1024, cols=115, per=11.23%, probes=99, dups=0
		              "j232": lines=500, size=1024, cols=128, per=12.50%, probes=152, dups=0
		             "j264l": lines=500, size=1024, cols=119, per=11.62%, probes=73, dups=0
		             "j264h": lines=500, size=1024, cols=112, per=10.94%, probes=148, dups=0
		           "j264xor": lines=500, size=1024, cols=124, per=12.11%, probes=124, dups=0
		          "spooky32": lines=500, size=1024, cols=122, per=11.91%, probes=115, dups=0
		         "siphashal": lines=500, size=1024, cols=121, per=11.82%, probes=137, dups=0
		         "siphashah": lines=500, size=1024, cols=119, per=11.62%, probes=122, dups=0
		         "siphashbl": lines=500, size=1024, cols=119, per=11.62%, probes=97, dups=0
		         "siphashbh": lines=500, size=1024, cols=131, per=12.79%, probes=129, dups=0
		       "skein256xor": lines=500, size=1024, cols=137, per=13.38%, probes=163, dups=0
		           "sha1160": lines=500, size=1024, cols=117, per=11.43%, probes=111, dups=0
		        "keccak160l": lines=500, size=1024, cols=122, per=11.91%, probes=123, dups=0
	
	file="db/dic_prefix.txt"
		        "MaHash8v64": lines=500, size=1024, cols=114, per=11.13%, probes=117, dups=0
		             "j332c": lines=500, size=1024, cols=115, per=11.23%, probes=78, dups=0
		             "j332b": lines=500, size=1024, cols=121, per=11.82%, probes=133, dups=0
		              "j232": lines=500, size=1024, cols=127, per=12.40%, probes=115, dups=0
		             "j264l": lines=500, size=1024, cols=114, per=11.13%, probes=123, dups=0
		             "j264h": lines=500, size=1024, cols=128, per=12.50%, probes=169, dups=0
		           "j264xor": lines=500, size=1024, cols=113, per=11.04%, probes=123, dups=0
		          "spooky32": lines=500, size=1024, cols=123, per=12.01%, probes=99, dups=0
		         "siphashal": lines=500, size=1024, cols=115, per=11.23%, probes=111, dups=0
		         "siphashah": lines=500, size=1024, cols=129, per=12.60%, probes=131, dups=0
		         "siphashbl": lines=500, size=1024, cols=113, per=11.04%, probes=57, dups=0
		         "siphashbh": lines=500, size=1024, cols=120, per=11.72%, probes=83, dups=0
		       "skein256xor": lines=500, size=1024, cols=115, per=11.23%, probes=125, dups=0
		           "sha1160": lines=500, size=1024, cols=107, per=10.45%, probes=103, dups=0
		        "keccak160l": lines=500, size=1024, cols=135, per=13.18%, probes=144, dups=0
	
	file="db/dic_variables.txt"
		        "MaHash8v64": lines=1842, size=4096, cols=425, per=10.38%, probes=292, dups=0
		             "j332c": lines=1842, size=4096, cols=421, per=10.28%, probes=324, dups=0
		             "j332b": lines=1842, size=4096, cols=415, per=10.13%, probes=295, dups=0
		              "j232": lines=1842, size=4096, cols=403, per=9.84%, probes=381, dups=0
		             "j264l": lines=1842, size=4096, cols=402, per=9.81%, probes=290, dups=0
		             "j264h": lines=1842, size=4096, cols=412, per=10.06%, probes=306, dups=0
		           "j264xor": lines=1842, size=4096, cols=394, per=9.62%, probes=313, dups=0
		          "spooky32": lines=1842, size=4096, cols=419, per=10.23%, probes=367, dups=0
		         "siphashal": lines=1842, size=4096, cols=428, per=10.45%, probes=301, dups=0
		         "siphashah": lines=1842, size=4096, cols=405, per=9.89%, probes=368, dups=0
		         "siphashbl": lines=1842, size=4096, cols=398, per=9.72%, probes=387, dups=0
		         "siphashbh": lines=1842, size=4096, cols=414, per=10.11%, probes=339, dups=0
		       "skein256xor": lines=1842, size=4096, cols=439, per=10.72%, probes=505, dups=0
		           "sha1160": lines=1842, size=4096, cols=407, per=9.94%, probes=342, dups=0
		        "keccak160l": lines=1842, size=4096, cols=419, per=10.23%, probes=290, dups=0
	
	file="db/dic_win32.txt"
		        "MaHash8v64": lines=1991, size=4096, cols=476, per=11.62%, probes=385, dups=0
		             "j332c": lines=1991, size=4096, cols=471, per=11.50%, probes=517, dups=0
		             "j332b": lines=1991, size=4096, cols=468, per=11.43%, probes=358, dups=0
		              "j232": lines=1991, size=4096, cols=511, per=12.48%, probes=513, dups=0
		             "j264l": lines=1991, size=4096, cols=460, per=11.23%, probes=435, dups=0
		             "j264h": lines=1991, size=4096, cols=485, per=11.84%, probes=401, dups=0
		           "j264xor": lines=1991, size=4096, cols=467, per=11.40%, probes=398, dups=0
		          "spooky32": lines=1991, size=4096, cols=472, per=11.52%, probes=404, dups=0
		         "siphashal": lines=1991, size=4096, cols=469, per=11.45%, probes=381, dups=0
		         "siphashah": lines=1991, size=4096, cols=492, per=12.01%, probes=519, dups=0
		         "siphashbl": lines=1991, size=4096, cols=461, per=11.25%, probes=437, dups=0
		         "siphashbh": lines=1991, size=4096, cols=495, per=12.08%, probes=473, dups=0
		       "skein256xor": lines=1991, size=4096, cols=460, per=11.23%, probes=419, dups=0
		           "sha1160": lines=1991, size=4096, cols=490, per=11.96%, probes=505, dups=0
		        "keccak160l": lines=1991, size=4096, cols=508, per=12.40%, probes=449, dups=0
	
	file="db/words-vak-alt.txt"
		        "MaHash8v64": lines=310594, size=1048576, cols=45979, per=4.38%, probes=19358, dups=12
		             "j332c": lines=310594, size=1048576, cols=46048, per=4.39%, probes=19335, dups=8
		             "j332b": lines=310594, size=1048576, cols=45909, per=4.38%, probes=19739, dups=4
		              "j232": lines=310594, size=1048576, cols=45904, per=4.38%, probes=18868, dups=12
		             "j264l": lines=310594, size=1048576, cols=46290, per=4.41%, probes=19276, dups=10
		             "j264h": lines=310594, size=1048576, cols=45781, per=4.37%, probes=19204, dups=6
		           "j264xor": lines=310594, size=1048576, cols=46216, per=4.41%, probes=19100, dups=9
		          "spooky32": lines=310594, size=1048576, cols=46029, per=4.39%, probes=19507, dups=12
		         "siphashal": lines=310594, size=1048576, cols=46195, per=4.41%, probes=19524, dups=10
		         "siphashah": lines=310594, size=1048576, cols=46201, per=4.41%, probes=19188, dups=11
		         "siphashbl": lines=310594, size=1048576, cols=45865, per=4.37%, probes=19254, dups=9
		         "siphashbh": lines=310594, size=1048576, cols=46067, per=4.39%, probes=19800, dups=9
		       "skein256xor": lines=310594, size=1048576, cols=45800, per=4.37%, probes=19301, dups=11
		           "sha1160": lines=310594, size=1048576, cols=46000, per=4.39%, probes=19472, dups=9
		        "keccak160l": lines=310594, size=1048576, cols=45880, per=4.38%, probes=19190, dups=8
	
	file="db/words-vak.txt"
		        "MaHash8v64": lines=326796, size=1048576, cols=51022, per=4.87%, probes=22870, dups=10
		             "j332c": lines=326796, size=1048576, cols=51054, per=4.87%, probes=23077, dups=7
		             "j332b": lines=326796, size=1048576, cols=50564, per=4.82%, probes=23056, dups=9
		              "j232": lines=326796, size=1048576, cols=50577, per=4.82%, probes=22727, dups=17
		             "j264l": lines=326796, size=1048576, cols=50784, per=4.84%, probes=22877, dups=15
		             "j264h": lines=326796, size=1048576, cols=50997, per=4.86%, probes=23268, dups=17
		           "j264xor": lines=326796, size=1048576, cols=50437, per=4.81%, probes=22535, dups=12
		          "spooky32": lines=326796, size=1048576, cols=50673, per=4.83%, probes=22751, dups=10
		         "siphashal": lines=326796, size=1048576, cols=50883, per=4.85%, probes=23335, dups=9
		         "siphashah": lines=326796, size=1048576, cols=50999, per=4.86%, probes=23337, dups=12
		         "siphashbl": lines=326796, size=1048576, cols=50792, per=4.84%, probes=23174, dups=9
		         "siphashbh": lines=326796, size=1048576, cols=51056, per=4.87%, probes=23243, dups=11
		       "skein256xor": lines=326796, size=1048576, cols=51001, per=4.86%, probes=22715, dups=8
		           "sha1160": lines=326796, size=1048576, cols=50976, per=4.86%, probes=22908, dups=13
		        "keccak160l": lines=326796, size=1048576, cols=50998, per=4.86%, probes=23003, dups=12
	
	leb@hula:~/gotest/src/github.com/tildeleb/hashland % 



