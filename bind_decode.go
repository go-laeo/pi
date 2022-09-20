package pi

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

var (
	ErrInvalidP = errors.New("p must be *T")
)

var fmap = make(map[any]map[int]string)
var tmap = make(map[any]map[int]reflect.StructField)

func decode(m url.Values, v reflect.Value) error {
	if v.Type().Kind() == reflect.Interface || v.Type().Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ErrInvalidP
	}

	t := v.Type()

	if v, ok := fmap[t]; !ok {
		fmap[t] = make(map[int]string)
	} else {
		if len(v) != t.NumField() {
			fmap[t] = make(map[int]string)
		}
	}
	if _, ok := tmap[t]; !ok {
		tmap[t] = make(map[int]reflect.StructField)
	}

	for i := 0; i < v.NumField(); i++ {

		tf, ok := tmap[t][i]
		if !ok {
			tf = t.Field(i)
			tmap[t][i] = tf
		}

		vf := v.Field(i)

		if !tf.Anonymous && !vf.CanSet() {
			continue
		}

		if tf.Anonymous {
			// must be a struct
			err := decode(m, vf)
			if err != nil {
				return err
			}

			continue
		}

		query, ok := fmap[t][i]
		if !ok {
			query, ok = tf.Tag.Lookup("query")
			if !ok {
				query = tf.Name
			}
			fmap[t][i] = query
		}

		if query == "-" {
			continue
		}

		switch vf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(m.Get(query), 10, 64)
			if err != nil {
				return fmt.Errorf("parse %s to int: %w", query, err)
			}

			vf.SetInt(n)
		case reflect.String:
			vf.SetString(m.Get(query))
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(m.Get(query), 64)
			if err != nil {
				return fmt.Errorf("parse %s to float: %w", query, err)
			}
			vf.SetFloat(n)
		case reflect.Bool:
			b, _ := strconv.ParseBool(m.Get(query))
			vf.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(m.Get(query), 10, 64)
			if err != nil {
				return fmt.Errorf("parse %s to uint: %w", query, err)
			}
			vf.SetUint(n)
		case reflect.Struct:
			err := decode(m, vf)
			if err != nil {
				return fmt.Errorf("parse %s to struct: %w", query, err)
			}
		}
	}
	return nil
}
