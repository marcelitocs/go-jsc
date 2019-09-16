package jsc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type ArrayInt struct {
	ArrayInt []int
}

func NewArrayInt() ArrayInt {
	return ArrayInt{ArrayInt: []int{}}
}

func (a ArrayInt) At(index int) int {
	l := len(a.ArrayInt)
	if index > l-1 || index < 0 {
		return 0
	}
	return a.ArrayInt[index]
}

func (a ArrayInt) First() int {
	if len(a.ArrayInt) == 0 {
		return 0
	}
	return a.ArrayInt[0]
}

func (a ArrayInt) Last() int {
	l := len(a.ArrayInt)
	if l == 0 {
		return 0
	}
	return a.ArrayInt[l-1]
}

func (a ArrayInt) Length() int {
	return len(a.ArrayInt)
}

func (a ArrayInt) IndexOf(value int) int {
	for i, v := range a.ArrayInt {
		if v == value {
			return i
		}
	}
	return -1
}

func (a ArrayInt) Filter(keep func(int) bool) ArrayInt {
	filtered := NewArrayInt()
	for _, v := range a.ArrayInt {
		if keep(v) {
			filtered.Push(v)
		}
	}
	return filtered
}

func (a ArrayInt) DeleteAt(index int) int {
	l := len(a.ArrayInt)
	if index > l-1 || index < 0 {
		return 0
	}
	removedValue := a.ArrayInt[index]
	copy(a.ArrayInt[index:], a.ArrayInt[index+1:])
	a.ArrayInt = a.ArrayInt[:len(a.ArrayInt)-1]
	return removedValue
}

func (a *ArrayInt) Push(value int) {
	if a.ArrayInt == nil {
		a.ArrayInt = []int{}
	}
	a.ArrayInt = append(a.ArrayInt, value)
}

func (a *ArrayInt) Pop() int {
	if len(a.ArrayInt) == 0 {
		return 0
	}

	var lastValue int
	lastValue, a.ArrayInt = a.ArrayInt[len(a.ArrayInt)-1], a.ArrayInt[:len(a.ArrayInt)-1]

	return lastValue
}

func (a *ArrayInt) Unshift(value int) {
	if a.ArrayInt == nil {
		a.ArrayInt = []int{}
	}
	a.ArrayInt = append([]int{value}, a.ArrayInt...)
}

func (a *ArrayInt) Shift() int {
	if len(a.ArrayInt) == 0 {
		return 0
	}

	var firstValue int
	firstValue, a.ArrayInt = a.ArrayInt[0], a.ArrayInt[1:]

	return firstValue
}

func (a ArrayInt) Value() (driver.Value, error) {
	if a.ArrayInt == nil {
		return nil, nil
	}
	return json.Marshal(a.ArrayInt)
}

func (a *ArrayInt) Scan(src interface{}) error {
	if src == nil {
		a.ArrayInt = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		a.ArrayInt = nil
		return errors.New("type assertion .([]byte) failed")
	}

	var sliceInt []int
	err := json.Unmarshal(source, &sliceInt)
	if err != nil {
		a.ArrayInt = nil
		return err
	}

	if sliceInt == nil {
		a.ArrayInt = nil
		return nil
	}

	a.ArrayInt = sliceInt

	return nil
}

func (a ArrayInt) MarshalJSON() ([]byte, error) {
	if a.ArrayInt != nil {
		return json.Marshal(a.ArrayInt)
	} else {
		return json.Marshal(nil)
	}
}

func (a *ArrayInt) UnmarshalJSON(data []byte) error {
	var x []int
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil {
		a.ArrayInt = nil
	}
	a.ArrayInt = x
	return nil
}

func MarshalArrayInt(a ArrayInt) graphql.Marshaler {
	if a.ArrayInt == nil {
		return graphql.Null
	}

	bytesResult, err := json.Marshal(a.ArrayInt)
	if err != nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bytesResult)
	})
}

func UnmarshalArrayInt(v interface{}) (ArrayInt, error) {
	switch t := v.(type) {
	case []int:
		return ArrayInt{ArrayInt: t}, nil
	case ArrayInt:
		return t, nil
	case []interface{}:
		sliceInt := make([]int, len(t))
		for i, v := range t {
			var ok bool
			if sliceInt[i], ok = v.(int); !ok {
				return ArrayInt{}, fmt.Errorf("%T is not a int", v)
			}
		}
		return ArrayInt{ArrayInt: sliceInt}, nil
	default:
		return ArrayInt{}, fmt.Errorf("%T is not an array of ints", t)
	}
}
