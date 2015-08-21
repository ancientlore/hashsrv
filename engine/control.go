package engine

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

func (e *Engine) push() error {
	b := e.stack.Top()
	if b == nil {
		return errors.New("push - no value")
	}
	e.stack.Push(b)
	return nil
}

func (e *Engine) pop() error {
	b := e.stack.Pop()
	if b == nil {
		return errors.New("pop: Stack empty")
	}
	return nil
}

func (e *Engine) load() error {
	var str string
	str, err := e.stack.PopString()
	if err == nil {
		v := e.GetVariable(str)
		if v == nil {
			err = errors.New("Nil or no value called " + str)
		} else {
			e.stack.Push(v)
		}
	}
	return err
}

func (e *Engine) save() error {
	var str string
	str, err := e.stack.PopString()
	if err == nil {
		v := e.stack.Pop()
		if v == nil {
			err = errors.New("Cannot save - stack empty")
		} else {
			e.SetVariable(str, v)
		}
	}
	return err
}

func (e *Engine) swap() error {
	bs2 := e.stack.Pop()
	bs1 := e.stack.Pop()
	if bs1 == nil || bs2 == nil {
		return errors.New("swap: expected 2 values on stack")
	} else {
		e.stack.Push(bs2)
		e.stack.Push(bs1)
	}
	return nil
}

func (e *Engine) append() error {
	bs2 := e.stack.Pop()
	bs1 := e.stack.Pop()
	if bs1 == nil || bs2 == nil {
		return errors.New("append: expected 2 values on stack")
	} else {
		bs1 = append(bs1, bs2...)
		e.stack.Push(bs1)
	}
	return nil
}

func (e *Engine) slice() error {
	var err error
	var start, end int
	var d []byte
	end, err = e.stack.PopInt()
	if err == nil {
		start, err = e.stack.PopInt()
		if err == nil {
			d = e.stack.Pop()
			if d == nil {
				err = errors.New("slice: expected 3 values on the stack")
			} else {
				if start > len(d) || end > len(d) {
					err = errors.New("Out of range")
				} else if end < 0 && start < 0 {
					// do nothing
				} else if end >= 0 && start >= 0 {
					if start > end {
						err = errors.New("Start greater than end")
					} else {
						d = d[start:end]
					}
				} else if end < 0 {
					d = d[start:]
				} else if start < 0 {
					d = d[:end]
				}
			}
		}
	}
	if err == nil {
		e.stack.Push(d)
	}
	return err
}

func (e *Engine) left() error {
	// compose with other routines
	e.stack.Push([]byte("-1"))
	err := e.swap()
	if err == nil {
		err = e.slice()
	}
	return err
}

func (e *Engine) right() error {
	c, err := e.stack.PopInt()
	if err == nil {
		b := e.stack.Top()
		if b == nil {
			err = errors.New("right: expected 2 values on stack")
		} else {
			T := len(b)
			e.stack.Push([]byte(fmt.Sprintf("%d", T-c)))
			e.stack.Push([]byte(fmt.Sprintf("%d", T)))
			err = e.slice()
		}
	}
	return err
}

func (e *Engine) len() error {
	b := e.stack.Top()
	if b == nil {
		return errors.New("len: expected 1 value on stack")
	}
	e.stack.Push([]byte(fmt.Sprintf("%d", len(b))))
	return nil
}

func (e *Engine) snip() error {
	var err error
	var pos int
	var d []byte
	pos, err = e.stack.PopInt()
	if err == nil {
		d = e.stack.Pop()
		if d == nil {
			err = errors.New("snip: expected 2 values on the stack")
		} else {
			if pos >= len(d) {
				e.stack.Push(d)
				e.stack.Push(d[len(d):])
			} else if pos <= 0 {
				e.stack.Push(d[:0])
				e.stack.Push(d)
				// do nothing
			} else {
				e.stack.Push(d[:pos])
				e.stack.Push(d[pos:])
			}
		}
	}
	return err
}

func (e *Engine) eq() error {
	val1 := e.stack.Pop()
	val2 := e.stack.Pop()
	if val1 == nil || val2 == nil {
		return errors.New("eq: expected 2 values on the stack")
	}
	if bytes.Compare(val1, val2) != 0 {
		return errors.New("Values not equal")
	}
	return nil
}

func (e *Engine) neq() error {
	val1 := e.stack.Pop()
	val2 := e.stack.Pop()
	if val1 == nil || val2 == nil {
		return errors.New("neq: expected 2 values on the stack")
	}
	if bytes.Compare(val1, val2) == 0 {
		return errors.New("Values not expected to be equal")
	}
	return nil
}

func (e *Engine) call() error {
	nm, err := e.stack.PopString()
	if err != nil {
		return err
	}

	f, ok := e.values[nm]
	if !ok {
		return errors.New(fmt.Sprintf("call: cannot find %s", nm))
	}

	p := strings.TrimPrefix(string(f), "/")
	var commands []string
	if p != "" {
		commands = strings.Split(p, "/")
	} else {
		commands = make([]string, 0)
	}

	return e.exec(commands)
}
