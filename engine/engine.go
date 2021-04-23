package engine

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	defaultKey = "& fri3d Gr33n tomat0s wiTh g0rill4_glu3 sauce!"
)

// funcInfo stores information about the function definitions
type funcInfo struct {
	f    func() error
	In   string
	Out  string
	Desc string
}

// The Engine is the processing logic of the hash server
type Engine struct {
	stack     *Stack
	values    map[string][]byte
	funcMap   map[string]funcInfo
	logBuf    *bytes.Buffer
	DebugMode bool
}

// New creates a new engine
func New() *Engine {
	e := new(Engine)
	e.logBuf = new(bytes.Buffer)
	e.initMap()
	e.Reset()
	return e
}

// Reset returns the engine to the initial state
func (e *Engine) Reset() {
	e.stack = NewStack()
	e.values = make(map[string][]byte)
	e.logBuf.Reset()
	e.DebugMode = false
	e.initVars()
}

// initVars sets the default variables
func (e *Engine) initVars() {
	var hmacFmt = func(alg string) []byte { return []byte(fmt.Sprintf("/key/load/%s/hmac-%s", alg, alg)) }

	var encryptFmt = func(alg, hash string) []byte {
		return []byte(fmt.Sprintf("/%s-blocksize/rand/push/iv/save/key/load/%s/%s-cfb/iv/load/swap/append", alg, hash, alg))
	}
	var decryptFmt = func(alg, hash string) []byte {
		return []byte(fmt.Sprintf("/%s-blocksize/snip/swap/key/load/%s/un%s-cfb", alg, hash, alg))
	}

	var signFmt = func(alg string) []byte { return []byte(fmt.Sprintf("/push/hash-hmac-%s/call/swap/append", alg)) }
	var checksigFmt = func(alg string) []byte {
		return []byte(fmt.Sprintf("/%s-len/snip/swap/temp/save/push/hash-hmac-%s/call/temp/load/eq", alg, alg))
	}

	// Add the default key
	e.SetVariable("key", []byte(defaultKey))

	// encrypt and sign the data
	e.SetVariable("encrypt-twofish", encryptFmt("twofish", "sha256"))
	e.SetVariable("decrypt-twofish", decryptFmt("twofish", "sha256"))
	e.SetVariable("encrypt-blowfish", encryptFmt("blowfish", "sha256"))
	e.SetVariable("decrypt-blowfish", decryptFmt("blowfish", "sha256"))
	e.SetVariable("encrypt-aes", encryptFmt("aes", "md5"))
	e.SetVariable("decrypt-aes", decryptFmt("aes", "md5"))
	e.SetVariable("encrypt-des", encryptFmt("des", "crc64-iso"))
	e.SetVariable("decrypt-des", decryptFmt("des", "crc64-iso"))
	e.SetVariable("encrypt-3des", encryptFmt("3des", "push/md5/swap/crc64-iso/append"))
	e.SetVariable("decrypt-3des", decryptFmt("3des", "push/md5/swap/crc64-iso/append"))

	// Add shortcuts that use the value stored in "key" to do an HMAC
	e.SetVariable("hash-hmac-md5", hmacFmt("md5"))
	e.SetVariable("hash-hmac-sha1", hmacFmt("sha1"))
	e.SetVariable("hash-hmac-sha224", hmacFmt("sha224"))
	e.SetVariable("hash-hmac-sha256", hmacFmt("sha256"))
	e.SetVariable("hash-hmac-sha384", hmacFmt("sha384"))
	e.SetVariable("hash-hmac-sha512", hmacFmt("sha512"))
	e.SetVariable("hash-hmac-ripemd160", hmacFmt("ripemd160"))

	// sign and check signature
	e.SetVariable("sign-md5", signFmt("md5"))
	e.SetVariable("checksig-md5", checksigFmt("md5"))
	e.SetVariable("sign-sha1", signFmt("sha1"))
	e.SetVariable("checksig-sha1", checksigFmt("sha1"))
	e.SetVariable("sign-sha224", signFmt("sha224"))
	e.SetVariable("checksig-sha224", checksigFmt("sha224"))
	e.SetVariable("sign-sha256", signFmt("sha256"))
	e.SetVariable("checksig-sha256", checksigFmt("sha256"))
	e.SetVariable("sign-sha384", signFmt("sha384"))
	e.SetVariable("checksig-sha384", checksigFmt("sha384"))
	e.SetVariable("sign-sha512", signFmt("sha512"))
	e.SetVariable("checksig-sha512", checksigFmt("sha512"))
	e.SetVariable("sign-ripemd160", signFmt("ripemd160"))
	e.SetVariable("checksig-ripemd160", checksigFmt("ripemd160"))

	// encrypt and sign
	e.SetVariable("encrypt-sign-twofish", []byte("/encrypt-twofish/call/sign-sha256/call"))
	e.SetVariable("decrypt-sign-twofish", []byte("/checksig-sha256/call/decrypt-twofish/call"))
	e.SetVariable("encrypt-sign-blowfish", []byte("/encrypt-blowfish/call/sign-sha256/call"))
	e.SetVariable("decrypt-sign-blowfish", []byte("/checksig-sha256/call/decrypt-blowfish/call"))
	e.SetVariable("encrypt-sign-aes", []byte("/encrypt-aes/call/sign-sha256/call"))
	e.SetVariable("decrypt-sign-aes", []byte("/checksig-sha256/call/decrypt-aes/call"))
	e.SetVariable("encrypt-sign-des", []byte("/encrypt-des/call/sign-sha256/call"))
	e.SetVariable("decrypt-sign-des", []byte("/checksig-sha256/call/decrypt-des/call"))
	e.SetVariable("encrypt-sign-3des", []byte("/encrypt-3des/call/sign-sha256/call"))
	e.SetVariable("decrypt-sign-3des", []byte("/checksig-sha256/call/decrypt-3des/call"))
}

