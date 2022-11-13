package config

import (
	"errors"
	"reflect"
)

// Merge two structs. Values in src will overwrite those in dst, provided
// that they do not contains the zero value for it's type.
func Merge(dst any, src any) error {
	dstv := reflect.ValueOf(dst)
	if dstv.Kind() != reflect.Pointer || dstv.Elem().Kind() != reflect.Struct {
		return errors.New("dst must be a pointer to a struct")
	}
	srcv := reflect.ValueOf(src)
	if srcv.Kind() == reflect.Pointer {
		if srcv.Elem().Kind() != reflect.Struct {
			return errors.New("src must either be a struct or a pointer to a struct")
		}
	} else {
		if srcv.Kind() != reflect.Struct {
			return errors.New("src must either be a struct or a pointer to a struct")
		}
	}

	dstv = reflect.ValueOf(dst).Elem()
	setFields(dstv, reflect.ValueOf(src))
	return nil
}

// setFields sets values from src into dst provided that they are not their
// zero value for it's type.
func setFields(dst reflect.Value, src reflect.Value) {
	if isEmpty(src) {
		return
	}
	if src.Kind() == reflect.Pointer {
		src = src.Elem()
	}
	for i := 0; i < src.NumField(); i++ {
		if isEmpty(src.Field(i)) && src.Field(i).Kind() != reflect.Bool {
			continue
		}
		if src.Field(i).Kind() == reflect.Pointer && src.Field(i).Elem().Kind() == reflect.Struct {
			setFields(dst.Field(i).Elem(), src.Field(i).Elem())
		} else if src.Field(i).Kind() == reflect.Struct {
			setFields(dst.Field(i), src.Field(i))
		} else {
			setValue(dst.Field(i), src.Field(i))
		}
	}
}

// setValue sets the value of src into dst.
func setValue(dst reflect.Value, src reflect.Value) {
	if dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type()).Elem())
		}
	} else if dst.Kind() == reflect.Struct {
		if dst.IsZero() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
	}
	dst.Set(src)
}

// isEmpty checks if the provided reflect.Value is either
// nil if a pointer, or zero value for all other cases.
func isEmpty(v reflect.Value) bool {
	if v.Kind() == reflect.Pointer {
		return v.IsNil()
	} else {
		return v.IsZero()
	}
}
