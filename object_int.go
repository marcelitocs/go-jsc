package jsc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type ObjectInt struct {
	ObjectInt map[string]int
}

func NewObjectInt() ObjectInt {
	return ObjectInt{ObjectInt: map[string]int{}}
}

func (o ObjectInt) Exists(key string) (exists bool) {
	_, ok := o.ObjectInt[key]
	return ok
}

func (o *ObjectInt) Set(key string, valueData int) {
	if o.ObjectInt == nil {
		o.ObjectInt = map[string]int{}
	}
	o.ObjectInt[key] = valueData
}

func (o ObjectInt) Get(key string) (value int, ok bool) {
	if o.ObjectInt == nil {
		return 0, false
	}
	value, ok = o.ObjectInt[key]
	return value, ok
}

func (o ObjectInt) Get2(key string) (value int) {
	value, _ = o.Get(key)
	return value
}

func (o ObjectInt) Delete(key string) int {
	if o.ObjectInt == nil {
		return 0
	}
	removedValue, exists := o.ObjectInt[key]
	if !exists {
		return 0
	}
	delete(o.ObjectInt, key)
	return removedValue
}

func (o ObjectInt) Value() (driver.Value, error) {
	if o.ObjectInt == nil {
		return nil, nil
	}
	return json.Marshal(o.ObjectInt)
}

func (o *ObjectInt) Scan(src interface{}) error {
	if src == nil {
		o.ObjectInt = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		o.ObjectInt = nil
		return errors.New("type assertion .([]byte) failed")
	}

	var mapInt map[string]int
	err := json.Unmarshal(source, &mapInt)
	if err != nil {
		o.ObjectInt = nil
		return err
	}

	if mapInt == nil {
		o.ObjectInt = nil
		return nil
	}

	o.ObjectInt = mapInt

	return nil
}

func (o ObjectInt) MarshalJSON() ([]byte, error) {
	if o.ObjectInt != nil {
		return json.Marshal(o.ObjectInt)
	} else {
		return json.Marshal(nil)
	}
}

func (o *ObjectInt) UnmarshalJSON(data []byte) error {
	var x map[string]int
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil {
		o.ObjectInt = nil
	}
	o.ObjectInt = x
	return nil
}

func MarshalObjectInt(o ObjectInt) graphql.Marshaler {
	if o.ObjectInt == nil {
		return graphql.Null
	}

	bytesResult, err := json.Marshal(o.ObjectInt)
	if err != nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bytesResult)
	})
}

func UnmarshalObjectInt(v interface{}) (ObjectInt, error) {
	switch t := v.(type) {
	case map[string]int:
		return ObjectInt{ObjectInt: t}, nil
	case ObjectInt:
		return t, nil
	case map[string]interface{}:
		mapInt := map[string]int{}
		for i, v := range t {
			var ok bool
			if mapInt[i], ok = v.(int); !ok {
				return ObjectInt{}, fmt.Errorf("%T is not a int", v)
			}
		}
		return ObjectInt{ObjectInt: mapInt}, nil
	default:
		return ObjectInt{}, fmt.Errorf("%T is not an object of ints", t)
	}
}
