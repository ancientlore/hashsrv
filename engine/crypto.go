package engine

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"fmt"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/twofish"
)

/*
	TODO:
	[X] crypto/aes
	[X] crypto/des
	[ ] crypto/dsa?
	[ ] crypto/ecdsa?
	[ ] crypto/elliptic?
	[ ] crypto/rc4
	[ ] crypto/rsa?
	[X] blowfish
	[X] twofish

	block modes:
	[X] CFB
	[ ] CBC - requires padding
	[X] CTR
	[X] OFB

	https://code.google.com/p/go/source/browse/?repo=crypto
	https://godoc.org/code.google.com/p/go.crypto/twofish
	http://golang.org/pkg/crypto/
*/

func (e *Engine) cfb(cipherBlock func(key []byte) (cipher.Block, error)) error {
	key := e.stack.Pop()
	iv := e.stack.Pop()
	plaintext := e.stack.Pop()
	if key == nil || iv == nil || plaintext == nil {
		return errors.New("Expected data, IV, and key on the stack")
	}
	block, err := cipherBlock(key)
	if err != nil {
		return err
	}
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)
	e.stack.Push(ciphertext)
	return nil
}

func (e *Engine) uncfb(cipherBlock func(key []byte) (cipher.Block, error)) error {
	key := e.stack.Pop()
	iv := e.stack.Pop()
	ciphertext := e.stack.Pop()
	if key == nil || iv == nil || ciphertext == nil {
		return errors.New("Expected data, IV, and key on the stack")
	}
	block, err := cipherBlock(key)
	if err != nil {
		return err
	}
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)
	e.stack.Push(plaintext)
	return nil
}

func (e *Engine) ofb(cipherBlock func(key []byte) (cipher.Block, error)) error {
	key := e.stack.Pop()
	iv := e.stack.Pop()
	text1 := e.stack.Pop()
	if key == nil || iv == nil || text1 == nil {
		return errors.New("Expected data, IV, and key on the stack")
	}
	block, err := cipherBlock(key)
	if err != nil {
		return err
	}
	text2 := make([]byte, len(text1))
	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(text2, text1)
	e.stack.Push(text2)
	return nil
}

func (e *Engine) ctr(cipherBlock func(key []byte) (cipher.Block, error)) error {
	key := e.stack.Pop()
	iv := e.stack.Pop()
	text1 := e.stack.Pop()
	if key == nil || iv == nil || text1 == nil {
		return errors.New("Expected data, IV, and key on the stack")
	}
	block, err := cipherBlock(key)
	if err != nil {
		return err
	}
	text2 := make([]byte, len(text1))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(text2, text1)
	e.stack.Push(text2)
	return nil
}

// CFB

func (e *Engine) aes_cfb() error {
	return e.cfb(aes.NewCipher)
}

func (e *Engine) unaes_cfb() error {
	return e.uncfb(aes.NewCipher)
}

func (e *Engine) des_cfb() error {
	return e.cfb(des.NewCipher)
}

func (e *Engine) undes_cfb() error {
	return e.uncfb(des.NewCipher)
}

func (e *Engine) tripledes_cfb() error {
	return e.cfb(des.NewTripleDESCipher)
}

func (e *Engine) untripledes_cfb() error {
	return e.uncfb(des.NewTripleDESCipher)
}

func (e *Engine) blowfish_cfb() error {
	return e.cfb(func(key []byte) (cipher.Block, error) { return blowfish.NewCipher(key) })
}

func (e *Engine) unblowfish_cfb() error {
	return e.uncfb(func(key []byte) (cipher.Block, error) { return blowfish.NewCipher(key) })
}

func (e *Engine) blowfish_salt_cfb() error {
	salt := e.stack.Pop()
	return e.cfb(func(key []byte) (cipher.Block, error) { return blowfish.NewSaltedCipher(key, salt) })
}

func (e *Engine) unblowfish_salt_cfb() error {
	salt := e.stack.Pop()
	return e.uncfb(func(key []byte) (cipher.Block, error) { return blowfish.NewSaltedCipher(key, salt) })
}

func (e *Engine) twofish_cfb() error {
	return e.cfb(func(key []byte) (cipher.Block, error) { return twofish.NewCipher(key) })
}

func (e *Engine) untwofish_cfb() error {
	return e.uncfb(func(key []byte) (cipher.Block, error) { return twofish.NewCipher(key) })
}

// OFB

func (e *Engine) aes_ofb() error {
	return e.ofb(aes.NewCipher)
}

func (e *Engine) des_ofb() error {
	return e.ofb(des.NewCipher)
}

func (e *Engine) tripledes_ofb() error {
	return e.ofb(des.NewTripleDESCipher)
}

func (e *Engine) blowfish_ofb() error {
	return e.ofb(func(key []byte) (cipher.Block, error) { return blowfish.NewCipher(key) })
}

func (e *Engine) blowfish_salt_ofb() error {
	salt := e.stack.Pop()
	return e.ofb(func(key []byte) (cipher.Block, error) { return blowfish.NewSaltedCipher(key, salt) })
}

func (e *Engine) twofish_ofb() error {
	return e.ofb(func(key []byte) (cipher.Block, error) { return twofish.NewCipher(key) })
}

// CTR

func (e *Engine) aes_ctr() error {
	return e.ctr(aes.NewCipher)
}

func (e *Engine) des_ctr() error {
	return e.ctr(des.NewCipher)
}

func (e *Engine) tripledes_ctr() error {
	return e.ctr(des.NewTripleDESCipher)
}

func (e *Engine) blowfish_ctr() error {
	return e.ctr(func(key []byte) (cipher.Block, error) { return blowfish.NewCipher(key) })
}

func (e *Engine) blowfish_salt_ctr() error {
	salt := e.stack.Pop()
	return e.ctr(func(key []byte) (cipher.Block, error) { return blowfish.NewSaltedCipher(key, salt) })
}

func (e *Engine) twofish_ctr() error {
	return e.ctr(func(key []byte) (cipher.Block, error) { return twofish.NewCipher(key) })
}

// Block size

func (e *Engine) aes_blocksize() error {
	e.stack.Push([]byte(fmt.Sprintf("%d", aes.BlockSize)))
	return nil
}

func (e *Engine) des_blocksize() error {
	e.stack.Push([]byte(fmt.Sprintf("%d", des.BlockSize)))
	return nil
}

func (e *Engine) blowfish_blocksize() error {
	e.stack.Push([]byte(fmt.Sprintf("%d", blowfish.BlockSize)))
	return nil
}

func (e *Engine) twofish_blocksize() error {
	e.stack.Push([]byte(fmt.Sprintf("%d", twofish.BlockSize)))
	return nil
}
