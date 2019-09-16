package jsc

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type ObjectString struct {
	ObjectString map[string]string
}

func NewObjectString() ObjectString {
	return ObjectString{ObjectString: map[string]string{}}
}

func (o ObjectString) Exists(key string) (exists bool) {
	_, ok := o.ObjectString[key]
	return ok
}

func (o *ObjectString) Set(key string, valueData string) {
	if o.ObjectString == nil {
		o.ObjectString = map[string]string{}
	}
	o.ObjectString[key] = valueData
}

func (o ObjectString) Get(key string) (value string, ok bool) {
	if o.ObjectString == nil {
		return "", false
	}
	value, ok = o.ObjectString[key]
	return value, ok
}

func (o ObjectString) Get2(key string) (value string) {
	value, _ = o.Get(key)
	return value
}

func (o ObjectString) Delete(key string) string {
	if o.ObjectString == nil {
		return ""
	}
	removedValue, exists := o.ObjectString[key]
	if !exists {
		return ""
	}
	delete(o.ObjectString, key)
	return removedValue
}

func (o ObjectString) Value() (driver.Value, error) {
	if o.ObjectString == nil {
		return nil, nil
	}
	return json.Marshal(o.ObjectString)
}

func (o *ObjectString) Scan(src interface{}) error {
	if src == nil {
		o.ObjectString = nil
		return nil
	}

	source, ok := src.([]byte)
	if !ok {
		o.ObjectString = nil
		return errors.New("type assertion .([]byte) failed")
	}

	var mapString map[string]string
	err := json.Unmarshal(source, &mapString)
	if err != nil {
		o.ObjectString = nil
		return err
	}

	if mapString == nil {
		o.ObjectString = nil
		return nil
	}

	o.ObjectString = mapString

	return nil
}

func (o ObjectString) MarshalJSON() ([]byte, error) {
	if o.ObjectString != nil {
		return json.Marshal(o.ObjectString)
	} else {
		return json.Marshal(nil)
	}
}

func (o *ObjectString) UnmarshalJSON(data []byte) error {
	var x map[string]string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil {
		o.ObjectString = nil
	}
	o.ObjectString = x
	return nil
}

func MarshalObjectString(o ObjectString) graphql.Marshaler {
	if o.ObjectString == nil {
		return graphql.Null
	}

	bytesResult, err := json.Marshal(o.ObjectString)
	if err != nil {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write(bytesResult)
	})
}

func UnmarshalObjectString(v interface{}) (ObjectString, error) {
	switch t := v.(type) {
	case map[string]string:
		return ObjectString{ObjectString: t}, nil
	case ObjectString:
		return t, nil
	case map[string]interface{}:
		mapString := map[string]string{}
		for i, v := range t {
			var ok bool
			if mapString[i], ok = v.(string); !ok {
				return ObjectString{}, fmt.Errorf("%T is not a string", v)
			}
		}
		return ObjectString{ObjectString: mapString}, nil
	default:
		return ObjectString{}, fmt.Errorf("%T is not an object of strings", t)
	}
}
