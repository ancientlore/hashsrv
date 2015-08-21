package engine

import (
	"container/list"
	"errors"
	"strconv"
)

type Stack list.List

func NewStack() *Stack {
	return (*Stack)(list.New())
}

func (stack *Stack) Len() int {
	list := (*list.List)(stack)
	return list.Len()
}

func (stack *Stack) Push(data []byte) {
	list := (*list.List)(stack)
	list.PushFront(data)
}

func (stack *Stack) Top() []byte {
	list := (*list.List)(stack)
	el := list.Front()
	if el == nil {
		return nil
	}
	val, ok := el.Value.([]byte)
	if !ok {
		panic("Why is it not a byte array?")
	}
	return val
}

func (stack *Stack) Pop() []byte {
	list := (*list.List)(stack)
	el := list.Front()
	if el == nil {
		return nil
	}
	val, ok := el.Value.([]byte)
	if !ok {
		panic("Why is it not a byte array?")
	}
	list.Remove(el)
	return val
}

func (stack *Stack) PopInt() (int, error) {
	val := stack.Pop()
	if val == nil {
		return -1, errors.New("Stack empty - integer required")
	}
	n, err := strconv.ParseInt(string(val), 10, 32)
	if err != nil {
		return -1, err
	}
	return int(n), nil
}

func (stack *Stack) PopString() (string, error) {
	val := stack.Pop()
	if val == nil {
		return "", errors.New("Stack empty - string required")
	}
	return string(val), nil
}

func (stack *Stack) ToArray() [][]byte {
	x := make([][]byte, 0, 4)
	list := (*list.List)(stack)
	for el := list.Back(); el != nil; el = el.Prev() {
		val, ok := el.Value.([]byte)
		if !ok {
			panic("Why is it not a byte array?")
		}
		x = append(x, val)
	}
	return x
}
