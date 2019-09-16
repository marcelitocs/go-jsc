package jsc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type ObjectFloat struct {
	ObjectFloat map[string]float64
}

func NewObjectFloat() ObjectFloat {
	return ObjectFloat{ObjectFloat: map[string]float64{}}
}

func (o ObjectFloat) Exists(key string) (exists bool) {
	_, ok := o.ObjectFloat[key]
	return ok
}

func (o *ObjectFloat) Set(key string, valueData float64) {
	if o.ObjectFloat == nil {
		o.ObjectFloat = map[string]float64{}
	}
	o.ObjectFloat[key] = valueData
}

func (o ObjectFloat) Get(key string) (value float64, ok bool) {
	if o.ObjectFloat == nil {
		return 0, false
	}
	value, ok = o.ObjectFloat[key]
	return value, ok
}

func (o ObjectFloat) Get2(key string) (value float64) {
	value, _ = o.Get(key)
	return value
}

func (o ObjectFloat) Delete(key string) float64 {
	if o.ObjectFloat == nil {
		return 0
	}
	removedValue, exists := o.ObjectFloat[key]
	if !exists {
		return 0
	}
	delete(o.ObjectFloat, key)
	return removedValue
}

func (o ObjectFloat) Value() (driver.Value, error) {
	if o.ObjectFloat == nil {
		return nil, nil
	}
	return json.Marshal(o.ObjectFloat)
}

func (o *ObjectFloat) Scan(src interface{}) error {
	if src == nil {
		o.ObjectFloat = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		o.ObjectFloat = nil
		return errors.New("type assertion .([]byte) failed")
	}

	var mapFloat map[string]float64
	err := json.Unmarshal(source, &mapFloat)
	if err != nil {
		o.ObjectFloat = nil
		return err
	}

	if mapFloat == nil {
		o.ObjectFloat = nil
		return nil
	}

	o.ObjectFloat = mapFloat

	return nil
}

func (o ObjectFloat) MarshalJSON() ([]byte, error) {
	if o.ObjectFloat != nil {
		return json.Marshal(o.ObjectFloat)
	} else {
		return json.Marshal(nil)
	}
}

func (o *ObjectFloat) UnmarshalJSON(data []byte) error {
	var x map[string]float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil {
		o.ObjectFloat = nil
	}
	o.ObjectFloat = x
	return nil
}

func MarshalObjectFloat(o ObjectFloat) graphql.Marshaler {
	if o.ObjectFloat == nil {
		return graphql.Null
	}

	bytesResult, err := json.Marshal(o.ObjectFloat)
	if err != nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bytesResult)
	})
}

func UnmarshalObjectFloat(v interface{}) (ObjectFloat, error) {
	switch t := v.(type) {
	case map[string]float64:
		return ObjectFloat{ObjectFloat: t}, nil
	case ObjectFloat:
		return t, nil
	case map[string]interface{}:
		mapFloat := map[string]float64{}
		for i, v := range t {
			var ok bool
			if mapFloat[i], ok = v.(float64); !ok {
				return ObjectFloat{}, fmt.Errorf("%T is not a float64", v)
			}
		}
		return ObjectFloat{ObjectFloat: mapFloat}, nil
	default:
		return ObjectFloat{}, fmt.Errorf("%T is not an object of float64", t)
	}
}