// SetVariable sets a variable for the engine to use
func (e *Engine) SetVariable(name string, value []byte) {
	e.values[strings.ToLower(name)] = value
}

// GetVariable returns the value of a variable in the engine
func (e *Engine) GetVariable(name string) []byte {
	return e.values[strings.ToLower(name)]
}

// PushStack is used to intialize the stack of the engine. It is typically
// used to set the initial value to operate on.
func (e *Engine) PushStack(value []byte) {
	e.stack.Push(value)
}

// Run executes the logic of the engine, returning the last value from the stack.
func (e *Engine) Run(commands []string) ([]byte, error) {
	var err error

	if e.DebugMode {
		e.logBuf.Write([]byte(htmlHeader))
	}

	//e.LogValues()

	err = e.exec(commands)
	if err != nil {
		if e.DebugMode {
			e.Log(err)
		} else {
			return nil, err
		}
	}

	e.LogStack()

	b := e.stack.Pop()
	if b == nil {
		err = errors.New("nothing left on the stack to return")
		if e.DebugMode {
			e.Log(err)
		} else {
			return nil, err
		}
	}

	if e.stack.Len() > 0 {
		err = fmt.Errorf("%d unused items on stack", e.stack.Len())
		if e.DebugMode {
			e.Log(err)
		} else {
			return nil, err
		}
	}

	if e.DebugMode {
		e.logBuf.Write([]byte(htmlFooter))
		b = e.logBuf.Bytes()
	}

	return b, nil
}

// exec runs the commands for an already initialized engine
func (e *Engine) exec(commands []string) error {
	var err error
	e.Log("exec /", strings.Join(commands, "/"))
	for _, s := range commands {
		e.LogStack()
		fd, ok := e.funcMap[strings.TrimSpace(s)]
		if ok {
			e.Logf("(%s) -> %s -> (%s)", fd.In, s, fd.Out)
			err = fd.f()
		} else {
			e.Logf("push %+q", s)
			e.stack.Push([]byte(s))
		}
		if err != nil {
			return err
		}
	}
	e.Log("end")
	return nil
}

