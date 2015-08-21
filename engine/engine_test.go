package engine

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type TestCase struct {
	name         string
	initialStack [][]byte
	commands     string
	result       []byte
}

var testCases = []TestCase{
	{name: "no commands", initialStack: [][]byte{[]byte("Hello")}, commands: "/", result: []byte("Hello")},

	// hashing
	{name: "md5", initialStack: [][]byte{[]byte("Hello")}, commands: "/md5/hex", result: []byte("8b1a9953c4611296a827abf8c47804d7")},
	{name: "sha1", initialStack: [][]byte{[]byte("Hello")}, commands: "/sha1/hex", result: []byte("f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0")},
	{name: "sha224", initialStack: [][]byte{[]byte("Hello")}, commands: "/sha224/hex", result: []byte("4149da18aa8bfc2b1e382c6c26556d01a92c261b6436dad5e3be3fcc")},
	{name: "sha256", initialStack: [][]byte{[]byte("Hello")}, commands: "/sha256/hex", result: []byte("185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969")},
	{name: "sha384", initialStack: [][]byte{[]byte("Hello")}, commands: "/sha384/hex", result: []byte("3519fe5ad2c596efe3e276a6f351b8fc0b03db861782490d45f7598ebd0ab5fd5520ed102f38c4a5ec834e98668035fc")},
	{name: "sha512", initialStack: [][]byte{[]byte("Hello")}, commands: "/sha512/hex", result: []byte("3615f80c9d293ed7402687f94b22d58e529b8cc7916f8fac7fddf7fbd5af4cf777d3d795a7a00a16bf7e7f3fb9561ee9baae480da9fe7a18769e71886b03f315")},
	{name: "ripemd160", initialStack: [][]byte{[]byte("Hello")}, commands: "/ripemd160/hex", result: []byte("d44426aca8ae0a69cdbc4021c64fa5ad68ca32fe")},
	{name: "hmac-md5", initialStack: [][]byte{[]byte("TheData"), []byte("TheKey")}, commands: "/md5/hmac-md5/hex", result: []byte("05dd5de8c3fe0ec39161f287c81b2ff9")},
	{name: "hmac-sha1", initialStack: [][]byte{[]byte("TheData"), []byte("TheKey")}, commands: "/sha1/hmac-sha1/hex", result: []byte("4175329c2ece3d097adfec866022a02aa8ccf2d8")},
	{name: "hmac-sha224", initialStack: [][]byte{[]byte("TheData"), []byte("TheKey")}, commands: "/sha224/hmac-sha224/hex", result: []byte("fa38d389dfd66966b0408e61b366d330f52eff296604a1b3c3a863db")},
	{name: "hmac-sha256", initialStack: [][]byte{[]byte("TheData"), []byte("TheKey")}, commands: "/sha256/hmac-sha256/hex", result: []byte("21ad8c7172c3ead1627075d305785587d18b641758ed07ebe5b85c6095f778cf")},
	{name: "hmac-sha384", initialStack: [][]byte{[]byte("TheData"), []byte("TheKey")}, commands: "/sha384/hmac-sha384/hex", result: []byte("83ebc1d46f3570a820eae1893e1be5f07949dd3b53ce8f11c4b793fc6dc5131b6f06f791b9a3af9b4f1c5b44d2c2794f")},
	{name: "hmac-sha512", initialStack: [][]byte{[]byte("TheData"), []byte("TheKey")}, commands: "/sha512/hmac-sha512/hex", result: []byte("5d63666a22c5150a8c1a7aacf95a5f5419e5c8efd1755519b6ec0e35ece20df5c50e948a997ece52f8a004d82fff7faa90ee991c58e425b43dc221476d0a4122")},
	{name: "hmac-ripemd160", initialStack: [][]byte{[]byte("TheData"), []byte("TheKey")}, commands: "/ripemd160/hmac-ripemd160/hex", result: []byte("a18d7ea2da1ffdd5ef0a877979db4ab05e5096b0")},
	{name: "rand", initialStack: [][]byte{[]byte("45")}, commands: "/rand/len/swap/pop", result: []byte("45")},
	{name: "md5 len", initialStack: [][]byte{[]byte("Hello")}, commands: "/push/md5/len/swap/pop/md5-len/eq", result: []byte("Hello")},
	{name: "sha1 len", initialStack: [][]byte{[]byte("Hello")}, commands: "/push/sha1/len/swap/pop/sha1-len/eq", result: []byte("Hello")},
	{name: "sha224 len", initialStack: [][]byte{[]byte("Hello")}, commands: "/push/sha224/len/swap/pop/sha224-len/eq", result: []byte("Hello")},
	{name: "sha256 len", initialStack: [][]byte{[]byte("Hello")}, commands: "/push/sha256/len/swap/pop/sha256-len/eq", result: []byte("Hello")},
	{name: "sha384 len", initialStack: [][]byte{[]byte("Hello")}, commands: "/push/sha384/len/swap/pop/sha384-len/eq", result: []byte("Hello")},
	{name: "sha512 len", initialStack: [][]byte{[]byte("Hello")}, commands: "/push/sha512/len/swap/pop/sha512-len/eq", result: []byte("Hello")},
	{name: "ripemd160 len", initialStack: [][]byte{[]byte("Hello")}, commands: "/push/ripemd160/len/swap/pop/ripemd160-len/eq", result: []byte("Hello")},

	// compression
	{name: "snappy", initialStack: [][]byte{[]byte("This is some data we might compress")}, commands: "/snappy/unsnappy", result: []byte("This is some data we might compress")},
	{name: "snappy2", initialStack: [][]byte{[]byte("Hello this is a test")}, commands: "/snappy/hex", result: []byte("144c48656c6c6f207468697320697320612074657374")},
	{name: "zlib", initialStack: [][]byte{[]byte("This is some data we might compress")}, commands: "/zlib/unzlib", result: []byte("This is some data we might compress")},
	{name: "deflate", initialStack: [][]byte{[]byte("This is some data we might compress")}, commands: "/1/deflate/inflate", result: []byte("This is some data we might compress")},
	{name: "gzip", initialStack: [][]byte{[]byte("This is some data we might compress")}, commands: "/1/gzip/ungzip", result: []byte("This is some data we might compress")},
	{name: "lzw-lsb", initialStack: [][]byte{[]byte("This is some data we might compress")}, commands: "/8/lzw-lsb/8/unlzw-lsb", result: []byte("This is some data we might compress")},
	{name: "lzw-msb", initialStack: [][]byte{[]byte("This is some data we might compress")}, commands: "/8/lzw-msb/8/unlzw-msb", result: []byte("This is some data we might compress")},
	{name: "unbzip2", initialStack: [][]byte{[]byte("QlpoMTFBWSZTWdTk/8cAAAPTgAAQQAAEAC7i3IAgADFDTTAARNMJoMalXGTN1WHLwFdlq1ZnNsE2TiAnxdyRThQkNTk/8cA=")}, commands: "/unbase64/unbzip2", result: []byte("This is some data we might compress\n")},

	// control
	{name: "pushpop", initialStack: [][]byte{[]byte("This is some data")}, commands: "/FOO/push/pop/pop", result: []byte("This is some data")},
	{name: "loadsave", initialStack: [][]byte{[]byte("This is some data")}, commands: "/FOO/save/FOO/load", result: []byte("This is some data")},
	{name: "len", initialStack: [][]byte{[]byte("This is some data")}, commands: "/len/swap/pop", result: []byte("17")},
	{name: "swap", initialStack: [][]byte{[]byte("This is some data")}, commands: "/FOO/swap/pop", result: []byte("FOO")},
	{name: "slice1", initialStack: [][]byte{[]byte("ABCD")}, commands: "/-1/-1/slice", result: []byte("ABCD")},
	{name: "slice2", initialStack: [][]byte{[]byte("ABCD")}, commands: "/0/-1/slice", result: []byte("ABCD")},
	{name: "slice3", initialStack: [][]byte{[]byte("ABCD")}, commands: "/-1/4/slice", result: []byte("ABCD")},
	{name: "slice4", initialStack: [][]byte{[]byte("ABCD")}, commands: "/0/4/slice", result: []byte("ABCD")},
	{name: "slice5", initialStack: [][]byte{[]byte("ABCD")}, commands: "/0/1/slice", result: []byte("A")},
	{name: "slice6", initialStack: [][]byte{[]byte("ABCD")}, commands: "/3/4/slice", result: []byte("D")},
	{name: "left1", initialStack: [][]byte{[]byte("ABCD")}, commands: "/4/left", result: []byte("ABCD")},
	{name: "left2", initialStack: [][]byte{[]byte("ABCD")}, commands: "/1/left", result: []byte("A")},
	{name: "left3", initialStack: [][]byte{[]byte("ABCD")}, commands: "/0/left", result: []byte("")},
	{name: "right1", initialStack: [][]byte{[]byte("ABCD")}, commands: "/4/right", result: []byte("ABCD")},
	{name: "right2", initialStack: [][]byte{[]byte("ABCD")}, commands: "/1/right", result: []byte("D")},
	{name: "right3", initialStack: [][]byte{[]byte("ABCD")}, commands: "/0/right", result: []byte("")},
	{name: "append", initialStack: [][]byte{[]byte("ABCD")}, commands: "/EFGH/append", result: []byte("ABCDEFGH")},
	{name: "snip", initialStack: [][]byte{[]byte("ABC")}, commands: "/-1/snip/-/swap/append/append", result: []byte("-ABC")},
	{name: "snip", initialStack: [][]byte{[]byte("ABC")}, commands: "/0/snip/-/swap/append/append", result: []byte("-ABC")},
	{name: "snip", initialStack: [][]byte{[]byte("ABC")}, commands: "/1/snip/-/swap/append/append", result: []byte("A-BC")},
	{name: "snip", initialStack: [][]byte{[]byte("ABC")}, commands: "/2/snip/-/swap/append/append", result: []byte("AB-C")},
	{name: "snip", initialStack: [][]byte{[]byte("ABC")}, commands: "/3/snip/-/swap/append/append", result: []byte("ABC-")},
	{name: "snip", initialStack: [][]byte{[]byte("ABC")}, commands: "/4/snip/-/swap/append/append", result: []byte("ABC-")},
	{name: "eq", initialStack: [][]byte{[]byte("ABC")}, commands: "/DEF/DEF/eq", result: []byte("ABC")},
	{name: "neq", initialStack: [][]byte{[]byte("ABC")}, commands: "/DEF/EFG/neq", result: []byte("ABC")},

	// call
	{name: "call twofish", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-twofish/call/decrypt-twofish/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call blowfish", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-blowfish/call/decrypt-blowfish/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call aes", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-aes/call/decrypt-aes/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call des", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-des/call/decrypt-des/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call 3des", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-3des/call/decrypt-3des/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call hmac-md5", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/hash-hmac-md5/call/hex", result: []byte("7b9421e56f37e2bbce23a96fb71c02d2")},
	{name: "call hmac-sha1", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/hash-hmac-sha1/call/hex", result: []byte("5bbf8283b20d5c370e7d63696d3e98f9cd6926c0")},
	{name: "call hmac-sha224", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/hash-hmac-sha224/call/hex", result: []byte("02464ba9f03c08cc5297226fd04f2682194501ae226a413c610b4d55")},
	{name: "call hmac-sha256", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/hash-hmac-sha256/call/hex", result: []byte("bc57d22ccf1453762434c26319fa996683fbe6c9a1c85bb7779adbc59d643c76")},
	{name: "call hmac-sha384", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/hash-hmac-sha384/call/hex", result: []byte("7c9f2c1482e7de7328a8091e09f93538d7534a2ca7c9660438ef0f76b7a38c92f41aa9f428988632922d0996a8e0b4ed")},
	{name: "call hmac-sha512", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/hash-hmac-sha512/call/hex", result: []byte("d537850403eedb5220f7813979af7d48c8d9ebef7bb7111e99a7be528cc052488a164e153f37ff24909075dbb7d689ec379661633345a7fdbd5cec8140ec9183")},
	{name: "call hmac-ripemd160", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/hash-hmac-ripemd160/call/hex", result: []byte("eff69d3889ef5a8b4770370e21aa3e0b5c59f9f6")},

	{name: "call twofish signed", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-sign-twofish/call/decrypt-sign-twofish/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call blowfish signed", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-sign-blowfish/call/decrypt-sign-blowfish/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call aes signed", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-sign-aes/call/decrypt-sign-aes/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call des signed", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-sign-des/call/decrypt-sign-des/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{name: "call 3des signed", initialStack: [][]byte{[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, commands: "/encrypt-sign-3des/call/decrypt-sign-3des/call", result: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},

	// encoding
	{name: "hex", initialStack: [][]byte{[]byte("This is some data we might encode")}, commands: "/hex/unhex", result: []byte("This is some data we might encode")},
	{name: "ascii85", initialStack: [][]byte{[]byte("This is some data we might encode")}, commands: "/ascii85/unascii85", result: []byte("This is some data we might encode")},
	{name: "base32", initialStack: [][]byte{[]byte("This is some data we might encode")}, commands: "/base32/unbase32", result: []byte("This is some data we might encode")},
	{name: "base32-hex", initialStack: [][]byte{[]byte("This is some data we might encode")}, commands: "/base32-hex/unbase32-hex", result: []byte("This is some data we might encode")},
	{name: "base64", initialStack: [][]byte{[]byte("This is some data we might encode")}, commands: "/base64/unbase64", result: []byte("This is some data we might encode")},
	{name: "base64-url", initialStack: [][]byte{[]byte("This is some data we might encode")}, commands: "/base64-url/unbase64-url", result: []byte("This is some data we might encode")},

	// crypto
	{name: "aes-cfb", initialStack: [][]byte{}, commands: "/aes-blocksize/rand/push/ABCDEF/swap/mykey/md5/aes-cfb/swap/mykey/md5/unaes-cfb", result: []byte("ABCDEF")},
	{name: "aes-ofb", initialStack: [][]byte{}, commands: "/aes-blocksize/rand/push/ABCDEF/swap/mykey/md5/aes-ofb/swap/mykey/md5/aes-ofb", result: []byte("ABCDEF")},
	{name: "aes-ctr", initialStack: [][]byte{}, commands: "/aes-blocksize/rand/push/ABCDEF/swap/mykey/md5/aes-ctr/swap/mykey/md5/aes-ctr", result: []byte("ABCDEF")},
	{name: "aes-blocksize", initialStack: [][]byte{}, commands: "/aes-blocksize", result: []byte("16")},
	{name: "des-cfb", initialStack: [][]byte{}, commands: "/des-blocksize/rand/push/ABCDEF/swap/mykey/md5/8/left/des-cfb/swap/mykey/md5/8/left/undes-cfb", result: []byte("ABCDEF")},
	{name: "des-ofb", initialStack: [][]byte{}, commands: "/des-blocksize/rand/push/ABCDEF/swap/mykey/md5/8/left/des-ofb/swap/mykey/md5/8/left/des-ofb", result: []byte("ABCDEF")},
	{name: "des-ctr", initialStack: [][]byte{}, commands: "/des-blocksize/rand/push/ABCDEF/swap/mykey/md5/8/left/des-ctr/swap/mykey/md5/8/left/des-ctr", result: []byte("ABCDEF")},
	{name: "3des-cfb", initialStack: [][]byte{}, commands: "/des-blocksize/rand/push/ABCDEF/swap/mykey/md5/push/8/left/append/3des-cfb/swap/mykey/md5/push/8/left/append/un3des-cfb", result: []byte("ABCDEF")},
	{name: "3des-ofb", initialStack: [][]byte{}, commands: "/des-blocksize/rand/push/ABCDEF/swap/mykey/md5/push/8/left/append/3des-ofb/swap/mykey/md5/push/8/left/append/3des-ofb", result: []byte("ABCDEF")},
	{name: "3des-ctr", initialStack: [][]byte{}, commands: "/des-blocksize/rand/push/ABCDEF/swap/mykey/md5/push/8/left/append/3des-ctr/swap/mykey/md5/push/8/left/append/3des-ctr", result: []byte("ABCDEF")},
	{name: "des-blocksize", initialStack: [][]byte{}, commands: "/des-blocksize", result: []byte("8")},
	{name: "blowfish-cfb", initialStack: [][]byte{}, commands: "/blowfish-blocksize/rand/push/ABCDEF/swap/mykey/sha1/blowfish-cfb/swap/mykey/sha1/unblowfish-cfb", result: []byte("ABCDEF")},
	{name: "blowfish-ofb", initialStack: [][]byte{}, commands: "/blowfish-blocksize/rand/push/ABCDEF/swap/mykey/sha1/blowfish-ofb/swap/mykey/sha1/blowfish-ofb", result: []byte("ABCDEF")},
	{name: "blowfish-ctr", initialStack: [][]byte{}, commands: "/blowfish-blocksize/rand/push/ABCDEF/swap/mykey/sha1/blowfish-ctr/swap/mykey/sha1/blowfish-ctr", result: []byte("ABCDEF")},
	{name: "blowfish-salt-cfb", initialStack: [][]byte{}, commands: "/blowfish-blocksize/rand/push/ABCDEF/swap/mykey/sha1/345/blowfish-salt-cfb/swap/mykey/sha1/345/unblowfish-salt-cfb", result: []byte("ABCDEF")},
	{name: "blowfish-salt-ofb", initialStack: [][]byte{}, commands: "/blowfish-blocksize/rand/push/ABCDEF/swap/mykey/sha1/345/blowfish-salt-ofb/swap/mykey/sha1/345/blowfish-salt-ofb", result: []byte("ABCDEF")},
	{name: "blowfish-salt-ctr", initialStack: [][]byte{}, commands: "/blowfish-blocksize/rand/push/ABCDEF/swap/mykey/sha1/345/blowfish-salt-ctr/swap/mykey/sha1/345/blowfish-salt-ctr", result: []byte("ABCDEF")},
	{name: "blowfish-blocksize", initialStack: [][]byte{}, commands: "/blowfish-blocksize", result: []byte("8")},
	{name: "twofish-cfb", initialStack: [][]byte{}, commands: "/twofish-blocksize/rand/push/ABCDEF/swap/mykey/sha256/twofish-cfb/swap/mykey/sha256/untwofish-cfb", result: []byte("ABCDEF")},
	{name: "twofish-ofb", initialStack: [][]byte{}, commands: "/twofish-blocksize/rand/push/ABCDEF/swap/mykey/sha256/twofish-ofb/swap/mykey/sha256/twofish-ofb", result: []byte("ABCDEF")},
	{name: "twofish-ctr", initialStack: [][]byte{}, commands: "/twofish-blocksize/rand/push/ABCDEF/swap/mykey/sha256/twofish-ctr/swap/mykey/sha256/twofish-ctr", result: []byte("ABCDEF")},
	{name: "twofish-blocksize", initialStack: [][]byte{}, commands: "/twofish-blocksize", result: []byte("16")},
}

func TestEngine(t *testing.T) {
	eng := New()
	for _, testCase := range testCases {
		// Initialize engine and initial value
		eng.Reset()
		for i, entry := range testCase.initialStack {
			eng.PushStack(entry)
			if i == 0 {
				eng.SetVariable("body", entry)
			}
		}

		// process commands
		arr := strings.Split(strings.TrimPrefix(testCase.commands, "/"), "/")
		if len(arr) == 1 && arr[0] == "" {
			arr = make([]string, 0)
		}
		b, err := eng.Run(arr)
		if err != nil {
			t.Errorf("test case \"%s\": %s", testCase.name, err.Error())
		}

		if bytes.Compare(testCase.result, b) != 0 {
			t.Errorf("test case \"%s\": unexpected result\n\t\texpected: %v\n\t\tgot:      %v", testCase.name, testCase.result, b)
		}
	}
}
