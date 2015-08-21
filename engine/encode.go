package engine

import (
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

func (e *Engine) hex() error {
	b := e.stack.Pop()
	enc := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(enc, b)
	e.stack.Push(enc)
	return nil
}

func (e *Engine) unhex() error {
	b := e.stack.Pop()
	dec := make([]byte, hex.DecodedLen(len(b)))
	n, err := hex.Decode(dec, b)
	if err == nil {
		e.stack.Push(dec[0:n])
	}
	return err
}

func (e *Engine) ascii85() error {
	b := e.stack.Pop()
	enc := make([]byte, ascii85.MaxEncodedLen(len(b)))
	sz := ascii85.Encode(enc, b)
	e.stack.Push(enc[0:sz])
	return nil
}

func (e *Engine) unascii85() error {
	b := e.stack.Pop()
	dec := make([]byte, len(b))
	sz, _, err := ascii85.Decode(dec, b, true)
	if err == nil {
		e.stack.Push(dec[0:sz])
	}
	return err
}

func (e *Engine) base32() error {
	b := e.stack.Pop()
	enc := make([]byte, base32.StdEncoding.EncodedLen(len(b)))
	base32.StdEncoding.Encode(enc, b)
	e.stack.Push(enc)
	return nil
}

func (e *Engine) unbase32() error {
	b := e.stack.Pop()
	dec := make([]byte, base32.StdEncoding.DecodedLen(len(b)))
	n, err := base32.StdEncoding.Decode(dec, b)
	if err == nil {
		e.stack.Push(dec[0:n])
	}
	return err
}

func (e *Engine) base32_hex() error {
	b := e.stack.Pop()
	enc := make([]byte, base32.HexEncoding.EncodedLen(len(b)))
	base32.HexEncoding.Encode(enc, b)
	e.stack.Push(enc)
	return nil
}

func (e *Engine) unbase32_hex() error {
	b := e.stack.Pop()
	dec := make([]byte, base32.HexEncoding.DecodedLen(len(b)))
	n, err := base32.HexEncoding.Decode(dec, b)
	if err == nil {
		e.stack.Push(dec[0:n])
	}
	return err
}

func (e *Engine) base64() error {
	b := e.stack.Pop()
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(enc, b)
	e.stack.Push(enc)
	return nil
}

func (e *Engine) unbase64() error {
	b := e.stack.Pop()
	dec := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	n, err := base64.StdEncoding.Decode(dec, b)
	if err == nil {
		e.stack.Push(dec[0:n])
	}
	return err
}

func (e *Engine) base64_url() error {
	b := e.stack.Pop()
	enc := make([]byte, base64.URLEncoding.EncodedLen(len(b)))
	base64.URLEncoding.Encode(enc, b)
	e.stack.Push(enc)
	return nil
}

func (e *Engine) unbase64_url() error {
	b := e.stack.Pop()
	dec := make([]byte, base64.URLEncoding.DecodedLen(len(b)))
	n, err := base64.URLEncoding.Decode(dec, b)
	if err == nil {
		e.stack.Push(dec[0:n])
	}
	return err
}