func (e *Engine) initMap() {
	e.funcMap = map[string]funcInfo{
		// hashing
		"md5":       {f: e.md5, In: "Data", Out: "Hash", Desc: "Hashes data using MD5"},
		"sha1":      {f: e.sha1, In: "Data", Out: "Hash", Desc: "Hashes data using SHA1"},
		"sha224":    {f: e.sha224, In: "Data", Out: "Hash", Desc: "Hashes data using SHA224"},
		"sha256":    {f: e.sha256, In: "Data", Out: "Hash", Desc: "Hashes data using SHA256"},
		"sha384":    {f: e.sha384, In: "Data", Out: "Hash", Desc: "Hashes data using SHA384"},
		"sha512":    {f: e.sha512, In: "Data", Out: "Hash", Desc: "Hashes data using SHA512"},
		"ripemd160": {f: e.ripemd160, In: "Data", Out: "Hash", Desc: "Hashes data using RIPEMD160"},
		"rand":      {f: e.rand, In: "Count", Out: "Data", Desc: "Generates cryptographically random bytes given the count on the stack"},

		"md5-len":       {f: e.md5_len, In: "", Out: "16", Desc: "Returns the number of bytes for MD5"},
		"sha1-len":      {f: e.sha1_len, In: "", Out: "20", Desc: "Returns the number of bytes for  SHA1"},
		"sha224-len":    {f: e.sha224_len, In: "", Out: "28", Desc: "Returns the number of bytes for  SHA224"},
		"sha256-len":    {f: e.sha256_len, In: "", Out: "32", Desc: "Returns the number of bytes for  SHA256"},
		"sha384-len":    {f: e.sha384_len, In: "", Out: "48", Desc: "Returns the number of bytes for  SHA384"},
		"sha512-len":    {f: e.sha512_len, In: "", Out: "64", Desc: "Returns the number of bytes for  SHA512"},
		"ripemd160-len": {f: e.ripemd160_len, In: "", Out: "20", Desc: "Returns the number of bytes for  RIPEMD160"},

		// HMAC hashing
		"hmac-md5":       {f: e.hmac_md5, In: "Data, Key", Out: "Hash", Desc: "HMAC hashes data using MD5"},
		"hmac-sha1":      {f: e.hmac_sha1, In: "Data, Key", Out: "Hash", Desc: "HMAC hashes data using SHA1"},
		"hmac-sha224":    {f: e.hmac_sha224, In: "Data, Key", Out: "Hash", Desc: "HMAC hashes data using SHA224"},
		"hmac-sha256":    {f: e.hmac_sha256, In: "Data, Key", Out: "Hash", Desc: "HMAC hashes data using SHA256"},
		"hmac-sha384":    {f: e.hmac_sha384, In: "Data, Key", Out: "Hash", Desc: "HMAC hashes data using SHA384"},
		"hmac-sha512":    {f: e.hmac_sha512, In: "Data, Key", Out: "Hash", Desc: "HMAC hashes data using SHA512"},
		"hmac-ripemd160": {f: e.hmac_ripemd160, In: "Data, Key", Out: "Hash", Desc: "HMAC hashes data using RIPEMD160"},

		// encoding
		"hex":          {f: e.hex, In: "Data", Out: "EncodedData", Desc: "Encode the data to hex"},
		"unhex":        {f: e.unhex, In: "EncodedData", Out: "Data", Desc: "Decode the data from hex"},
		"ascii85":      {f: e.ascii85, In: "Data", Out: "EncodedData", Desc: "Encode the data to ascii-85"},
		"unascii85":    {f: e.unascii85, In: "EncodedData", Out: "Data", Desc: "Decode the data from ascii-85"},
		"base32":       {f: e.base32, In: "Data", Out: "EncodedData", Desc: "Encode the data to base32"},
		"unbase32":     {f: e.unbase32, In: "EncodedData", Out: "Data", Desc: "Decode the data from base32"},
		"base32-hex":   {f: e.base32_hex, In: "Data", Out: "EncodedData", Desc: "Encode the data to base32 hex"},
		"unbase32-hex": {f: e.unbase32_hex, In: "EncodedData", Out: "Data", Desc: "Decode the data from base32 hex"},
		"base64":       {f: e.base64, In: "Data", Out: "EncodedData", Desc: "Encode the data to base64"},
		"unbase64":     {f: e.unbase64, In: "EncodedData", Out: "Data", Desc: "Decode the data from base64"},
		"base64-url":   {f: e.base64_url, In: "Data", Out: "EncodedData", Desc: "Encode the data to base64 url"},
		"unbase64-url": {f: e.unbase64_url, In: "EncodedData", Out: "Data", Desc: "Decode the data from base64 url"},

		// checksum hashing
		"adler32":          {f: e.adler32, In: "Data", Out: "Checksum", Desc: "Compute the Adler-32 checksum"},
		"crc32":            {f: e.crc32, In: "Data", Out: "Checksum", Desc: "Compute the CRC-32 checksum using the IEEE polynomial"},
		"crc32-ieee":       {f: e.crc32_ieee, In: "Data", Out: "Checksum", Desc: "Compute the CRC-32 checksum using the IEEE polynomial"},
		"crc32-castagnoli": {f: e.crc32_castagnoli, In: "Data", Out: "Checksum", Desc: "Compute the CRC-32 checksum using the Castagnoli polynomial"},
		"crc32-koopman":    {f: e.crc32_koopman, In: "Data", Out: "Checksum", Desc: "Compute the CRC-32 checksum using the Koopman polynomial"},
		"crc64-iso":        {f: e.crc64_iso, In: "Data", Out: "Checksum", Desc: "Compute the CRC-64 checksum using the ISO polynomial"},
		"crc64-ecma":       {f: e.crc64_ecma, In: "Data", Out: "Checksum", Desc: "Compute the CRC-64 checksum using the ECMA polynomial"},
		"fnv32":            {f: e.fnv32, In: "Data", Out: "Hash", Desc: "Compute the FNV-1 non-cryptographic hash for 32-bits"},
		"fnv32a":           {f: e.fnv32a, In: "Data", Out: "Hash", Desc: "Compute the FNV-1a non-cryptographic hash for 32-bits"},
		"fnv64":            {f: e.fnv64, In: "Data", Out: "Hash", Desc: "Compute the FNV-1 non-cryptographic hash for 64-bits"},
		"fnv64a":           {f: e.fnv64a, In: "Data", Out: "Hash", Desc: "Compute the FNV-1a non-cryptographic hash for 64-bits"},

		// Compression
		"snappy":    {f: e.snappy, In: "Data", Out: "Compressed", Desc: "Compresses data using the Snappy algorithm"},
		"unsnappy":  {f: e.unsnappy, In: "Compressed", Out: "Data", Desc: "Decompresses data using the Snappy algorithm"},
		"zlib":      {f: e.zlib, In: "Data", Out: "Compressed", Desc: "Compresses data using the zlib algorithm"},
		"unzlib":    {f: e.unzlib, In: "Compressed", Out: "Data", Desc: "Decompresses data using the zlib algorithm"},
		"deflate":   {f: e.deflate, In: "Data, Factor", Out: "Compressed", Desc: "Compresses data using the flate algorithm - stack contains a compression factor where -1 is default and 0-9 controls compression (0 is none, and 9 is the most)"},
		"inflate":   {f: e.inflate, In: "Compressed", Out: "Data", Desc: "Decompresses data using the flate algorithm"},
		"gzip":      {f: e.gzip, In: "Data, Factor", Out: "Compressed", Desc: "Compresses data using the gzip algorithm - stack contains a compression factor where -1 is default, 0 is none, 1 is best speed, and 9 is best size"},
		"ungzip":    {f: e.ungzip, In: "Compressed", Out: "Data", Desc: "Decompresses data using the gzip algorithm"},
		"unbzip2":   {f: e.unbzip2, In: "Compressed", Out: "Data", Desc: "Decompresses data using the bzip2 algorithm"},
		"lzw-msb":   {f: e.lzw_msb, In: "Data, Bits", Out: "Compressed", Desc: "Compresses data using the lzw algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses most significant bit ordering as used in the TIFF and PDF file formats."},
		"lzw-lsb":   {f: e.lzw_lsb, In: "Data, Bits", Out: "Compressed", Desc: "Compresses data using the lzw algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses least significant bit ordering as used in the GIF file format."},
		"unlzw-msb": {f: e.unlzw_msb, In: "Compressed, Bits", Out: "Data", Desc: "Decompresses data using the lzw algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses most significant bit ordering as used in the TIFF and PDF file formats."},
		"unlzw-lsb": {f: e.unlzw_lsb, In: "Compressed, Bits", Out: "Data", Desc: "Decompresses data using the lzw algorithm - stack contains the number of bits to use for literal codes, typically 8 but can be 2-8. This version uses least significant bit ordering as used in the GIF file format."},

		// control commands
		"push":   {f: e.push, In: "Data", Out: "Data, Data", Desc: "Duplicates the value on the top of the stack"},
		"pop":    {f: e.pop, In: "Data", Out: "", Desc: "Pops the value off the top of the stack (effectively discarding)"},
		"load":   {f: e.load, In: "Name", Out: "Value", Desc: "Pushes a named value from the dictinary onto the stack"},
		"save":   {f: e.save, In: "Value, Name", Out: "", Desc: "Pops a value from the stack and places it into the dictionary"},
		"swap":   {f: e.swap, In: "Val1, Val2", Out: "Val2, Val1", Desc: "Swaps the two values at the top of the stack"},
		"append": {f: e.append, In: "Val1, Val2", Out: "Appended", Desc: "Appends the value on the top of the stack to the previous value on the stack"},
		"slice":  {f: e.slice, In: "Data, Start, End", Out: "SliceOfData", Desc: "Slices the value on the stack, taking elements from start to end on the stack. Use -1 for values from the beginning or end. One example is /9/20/slice which takes elements 9 through 19, or /2/-1/slice which takes elements 2 through the end."},
		"len":    {f: e.len, In: "Data", Out: "Data, Length", Desc: "Pushes the length of the value on the stack in bytes onto the stack"},
		"left":   {f: e.left, In: "Data, Count", Out: "SliceOfData", Desc: "Takes the leftmost bytes of data"},
		"right":  {f: e.right, In: "Data, Count", Out: "SliceOfData", Desc: "Takes the rightmost bytes of data"},
		"snip":   {f: e.snip, In: "Data, Position", Out: "Data1, Data2", Desc: "Snips the data in half at the given position, resulting in two values on the stack"},
		"eq":     {f: e.eq, In: "Data1, Data2", Out: "", Desc: "Fails the command unless the two data elements are equal"},
		"neq":    {f: e.neq, In: "Data1, Data2", Out: "", Desc: "Fails the command unless the two data elements are not equal"},
		"call":   {f: e.call, In: "Name", Out: "(varies)", Desc: "Loads the named value from the dictionary and executes the commands contained there (formatted like normal - /md5/hex for example)"},

		// Crypto
		"aes-cfb":       {f: e.aes_cfb, In: "PlainData, IV, Key", Out: "CipherData", Desc: "Encrypts data using the given IV and 16-byte Key, placing the ciphertext back on the stack. Uses AES encryption and the CFB block mode."},
		"unaes-cfb":     {f: e.unaes_cfb, In: "CipherData, IV, Key", Out: "PlainData", Desc: "Decrypts data using the given IV and 16-byte Key, placing the plaintext back on the stack. Uses AES encryption and the CFB block mode."},
		"aes-ofb":       {f: e.aes_ofb, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 16-byte Key, placing the result back on the stack. Uses AES encryption and the OFB block mode."},
		"aes-ctr":       {f: e.aes_ctr, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 16-byte Key, placing the result back on the stack. Uses AES encryption and the CTR block mode."},
		"aes-blocksize": {f: e.aes_blocksize, In: "", Out: "16", Desc: "Pushes the AES block size on the stack"},

		"des-cfb":   {f: e.des_cfb, In: "PlainData, IV, Key", Out: "CipherData", Desc: "Encrypts data using the given IV and 8-byte Key, placing the ciphertext back on the stack. Uses DES encryption and the CFB block mode."},
		"undes-cfb": {f: e.undes_cfb, In: "CipherData, IV, Key", Out: "PlainData", Desc: "Decrypts data using the given IV and 8-byte Key, placing the plaintext back on the stack. Uses DES encryption and the CFB block mode."},
		"des-ofb":   {f: e.des_ofb, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 8-byte Key, placing the result back on the stack. Uses DES encryption and the OFB block mode."},
		"des-ctr":   {f: e.des_ctr, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 8-byte Key, placing the result back on the stack. Uses DES encryption and the CTR block mode."},

		"3des-cfb":       {f: e.tripledes_cfb, In: "PlainData, IV, Key", Out: "CipherData", Desc: "Encrypts data using the given IV and 24-byte Key, placing the ciphertext back on the stack. Uses Triple DES encryption and the CFB block mode."},
		"un3des-cfb":     {f: e.untripledes_cfb, In: "CipherData, IV, Key", Out: "PlainData", Desc: "Decrypts data using the given IV and 24-byte Key, placing the plaintext back on the stack. Uses Triple DES encryption and the CFB block mode."},
		"3des-ofb":       {f: e.tripledes_ofb, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 24-byte Key, placing the result back on the stack. Uses Triple DES encryption and the OFB block mode."},
		"3des-ctr":       {f: e.tripledes_ctr, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 24-byte Key, placing the result back on the stack. Uses Triple DES encryption and the CTR block mode."},
		"des-blocksize":  {f: e.des_blocksize, In: "", Out: "8", Desc: "Pushes the DES block size on the stack"},
		"3des-blocksize": {f: e.des_blocksize, In: "", Out: "8", Desc: "Pushes the Triple DES block size on the stack"},

		"blowfish-cfb":       {f: e.blowfish_cfb, In: "PlainData, IV, Key", Out: "CipherData", Desc: "Encrypts data using the given IV and 1 to 56-byte Key, placing the ciphertext back on the stack. Uses Blowfish encryption and the CFB block mode."},
		"unblowfish-cfb":     {f: e.unblowfish_cfb, In: "CipherData, IV, Key", Out: "PlainData", Desc: "Decrypts data using the given IV and 1 to 56-byte Key, placing the plaintext back on the stack. Uses Blowfish encryption and the CFB block mode."},
		"blowfish-ofb":       {f: e.blowfish_ofb, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses Blowfish encryption and the OFB block mode."},
		"blowfish-ctr":       {f: e.blowfish_ctr, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses Blowfish encryption and the CTR block mode."},
		"blowfish-blocksize": {f: e.blowfish_blocksize, In: "", Out: "8", Desc: "Pushes the blowfish block size on the stack"},

		"blowfish-salt-cfb":   {f: e.blowfish_salt_cfb, In: "PlainData, IV, Key, Salt", Out: "CipherData", Desc: "Encrypts data using the given IV and 1 to 56-byte Key, placing the ciphertext back on the stack. Uses Blowfish encryption and the CFB block mode."},
		"unblowfish-salt-cfb": {f: e.unblowfish_salt_cfb, In: "CipherData, IV, Key, Salt", Out: "PlainData", Desc: "Decrypts data using the given IV and 1 to 56-byte Key, placing the plaintext back on the stack. Uses Blowfish encryption and the CFB block mode."},
		"blowfish-salt-ofb":   {f: e.blowfish_salt_ofb, In: "Data, IV, Key, Salt", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses Blowfish encryption and the OFB block mode."},
		"blowfish-salt-ctr":   {f: e.blowfish_salt_ctr, In: "Data, IV, Key, Salt", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 1 to 56-byte Key, placing the result back on the stack. Uses Blowfish encryption and the CTR block mode."},

		"twofish-cfb":       {f: e.twofish_cfb, In: "PlainData, IV, Key", Out: "CipherData", Desc: "Encrypts data using the given IV and 16, 24, or 32-byte Key, placing the ciphertext back on the stack. Uses Twofish encryption and the CFB block mode."},
		"untwofish-cfb":     {f: e.untwofish_cfb, In: "CipherData, IV, Key", Out: "PlainData", Desc: "Decrypts data using the given IV and 16, 24, or 32-byte Key, placing the plaintext back on the stack. Uses Twofish encryption and the CFB block mode."},
		"twofish-ofb":       {f: e.twofish_ofb, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 16, 24, or 32-byte Key, placing the result back on the stack. Uses Twofish encryption and the OFB block mode."},
		"twofish-ctr":       {f: e.twofish_ctr, In: "Data, IV, Key", Out: "RData", Desc: "Encrypts or decrypts data using the given IV and 16, 24, or 32-byte Key, placing the result back on the stack. Uses Twofish encryption and the CTR block mode."},
		"twofish-blocksize": {f: e.twofish_blocksize, In: "", Out: "16", Desc: "Pushes the twofish block size on the stack"},
	}
}
