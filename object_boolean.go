package jsc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type ObjectBoolean struct {
	ObjectBoolean map[string]bool
}

func NewObjectBoolean() ObjectBoolean {
	return ObjectBoolean{ObjectBoolean: map[string]bool{}}
}

func (o ObjectBoolean) Exists(key string) (exists bool) {
	_, ok := o.ObjectBoolean[key]
	return ok
}

func (o *ObjectBoolean) Set(key string, valueData bool) {
	if o.ObjectBoolean == nil {
		o.ObjectBoolean = map[string]bool{}
	}
	o.ObjectBoolean[key] = valueData
}

func (o ObjectBoolean) Get(key string) (value bool, ok bool) {
	if o.ObjectBoolean == nil {
		return false, false
	}
	value, ok = o.ObjectBoolean[key]
	return value, ok
}

func (o ObjectBoolean) Get2(key string) (value bool) {
	value, _ = o.Get(key)
	return value
}

func (o ObjectBoolean) Delete(key string) bool {
	if o.ObjectBoolean == nil {
		return false
	}
	removedValue, exists := o.ObjectBoolean[key]
	if !exists {
		return false
	}
	delete(o.ObjectBoolean, key)
	return removedValue
}

func (o ObjectBoolean) Value() (driver.Value, error) {
	if o.ObjectBoolean == nil {
		return nil, nil
	}
	return json.Marshal(o.ObjectBoolean)
}

func (o *ObjectBoolean) Scan(src interface{}) error {
	if src == nil {
		o.ObjectBoolean = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		o.ObjectBoolean = nil
		return errors.New("type assertion .([]byte) failed")
	}

	var mapBoolean map[string]bool
	err := json.Unmarshal(source, &mapBoolean)
	if err != nil {
		o.ObjectBoolean = nil
		return err
	}

	if mapBoolean == nil {
		o.ObjectBoolean = nil
		return nil
	}

	o.ObjectBoolean = mapBoolean

	return nil
}

func (o ObjectBoolean) MarshalJSON() ([]byte, error) {
	if o.ObjectBoolean != nil {
		return json.Marshal(o.ObjectBoolean)
	} else {
		return json.Marshal(nil)
	}
}

func (o *ObjectBoolean) UnmarshalJSON(data []byte) error {
	var x map[string]bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil {
		o.ObjectBoolean = nil
	}
	o.ObjectBoolean = x
	return nil
}

func MarshalObjectBoolean(o ObjectBoolean) graphql.Marshaler {
	if o.ObjectBoolean == nil {
		return graphql.Null
	}

	bytesResult, err := json.Marshal(o.ObjectBoolean)
	if err != nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bytesResult)
	})
}

func UnmarshalObjectBoolean(v interface{}) (ObjectBoolean, error) {
	switch t := v.(type) {
	case map[string]bool:
		return ObjectBoolean{ObjectBoolean: t}, nil
	case ObjectBoolean:
		return t, nil
	case map[string]interface{}:
		mapBoolean := map[string]bool{}
		for i, v := range t {
			var ok bool
			if mapBoolean[i], ok = v.(bool); !ok {
				return ObjectBoolean{}, fmt.Errorf("%T is not a boolean", v)
			}
		}
		return ObjectBoolean{ObjectBoolean: mapBoolean}, nil
	default:
		return ObjectBoolean{}, fmt.Errorf("%T is not an object of booleans", t)
	}
}
