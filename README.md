![hashsrv](media/logo.png)

**hashsrv** is a web service that performs hashing, encryption, encoding, and compression.

[![Go Reference](https://pkg.go.dev/badge/github.com/ancientlore/hashsrv.svg)](https://pkg.go.dev/github.com/ancientlore/hashsrv)

A [configuration file](hashsrv.config) in [TOML](https://github.com/mojombo/toml) format is used to set up hashsrv,
but environment variables and command-line options may be used as well.

Using hashsrv
-------------

hashsrv URLs are composed of commands that describe what to do with the given data. For instance, posting data to:

	/md5/hex

will calculate the MD5 hash of the posted data, convert it to hex encoding, and respond with the result.

hashsrv implements a simple processing engine that has a stack and a dictionary to store variables. Initially, the data posted via HTTP is pushed onto the stack. Most operations consume data from the stack and push their results onto the stack.

Additional arguments to operations can be placed onto the stack as literals. For instance, to generate 20 bytes of cryptographically random data and convert it to base64, use:

	/20/rand/base64

You should issue a GET request for that because no POST data is required.

Items in the URL that are not keywords are placed onto the stack. At the end of the list of commands, the stack should have a single value to use as the result of the request, or else an error occurs.

Named variables can be saved and loaded from a dictionary. See the load and save commands. The dictionary is initialized with HTTP headers that begin with `Hashsrv-` (with the prefix removed). So, to pass a variable called `key` into the dictionary, you can send an HTTP header called `Hashsrv-Key`.

As a convenience, the dictionary is initialized with the following values:

* body - the original request body
* key - initialized with a default key
* A number of standard combinations that you can invoke with the `call` command.

### Debug Mode

To output a debug view instead of the result, add ?debug=1 to the URL.

[Try It!](http://served.ancientlore.io:8080/Hello%20World/32/rand/md5/hmac-md5/hex?debug=1)

### Hash Functions

Command          | Stack in     | Stack out   | Description
-----------------|--------------|-------------|--------------------------------------------------------------------------------
md5              | Data         | Hash        | Hashes data using [MD5](http://golang.org/pkg/crypto/md5/)
sha1             | Data         | Hash        | Hashes data using [SHA1](http://golang.org/pkg/crypto/sha1/)
sha224           | Data         | Hash        | Hashes data using [SHA224](http://golang.org/pkg/crypto/sha256/)
sha256           | Data         | Hash        | Hashes data using [SHA256](http://golang.org/pkg/crypto/sha256/)
sha384           | Data         | Hash        | Hashes data using [SHA384](http://golang.org/pkg/crypto/sha512/)
sha512           | Data         | Hash        | Hashes data using [SHA512](http://golang.org/pkg/crypto/sha512/)
ripemd160        | Data         | Hash        | Hashes data using [RIPEMD160](http://golang.org/x/crypto/ripemd160)
hmac-md5         | Data, Key    | Hash        | [HMAC](http://golang.org/pkg/crypto/hmac/) hashes data using MD5
hmac-sha1        | Data, Key    | Hash        | [HMAC](http://golang.org/pkg/crypto/hmac/) hashes data using SHA1
hmac-sha224      | Data, Key    | Hash        | [HMAC](http://golang.org/pkg/crypto/hmac/) hashes data using SHA2 224-bit
hmac-sha256      | Data, Key    | Hash        | [HMAC](http://golang.org/pkg/crypto/hmac/) hashes data using SHA2 256-bit
hmac-sha384      | Data, Key    | Hash        | [HMAC](http://golang.org/pkg/crypto/hmac/) hashes data using SHA2 384-bit
hmac-sha512      | Data, Key    | Hash        | [HMAC](http://golang.org/pkg/crypto/hmac/) hashes data using SHA2 512-bit
hmac-ripemd160   | Data, Key    | Hash        | [HMAC](http://golang.org/pkg/crypto/hmac/) hashes data using RIPEMD160
md5-len          |              | 16          | Returns the number of bytes for MD5
sha1-len         |              | 20          | Returns the number of bytes for  SHA1
sha224-len       |              | 28          | Returns the number of bytes for SHA224
sha256-len       |              | 32          | Returns the number of bytes for SHA256
sha384-len       |              | 48          | Returns the number of bytes for SHA384
sha512-len       |              | 64          | Returns the number of bytes for SHA512
ripemd160-len    |              | 20          | Returns the number of bytes for  RIPEMD160

**Note:** When using HMAC, it is customary to hash the key using the same hash function defined for that version of HMAC. You must do that yourself. For instance, when using hmac-sha256, the key should be hashed with sha256 and then used for HMAC.

### Encoding Functions

Command          | Stack in     | Stack out   | Description
-----------------|--------------|-------------|--------------------------------------------------------------------------------
hex              | Data         | EncodedData | Encode the data as [hex](http://golang.org/pkg/encoding/hex/)
unhex            | EncodedData  | Data        | Decode the data as [hex](http://golang.org/pkg/encoding/hex/)
ascii85          | Data         | EncodedData | Encode data as [ASCII-85](http://golang.org/pkg/encoding/ascii85/)
unascii85        | EncodedData  | Data        | Decode data as [ASCII-85](http://golang.org/pkg/encoding/ascii85/)
base32           | Data         | EncodedData | Encode data as [BASE-32](http://golang.org/pkg/encoding/base32/)
unbase32         | EncodedData  | Data        | Decode data as [BASE-32](http://golang.org/pkg/encoding/base32/)
base32-hex       | Data         | EncodedData | Encode data as [BASE-32 Hex](http://golang.org/pkg/encoding/base32/)
unbase32-hex     | EncodedData  | Data        | Decode data as [BASE-32 Hex](http://golang.org/pkg/encoding/base32/)
base64           | Data         | EncodedData | Encode data as [BASE-64](http://golang.org/pkg/encoding/base64/)
unbase64         | EncodedData  | Data        | Decode data as [BASE-64](http://golang.org/pkg/encoding/base64/)
base64-url       | Data         | EncodedData | Encode data as [BASE-64 URL](http://golang.org/pkg/encoding/base64/)
unbase64-url     | EncodedData  | Data        | Decode data as [BASE-64 URL](http://golang.org/pkg/encoding/base64/)

### Checksum Functions

Command          | Stack in     | Stack out   | Description
-----------------|--------------|-------------|--------------------------------------------------------------------------------
adler32          | Data         | Checksum    | Compute the [Adler-32](http://golang.org/pkg/hash/adler32/) checksum
crc32            | Data         | Checksum    | Compute the [CRC-32](http://golang.org/pkg/hash/crc32/) checksum using the IEEE polynomial
crc32-ieee       | Data         | Checksum    | Compute the [CRC-32](http://golang.org/pkg/hash/crc32/) checksum using the IEEE polynomial
crc32-castagnoli | Data         | Checksum    | Compute the [CRC-32](http://golang.org/pkg/hash/crc32/) checksum using the Castagnoli polynomial
crc32-koopman    | Data         | Checksum    | Compute the [CRC-32](http://golang.org/pkg/hash/crc32/) checksum using the Koopman polynomial
crc64-iso        | Data         | Checksum    | Compute the [CRC-64](http://golang.org/pkg/hash/crc64/) checksum using the ISO polynomial
crc64-ecma       | Data         | Checksum    | Compute the [CRC-64](http://golang.org/pkg/hash/crc64/) checksum using the ECMA polynomial
fnv32            | Data         | Hash        | Compute the [FNV-1](http://golang.org/pkg/hash/fnv/) non-cryptographic hash for 32-bits
fnv32a           | Data         | Hash        | Compute the [FNV-1a](http://golang.org/pkg/hash/fnv/) non-cryptographic hash for 32-bits
fnv64            | Data         | Hash        | Compute the [FNV-1](http://golang.org/pkg/hash/fnv/) non-cryptographic hash for 64-bits
fnv64a           | Data         | Hash        | Compute the [FNV-1a](http://golang.org/pkg/hash/fnv/) non-cryptographic hash for 64-bits

### Compression Functions

Command          | Stack in     | Stack out   | Description
-----------------|--------------|-------------|--------------------------------------------------------------------------------
snappy           | Data         | Compressed  | Compresses data using the [Snappy](https://code.google.com/p/snappy/) algorithm
unsnappy         | Compressed   | Data        | Decompresses data using the [Snappy](https://code.google.com/p/snappy/) algorithm
zlib             | Data         | Compressed  | Compresses data using the [zlib](http://golang.org/pkg/compress/zlib/) algorithm
unzlib           | Compressed   | Data        | Decompresses data using the [zlib](http://golang.org/pkg/compress/zlib/) algorithm
deflate          | Data, Factor | Compressed  | Compresses data using the [flate](http://golang.org/pkg/compress/flate/) algorithm - stack contains a compression factor where -1 is default and 0-9 controls compression (0 is none, and 9 is the most)
inflate          | Compressed   | Data        | Decompresses data using the [flate](http://golang.org/pkg/compress/flate/) algorithm
gzip             | Data, Factor | Compressed  | Compresses data using the [gzip](http://golang.org/pkg/compress/gzip/) algorithm - stack contains a compression factor where -1 is default, 0 is none, 1 is best speed, and 9 is best size
ungzip           | Compressed   | Data        | Decompresses data using the [gzip](http://golang.org/pkg/compress/gzip/) algorithm
unbzip2          | Compressed   | Data        | Decompresses data using the [bzip2](http://golang.org/pkg/compress/bzip2/) algorithm
lzw-lsb          | Data, Bits   | Compressed  | Compresses data using the [lzw](http://golang.org/pkg/compress/lzw/) algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses least significant bit ordering as used in the GIF file format.
unlzw-lsb        | Compressed, Bits | Data        | Decompresses data using the [lzw](http://golang.org/pkg/compress/lzw/) algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses least significant bit ordering as used in the GIF file format.
lzw-msb          | Data, Bits   | Compressed  | Compresses data using the [lzw](http://golang.org/pkg/compress/lzw/) algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses most significant bit ordering as used in the TIFF and PDF file formats.
unlzw-msb        | Compressed, Bits | Compressed  | Decompresses data using the [lzw](http://golang.org/pkg/compress/lzw/) algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses most significant bit ordering as used in the TIFF and PDF file formats.

### Control Functions

Command          | Stack in     | Stack out   | Description
-----------------|--------------|-------------|--------------------------------------------------------------------------------
push             | Data         | Data, Data  | Duplicates the value on the top of the stack
pop              | Data         |             | Pops the value off the top of the stack (effectively discarding)
load             | Name         | Value       | Pushes a named value from the dictinary onto the stack
save             | Value, Name  |             | Pops a value from the stack and places it into the dictionary
swap             | Val1, Val2   | Val2, Val1  | Swaps the two values at the top of the stack
append           | Val1, Val2   | Appended    | Appends the value on the top of the stack to the previous value on the stack
slice            | Data         | SliceOfData | Slices the value on the stack, taking elements from start to end on the stack. Use -1 for values from the beginning or end. One example is `/9/20/slice` which takes elements 9 through 19, or `/2/-1/slice` which takes elements 2 through the end.
len              | Data         | Data, Length| Pushes the length of the value on the stack in bytes onto the stack
left             | Data, Count  | SliceOfData | Takes the leftmost bytes of data
right            | Data, Count  | SliceOfData | Takes the rightmost bytes of data
snip             | Data, Pos    | Data1, Data2| Snips the data in half at the given position, resulting in two values on the stack
eq               | Data1, Data2 |             | Fails the command unless the two data elements are equal
neq              | Data1, Data2 |             | Fails the command unless the two data elements are not equal
call             | Name         | (Varies)    | Loads the named value from the dictionary and executes the commands contained there (formatted like normal - /md5/hex for example)

### Crypto Functions

Command          | Stack in     | Stack out   | Description
-----------------|--------------|-------------|--------------------------------------------------------------------------------
rand             | Count        | Data        | Generates cryptographically [random](http://golang.org/pkg/crypto/rand/) bytes given the count on the stack
aes-blocksize    |              | 16          | Pushes the AES block size on the stack
aes-cfb          | Data, IV, Key| Data        | Encrypts data using the given IV and 16-byte Key, placing the ciphertext back on the stack. Uses [AES](http://golang.org/pkg/crypto/aes/) encryption and the CFB block mode.
unaes-cfb        | Data, IV, Key| Data        | Decrypts data using the given IV and 16-byte Key, placing the plaintext back on the stack. Uses [AES](http://golang.org/pkg/crypto/aes/) encryption and the CFB block mode.
aes-ofb          | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 16-byte Key, placing the result back on the stack. Uses [AES](http://golang.org/pkg/crypto/aes/) encryption and the OFB block mode.
aes-ctr          | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 16-byte Key, placing the result back on the stack. Uses [AES](http://golang.org/pkg/crypto/aes/) encryption and the CTR block mode.
des-blocksize    |              | 8           | Pushes the DES block size on the stack
des-cfb          | Data, IV, Key| Data        | Encrypts data using the given IV and 8-byte Key, placing the ciphertext back on the stack. Uses [DES](http://golang.org/pkg/crypto/des/) encryption and the CFB block mode.
undes-cfb        | Data, IV, Key| Data        | Decrypts data using the given IV and 8-byte Key, placing the plaintext back on the stack. Uses [DES](http://golang.org/pkg/crypto/des/) encryption and the CFB block mode.
des-ofb          | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 8-byte Key, placing the result back on the stack. Uses [DES](http://golang.org/pkg/crypto/des/) encryption and the OFB block mode.
des-ctr          | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 8-byte Key, placing the result back on the stack. Uses [DES](http://golang.org/pkg/crypto/des/) encryption and the CTR block mode.
3des-blocksize   |              | 8           | Pushes the Triple DES block size on the stack
3des-cfb         | Data, IV, Key| Data        | Encrypts data using the given IV and 24-byte Key, placing the ciphertext back on the stack. Uses [Triple DES](http://golang.org/pkg/crypto/des/) encryption and the CFB block mode.
un3des-cfb       | Data, IV, Key| Data        | Decrypts data using the given IV and 24-byte Key, placing the plaintext back on the stack. Uses [Triple DES](http://golang.org/pkg/crypto/des/) encryption and the CFB block mode.
3des-ofb         | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 24-byte Key, placing the result back on the stack. Uses [Triple DES](http://golang.org/pkg/crypto/des/) encryption and the OFB block mode.
3des-ctr         | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 24-byte Key, placing the result back on the stack. Uses [Triple DES](http://golang.org/pkg/crypto/des/) encryption and the CTR block mode.
blowfish-blocksize|              | 8           | Pushes the blowfish block size on the stack
blowfish-cfb     | Data, IV, Key| Data        | Encrypts data using the given IV and 1 to 56-byte Key, placing the ciphertext back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the CFB block mode.
unblowfish-cfb   | Data, IV, Key| Data        | Decrypts data using the given IV and 1 to 56-byte Key, placing the plaintext back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the CFB block mode.
blowfish-ofb     | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the OFB block mode.
blowfish-ctr     | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the CTR block mode.
blowfish-salt-cfb  | Data, IV, Key, Salt| Data        | Encrypts data using the given IV and 1 to 56-byte Key, placing the ciphertext back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the CFB block mode.
unblowfish-salt-cfb| Data, IV, Key, Salt| Data        | Decrypts data using the given IV and 1 to 56-byte Key, placing the plaintext back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the CFB block mode.
blowfish-salt-ofb  | Data, IV, Key, Salt| Data        | Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the OFB block mode.
blowfish-salt-ctr  | Data, IV, Key, Salt| Data        | Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses [Blowfish](https://godoc.org/golang.org/x/crypto/blowfish) encryption and the CTR block mode.
twofish-blocksize|              | 16          | Pushes the twofish block size on the stack
twofish-cfb      | Data, IV, Key| Data        | Encrypts data using the given IV and 16, 24, or 32-byte Key, placing the ciphertext back on the stack. Uses [Twofish](https://godoc.org/golang.org/x/crypto/twofish) encryption and the CFB block mode.
untwofish-cfb    | Data, IV, Key| Data        | Decrypts data using the given IV and 16, 24, or 32-byte Key, placing the plaintext back on the stack. Uses [Twofish](https://godoc.org/golang.org/x/crypto/twofish) encryption and the CFB block mode.
twofish-ofb      | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 16, 24, or 32-byte Key, placing the result back on the stack. Uses [Twofish](https://godoc.org/golang.org/x/crypto/twofish) encryption and the OFB block mode.
twofish-ctr      | Data, IV, Key| Data        | Encrypts or decrypts data using the given IV and 16, 24, or 32-byte Key, placing the result back on the stack. Uses [Twofish](https://godoc.org/golang.org/x/crypto/twofish) encryption and the CTR block mode.

#### Notes on encryption

The initialization vector (IV) is used by many routines. It does not need to be kept secure, but it should generally be random and different for each different encryption run. It can easily be generated with the `rand` function. However, you need to keep it for decryption. It is customary to put it at the beginning of the encrypted data. *These routines don't do that for you.*

Each encryption routine supports several [block modes](http://golang.org/pkg/crypto/cipher/). Some of the block modes are symmetrical - so you use the same function to encrypt and decrypt. Others are not.

Some routines require fixed key sizes, others are variable. Keys can be any data. It is usually considered more secure when these keys are relatively random or hashed.


### On the todo list

* Control - loop (to go through lines of text and do batch operations)
* Specialized - protect, unprotect

Examples
--------

URL                                                     | Result
--------------------------------------------------------|-----------------------------------------------------------------
POST /                                                  | Returns what you posted
POST /sha256                                            | Returns SHA256 hash as binary data
POST /sha256/hex                                        | Returns SHA256 hash as hex encoding
POST /unhex/snappy/hex                                  | Decodes hex data, compresses it using Snappy, and encodes the result to hex
GET /Hello%20World/32/rand/md5/hmac-md5/hex        | Pushes "Hello World" on the stach, generates 32 bytes of random data as the HMAC key (which is then hashed with md5), computes the HMAC-MD5 hash, and converts the result to hex. [Try It!](http://served.ancientlore.io:8080/Hello%20World/32/rand/md5/hmac-md5/hex)
POST /MyKeyHere/sha512/hmac-sha512/base64-url                  | Hashes the data with HMAC-SHA512 using the the sha512 hash of the key "MyKeyHere" and returns it as base64.

Running hashsrv
---------------

All you need is your configuration file and the hashsrv binary for your platform. You can run it manually or as a service on Windows or Linux (see below).

To test with the default [configuration file](hashsrv.config) and logging:

	./hashsrv -run

To run in the background:

	./hashsrv -run &

To keep it running after you log off:

	nohup ./hashsrv -run &

`nohup` is a Linux utility that keeps a process going after you log off.

### Installation

The only required files are the hashsrv binary and the configuration file. The hashsrv has minimal dependencies - just a few shared libraries that should already be on the operating system.

On all operating systems, you may override the configuration file location using the `HASHSRV_CONFIG` environment variable or the `-config` command-line option, which takes precedence. See below for where to place the configuration file when none of these are present.

The location of the configuration file is based on the location of the hashsrv binary. `/usr/bin` and `/bin` locations are replaced with `/etc` - so effectively, the configuration file is located in the `etc` folder that corresponds to the `bin` folder. If the binary is not in a `bin` folder, then the configuration file is expected to be in the same folder as the binary. Some examples are shown below.

| hashsrv binary location       | Default configuration file location  |
|-------------------------------|--------------------------------------|
| /bin/hashsrv                  | /etc/hashsrv.config                  |
| /usr/bin/hashsrv              | /etc/hashsrv.config                  |
| /usr/local/bin/hashsrv        | /usr/local/etc/hashsrv.config        |
| /usr/local/bin/foo/hashsrv    | /usr/local/etc/foo/hashsrv.config    |
| /usr/local/foo/bin/hashsrv    | /usr/local/foo/etc/hashsrv.config    |
| c:\hashsrv\bin\hashsrv.exe    | c:\hashsrv\etc\hashsrv.config        |
| c:\files\hashsrv.exe          | c:\files\hashsrv.config              |
| /home/michael/hashsrv         | /home/michael/hashsrv.config         |

### Running as a service

You can install the hashsrv as a service on Windows or Linux with Upstart. Use the `-install` and `-remove` options to install or remove the hashsrv.

#### Linux

On Linux, the -install option created a `HashSrv.conf` file in `/etc/init`. To start or stop the hashsrv, you can use:

	sudo start HashSrv
	sudo stop HashSrv

If the service doesn't start, most likely the configuration file has a problem.

#### Windows

On Windows, the hashsrv uses the Service API. Use the Service administration tool to start or stop the hashsrv.

Also, you will need to use the `-run` option if you want to run the application standalone (not as a service).

### Environment Variables

| Option         | Default                    | Description                                              |
|----------------|----------------------------|----------------------------------------------------------|
| HASHSRV_CONFIG | hashsrv.config (see above) | Specifies the default location of the configuration file |


### Command-Line Parameters

| Option      | Default                             | Description                               | 
|-------------|-------------------------------------|-------------------------------------------|
| -addr       | ":9009"                             | Address to serve                          |
| -config     | HASHSRV_CONFIG environment variable | Use to override the configuration file    |
| -cpuprofile |                                     | Write CPU profile to file                 |
| -memprofile |                                     | Write memory profile to file              |
| -help       | false                               | Show command help                         |
| -noisy      | false                               | Enable logging                            |
| -install    | false                               | Install hashsrv as a service              |
| -remove     | false                               | Remove the hashsrv service                |
| -run        | false                               | Run hashsrv standalone (not as a service) |
| -start      | false                               | Start the hashsrv service                 |
| -stop       | false                               | Stop the hashsrv service                  |


### Configuration File Parameters

| Option | Default | Description        |
|--------|---------|--------------------|
| addr   | ":9009" | Web server address |

