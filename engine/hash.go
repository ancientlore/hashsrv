package engine

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"

	"golang.org/x/crypto/ripemd160"
)

func computeHash(h hash.Hash, data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("no data provided to hash")
	}
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func computeHmac(hf func() hash.Hash, key []byte, data []byte) ([]byte, error) {
	if key == nil {
		return nil, errors.New("no key provided for hmac")
	}
	/*
		Now assuming key is already hashed by the user - more flexible
		hashedKey, err := computeHash(hf(), key)
		if err != nil {
			return nil, err
		}
	*/
	return computeHash(hmac.New(hf, key), data)
}

func (e *Engine) md5() error {
	data, err := computeHash(md5.New(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) sha1() error {
	data, err := computeHash(sha1.New(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) sha224() error {
	data, err := computeHash(sha256.New224(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) sha256() error {
	data, err := computeHash(sha256.New(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) sha384() error {
	data, err := computeHash(sha512.New384(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) sha512() error {
	data, err := computeHash(sha512.New(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) ripemd160() error {
	data, err := computeHash(ripemd160.New(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) hmac_md5() error {
	k := e.stack.Pop()
	data, err := computeHmac(md5.New, k, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) hmac_sha1() error {
	k := e.stack.Pop()
	data, err := computeHmac(sha1.New, k, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) hmac_sha224() error {
	k := e.stack.Pop()
	data, err := computeHmac(sha256.New224, k, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) hmac_sha256() error {
	k := e.stack.Pop()
	data, err := computeHmac(sha256.New, k, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) hmac_sha384() error {
	k := e.stack.Pop()
	data, err := computeHmac(sha512.New384, k, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) hmac_sha512() error {
	k := e.stack.Pop()
	data, err := computeHmac(sha512.New, k, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) hmac_ripemd160() error {
	k := e.stack.Pop()
	data, err := computeHmac(ripemd160.New, k, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) adler32() error {
	data, err := computeHash(adler32.New(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) crc32() error {
	data, err := computeHash(crc32.NewIEEE(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) crc32_ieee() error {
	data, err := computeHash(crc32.New(crc32.MakeTable(crc32.IEEE)), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) crc32_castagnoli() error {
	data, err := computeHash(crc32.New(crc32.MakeTable(crc32.Castagnoli)), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) crc32_koopman() error {
	data, err := computeHash(crc32.New(crc32.MakeTable(crc32.Koopman)), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) crc64_iso() error {
	data, err := computeHash(crc64.New(crc64.MakeTable(crc64.ISO)), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) crc64_ecma() error {
	data, err := computeHash(crc64.New(crc64.MakeTable(crc64.ECMA)), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) fnv32() error {
	data, err := computeHash(fnv.New32(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) fnv32a() error {
	data, err := computeHash(fnv.New32a(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) fnv64() error {
	data, err := computeHash(fnv.New64(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) fnv64a() error {
	data, err := computeHash(fnv.New64a(), e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) rand() error {
	var data []byte
	var sz int
	var err error
	sz, err = e.stack.PopInt()
	if err == nil {
		data = make([]byte, sz)
		_, err = rand.Read(data)
	}
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) md5_len() error {
	e.stack.Push([]byte("16"))
	return nil
}

func (e *Engine) sha1_len() error {
	e.stack.Push([]byte("20"))
	return nil
}

func (e *Engine) sha224_len() error {
	e.stack.Push([]byte("28"))
	return nil
}

func (e *Engine) sha256_len() error {
	e.stack.Push([]byte("32"))
	return nil
}

func (e *Engine) sha384_len() error {
	e.stack.Push([]byte("48"))
	return nil
}

func (e *Engine) sha512_len() error {
	e.stack.Push([]byte("64"))
	return nil
}

func (e *Engine) ripemd160_len() error {
	e.stack.Push([]byte("20"))
	return nil
}
