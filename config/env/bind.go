package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func bindEnv(env interface{}) error {
	if env == nil {
		return fmt.Errorf("env is nil")
	}
	rv := reflect.ValueOf(env)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("env must be a non-nil pointer")
	}
	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("env must point to a struct")
	}

	b := &binder{}
	b.walk(rv.Elem())

	if b.hasErrors() {
		return b.err()
	}
	return nil
}

type fieldErr struct {
	field string
	err   error
}

type binder struct {
	errs []fieldErr
}

func (b *binder) walk(rv reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		sf := rt.Field(i)

		if fv.Kind() == reflect.Struct {
			b.walk(fv)
			continue
		}

		envKey := sf.Tag.Get("env")
		if envKey == "" {
			continue
		}

		val, ok := resolveValue(envKey, sf.Tag.Get("default"))
		if !ok {
			continue
		}

		if err := setScalarValue(fv, val, sf); err != nil {
			b.errs = append(b.errs, fieldErr{field: sf.Name, err: err})
		}
	}
}

func resolveValue(envKey, def string) (string, bool) {
	if v, ok := os.LookupEnv(envKey); ok && v != "" {
		return v, true
	}
	if def != "" {
		return def, true
	}
	return "", false
}

func setScalarValue(fv reflect.Value, raw string, sf reflect.StructField) error {
	if !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", sf.Name)
	}
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(raw)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("parsing bool for %s: %w", sf.Name, err)
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("parsing int for %s: %w", sf.Name, err)
		}
		fv.SetInt(int64(i))
	default:
		return fmt.Errorf("unsupported kind %s for field %s", fv.Kind(), sf.Name)
	}
	return nil
}

func (b *binder) hasErrors() bool {
	return len(b.errs) > 0
}

func (b *binder) err() error {
	var sb strings.Builder
	sb.WriteString("binding errors:")
	for _, fe := range b.errs {
		sb.WriteString(" ")
		sb.WriteString(fe.field)
		sb.WriteString(": ")
		sb.WriteString(fe.err.Error())
		sb.WriteString(";")
	}
	return fmt.Errorf("%s", sb.String())
}
