package pi

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

var (
	ErrInvalidP = errors.New("p must be *T")
)

var fmap = make(map[string]map[int]string)
var tmap = make(map[string]map[int]reflect.StructField)

func decode(m url.Values, v reflect.Value) error {
	if v.Type().Kind() == reflect.Interface || v.Type().Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ErrInvalidP
	}

	t := v.Type()
	name := t.Name()

	if _, ok := fmap[name]; !ok {
		fmap[name] = make(map[int]string)
	}
	if _, ok := tmap[name]; !ok {
		tmap[name] = make(map[int]reflect.StructField)
	}

	for i := 0; i < v.NumField(); i++ {

		tf, ok := tmap[name][i]
		if !ok {
			tf = t.Field(i)
			tmap[name][i] = tf
		}

		vf := v.Field(i)

		if !tf.Anonymous && !vf.CanSet() {
			// println("skip ", t.Field(i).Name)
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

		query, ok := fmap[name][i]
		if !ok {
			query, ok = tf.Tag.Lookup("query")
			if !ok {
				query = tf.Name
			}
			fmap[name][i] = query
		}

		if query == "-" {
			continue
		}

		switch vf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(m.Get(query), 10, 64)
			if err != nil {
				return err
			}

			vf.SetInt(n)
		case reflect.String:
			vf.SetString(m.Get(query))
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(m.Get(query), 64)
			if err != nil {
				return err
			}
			vf.SetFloat(n)
		case reflect.Bool:
			b, _ := strconv.ParseBool(m.Get(query))
			vf.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(m.Get(query), 10, 64)
			if err != nil {
				return err
			}
			vf.SetUint(n)
		case reflect.Struct:
			err := decode(m, vf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
