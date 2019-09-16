package jsc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type ArrayString struct {
	ArrayString []string
}

func NewArrayString() ArrayString {
	return ArrayString{ArrayString: []string{}}
}

func (a ArrayString) At(index int) string {
	l := len(a.ArrayString)
	if index > l-1 || index < 0 {
		return ""
	}
	return a.ArrayString[index]
}

func (a ArrayString) First() string {
	if len(a.ArrayString) == 0 {
		return ""
	}
	return a.ArrayString[0]
}

func (a ArrayString) Last() string {
	l := len(a.ArrayString)
	if l == 0 {
		return ""
	}
	return a.ArrayString[l-1]
}

func (a ArrayString) Length() int {
	return len(a.ArrayString)
}

func (a ArrayString) IndexOf(value string) int {
	for i, v := range a.ArrayString {
		if v == value {
			return i
		}
	}
	return -1
}

func (a ArrayString) Filter(keep func(string) bool) ArrayString {
	filtered := NewArrayString()
	for _, v := range a.ArrayString {
		if keep(v) {
			filtered.Push(v)
		}
	}
	return filtered
}

func (a ArrayString) DeleteAt(index int) string {
	l := len(a.ArrayString)
	if index > l-1 || index < 0 {
		return ""
	}
	removedValue := a.ArrayString[index]
	copy(a.ArrayString[index:], a.ArrayString[index+1:])
	a.ArrayString = a.ArrayString[:len(a.ArrayString)-1]
	return removedValue
}

func (a *ArrayString) Push(value string) {
	if a.ArrayString == nil {
		a.ArrayString = []string{}
	}
	a.ArrayString = append(a.ArrayString, value)
}

func (a *ArrayString) Pop() string {
	if len(a.ArrayString) == 0 {
		return ""
	}

	var lastValue string
	lastValue, a.ArrayString = a.ArrayString[len(a.ArrayString)-1], a.ArrayString[:len(a.ArrayString)-1]

	return lastValue
}

func (a *ArrayString) Unshift(value string) {
	if a.ArrayString == nil {
		a.ArrayString = []string{}
	}
	a.ArrayString = append([]string{value}, a.ArrayString...)
}

func (a *ArrayString) Shift() string {
	if len(a.ArrayString) == 0 {
		return ""
	}

	var firstValue string
	firstValue, a.ArrayString = a.ArrayString[0], a.ArrayString[1:]

	return firstValue
}

func (a ArrayString) Value() (driver.Value, error) {
	if a.ArrayString == nil {
		return nil, nil
	}
	return json.Marshal(a.ArrayString)
}

func (a *ArrayString) Scan(src interface{}) error {
	if src == nil {
		a.ArrayString = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		a.ArrayString = nil
		return errors.New("type assertion .([]byte) failed")
	}

	var sliceString []string
	err := json.Unmarshal(source, &sliceString)
	if err != nil {
		a.ArrayString = nil
		return err
	}

	if sliceString == nil {
		a.ArrayString = nil
		return nil
	}

	a.ArrayString = sliceString

	return nil
}

func (a ArrayString) MarshalJSON() ([]byte, error) {
	if a.ArrayString != nil {
		return json.Marshal(a.ArrayString)
	} else {
		return json.Marshal(nil)
	}
}

func (a *ArrayString) UnmarshalJSON(data []byte) error {
	var x []string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil {
		a.ArrayString = nil
	}
	a.ArrayString = x
	return nil
}

func MarshalArrayString(a ArrayString) graphql.Marshaler {
	if a.ArrayString == nil {
		return graphql.Null
	}

	bytesResult, err := json.Marshal(a.ArrayString)
	if err != nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bytesResult)
	})
}

func UnmarshalArrayString(v interface{}) (ArrayString, error) {
	switch t := v.(type) {
	case []string:
		return ArrayString{ArrayString: t}, nil
	case ArrayString:
		return t, nil
	case []interface{}:
		sliceString := make([]string, len(t))
		for i, v := range t {
			var ok bool
			if sliceString[i], ok = v.(string); !ok {
				return ArrayString{}, fmt.Errorf("%T is not a string", v)
			}
		}
		return ArrayString{ArrayString: sliceString}, nil
	default:
		return ArrayString{}, fmt.Errorf("%T is not an array of strings", t)
	}
}
