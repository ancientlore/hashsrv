package engine

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"io"
	"io/ioutil"

	"github.com/golang/snappy"
)

func (e *Engine) snappy() error {
	data := snappy.Encode(nil, e.stack.Pop())
	e.stack.Push(data)
	return nil
}

func (e *Engine) unsnappy() error {
	data, err := snappy.Decode(nil, e.stack.Pop())
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) zlib() error {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(e.stack.Pop())
	w.Close()
	if err == nil {
		e.stack.Push(buf.Bytes())
	}
	return err
}

func (e *Engine) unzlib() error {
	buf := bytes.NewBuffer(e.stack.Pop())
	var r io.ReadCloser
	var err error
	var data []byte
	r, err = zlib.NewReader(buf)
	if err == nil {
		data, err = ioutil.ReadAll(r)
		r.Close()
	}
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) deflate() error {
	var buf bytes.Buffer
	var level int
	var err error
	level, err = e.stack.PopInt()
	if err == nil {
		var w *flate.Writer
		w, err = flate.NewWriter(&buf, level)
		if err == nil {
			_, err = w.Write(e.stack.Pop())
			w.Close()
		}
	}
	if err == nil {
		e.stack.Push(buf.Bytes())
	}
	return err
}

func (e *Engine) inflate() error {
	buf := bytes.NewBuffer(e.stack.Pop())
	r := flate.NewReader(buf)
	data, err := ioutil.ReadAll(r)
	r.Close()
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) gzip() error {
	var buf bytes.Buffer
	var level int
	var err error
	level, err = e.stack.PopInt()
	if err == nil {
		var w *gzip.Writer
		w, err = gzip.NewWriterLevel(&buf, level)
		if err == nil {
			_, err = w.Write(e.stack.Pop())
			w.Close()
		}
	}
	if err == nil {
		e.stack.Push(buf.Bytes())
	}
	return err
}

func (e *Engine) ungzip() error {
	buf := bytes.NewBuffer(e.stack.Pop())
	var r io.ReadCloser
	var data []byte
	var err error
	r, err = gzip.NewReader(buf)
	if err == nil {
		data, err = ioutil.ReadAll(r)
		r.Close()
	}
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) unbzip2() error {
	buf := bytes.NewBuffer(e.stack.Pop())
	r := bzip2.NewReader(buf)
	data, err := ioutil.ReadAll(r)
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) lzw_msb() error {
	var buf bytes.Buffer
	var litWidth int
	var err error
	litWidth, err = e.stack.PopInt()
	if err == nil {
		w := lzw.NewWriter(&buf, lzw.MSB, litWidth)
		_, err = w.Write(e.stack.Pop())
		w.Close()
	}
	if err == nil {
		e.stack.Push(buf.Bytes())
	}
	return err
}

func (e *Engine) lzw_lsb() error {
	var buf bytes.Buffer
	var litWidth int
	var err error
	litWidth, err = e.stack.PopInt()
	if err == nil {
		w := lzw.NewWriter(&buf, lzw.LSB, litWidth)
		_, err = w.Write(e.stack.Pop())
		w.Close()
	}
	if err == nil {
		e.stack.Push(buf.Bytes())
	}
	return err
}

func (e *Engine) unlzw_msb() error {
	var litWidth int
	var data []byte
	var err error
	litWidth, err = e.stack.PopInt()
	if err == nil {
		buf := bytes.NewBuffer(e.stack.Pop())
		r := lzw.NewReader(buf, lzw.MSB, litWidth)
		data, err = ioutil.ReadAll(r)
		r.Close()
	}
	if err == nil {
		e.stack.Push(data)
	}
	return err
}

func (e *Engine) unlzw_lsb() error {
	var litWidth int
	var data []byte
	var err error
	litWidth, err = e.stack.PopInt()
	if err == nil {
		buf := bytes.NewBuffer(e.stack.Pop())
		r := lzw.NewReader(buf, lzw.LSB, litWidth)
		data, err = ioutil.ReadAll(r)
		r.Close()
	}
	if err == nil {
		e.stack.Push(data)
	}
	return err
}
