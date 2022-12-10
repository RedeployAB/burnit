package config

import (
	"errors"
	"reflect"
)

// merge two structs. Values in src will overwrite those in dst, provided
// that they do not contains the zero value for it's type.
func merge(dst any, src any) error {
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
	if src.Kind() == reflect.Pointer {
		if src.IsNil() {
			return
		}
		src = src.Elem()
	}
	for i := 0; i < src.NumField(); i++ {
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
	}
	if src.IsZero() && src.Kind() != reflect.Bool {
		return
	}
	dst.Set(src)
}
