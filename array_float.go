package jsc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type ArrayFloat struct {
	ArrayFloat []float64
}

func NewArrayFloat() ArrayFloat {
	return ArrayFloat{ArrayFloat: []float64{}}
}

func (a ArrayFloat) At(index int) float64 {
	l := len(a.ArrayFloat)
	if index > l-1 || index < 0 {
		return 0
	}
	return a.ArrayFloat[index]
}

func (a ArrayFloat) First() float64 {
	if len(a.ArrayFloat) == 0 {
		return 0
	}
	return a.ArrayFloat[0]
}

func (a ArrayFloat) Last() float64 {
	l := len(a.ArrayFloat)
	if l == 0 {
		return 0
	}
	return a.ArrayFloat[l-1]
}

func (a ArrayFloat) Length() int {
	return len(a.ArrayFloat)
}

func (a ArrayFloat) IndexOf(value float64) int {
	for i, v := range a.ArrayFloat {
		if v == value {
			return i
		}
	}
	return -1
}

func (a ArrayFloat) Filter(keep func(float64) bool) ArrayFloat {
	filtered := NewArrayFloat()
	for _, v := range a.ArrayFloat {
		if keep(v) {
			filtered.Push(v)
		}
	}
	return filtered
}

func (a ArrayFloat) DeleteAt(index int) float64 {
	l := len(a.ArrayFloat)
	if index > l-1 || index < 0 {
		return 0
	}
	removedValue := a.ArrayFloat[index]
	copy(a.ArrayFloat[index:], a.ArrayFloat[index+1:])
	a.ArrayFloat = a.ArrayFloat[:len(a.ArrayFloat)-1]
	return removedValue
}

func (a *ArrayFloat) Push(value float64) {
	if a.ArrayFloat == nil {
		a.ArrayFloat = []float64{}
	}
	a.ArrayFloat = append(a.ArrayFloat, value)
}

func (a *ArrayFloat) Pop() float64 {
	if len(a.ArrayFloat) == 0 {
		return 0
	}

	var lastValue float64
	lastValue, a.ArrayFloat = a.ArrayFloat[len(a.ArrayFloat)-1], a.ArrayFloat[:len(a.ArrayFloat)-1]

	return lastValue
}

func (a *ArrayFloat) Unshift(value float64) {
	if a.ArrayFloat == nil {
		a.ArrayFloat = []float64{}
	}
	a.ArrayFloat = append([]float64{value}, a.ArrayFloat...)
}

func (a *ArrayFloat) Shift() float64 {
	if len(a.ArrayFloat) == 0 {
		return 0
	}

	var firstValue float64
	firstValue, a.ArrayFloat = a.ArrayFloat[0], a.ArrayFloat[1:]

	return firstValue
}

func (a ArrayFloat) Value() (driver.Value, error) {
	if a.ArrayFloat == nil {
		return nil, nil
	}
	return json.Marshal(a.ArrayFloat)
}

func (a *ArrayFloat) Scan(src interface{}) error {
	if src == nil {
		a.ArrayFloat = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		a.ArrayFloat = nil
		return errors.New("type assertion .([]byte) failed")
	}

	var sliceFloat []float64
	err := json.Unmarshal(source, &sliceFloat)
	if err != nil {
		a.ArrayFloat = nil
		return err
	}

	if sliceFloat == nil {
		a.ArrayFloat = nil
		return nil
	}

	a.ArrayFloat = sliceFloat

	return nil
}

func (a ArrayFloat) MarshalJSON() ([]byte, error) {
	if a.ArrayFloat != nil {
		return json.Marshal(a.ArrayFloat)
	} else {
		return json.Marshal(nil)
	}
}

func (a *ArrayFloat) UnmarshalJSON(data []byte) error {
	var x []float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil {
		a.ArrayFloat = nil
	}
	a.ArrayFloat = x
	return nil
}

func MarshalArrayFloat(a ArrayFloat) graphql.Marshaler {
	if a.ArrayFloat == nil {
		return graphql.Null
	}

	bytesResult, err := json.Marshal(a.ArrayFloat)
	if err != nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bytesResult)
	})
}

func UnmarshalArrayFloat(v interface{}) (ArrayFloat, error) {
	switch t := v.(type) {
	case []float64:
		return ArrayFloat{ArrayFloat: t}, nil
	case ArrayFloat:
		return t, nil
	case []interface{}:
		sliceFloat := make([]float64, len(t))
		for i, v := range t {
			var ok bool
			if sliceFloat[i], ok = v.(float64); !ok {
				return ArrayFloat{}, fmt.Errorf("%T is not a float64", v)
			}
		}
		return ArrayFloat{ArrayFloat: sliceFloat}, nil
	default:
		return ArrayFloat{}, fmt.Errorf("%T is not an array of float64", t)
	}
}
