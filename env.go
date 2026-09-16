package env

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// env.Unmarshaler is the interface for any type that should unmarshal a textual
// representation of itself from an environment variable
type Unmarshaler interface {
	UnmarshalEnv(string) error
}

// Loads environment variables into a pointer of a stuct, to easily load multiple values.
//
// Reads from `env:` struct tags if available, or tries converting
// the struct field name to snake case. Keep in mind the conversion is a little janky
// and I wouldn't risk it.
func Load(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("%s needs to be a pointer", rv.Type().Name())
	}
	if rv.IsNil() {
		return fmt.Errorf("recieved a nil pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("%s is not a struct", rv.Type().Name())
	}
	_, err := loadStruct(rv)
	return err
}

func loadStruct(rv reflect.Value) (set bool, err error) {
	for field, val := range rv.Fields() {
		if !val.CanSet() {
			continue
		}
		fieldSet, err := loadField(field, val)
		if err != nil {
			return false, fmt.Errorf("loading field %s: %w", field.Name, err)
		}
		set = set || fieldSet
	}
	return set, nil
}

func loadField(field reflect.StructField, val reflect.Value) (set bool, err error) {
	t := strings.Split(field.Tag.Get("env"), ",")[0]
	if t == "-" {
		return false, nil
	}
	if t == "" {
		t = stringToEnvVar(field.Name)
	}
	env := os.Getenv(t)

	envunmarshaler := reflect.TypeFor[Unmarshaler]()
	if val.Addr().Type().Implements(envunmarshaler) {
		set, err = loadUnmarshalEnv(field, val, env)
		if err != nil {
			return false, fmt.Errorf("loading env.Unmarshaler %s: %w", field.Name, err)
		}
		return set, nil
	}
	textunmarshaler := reflect.TypeFor[encoding.TextUnmarshaler]()
	if val.Addr().Type().Implements(textunmarshaler) {
		set, err = loadUnmarshalText(field, val, env)
		if err != nil {
			return false, fmt.Errorf("loading encoding.TextUnmarshaler %s: %w", field.Name, err)
		}
		return set, nil
	}

	switch val.Kind() {
	case reflect.Struct:
		set, err = loadStruct(val)
		if err != nil {
			return false, fmt.Errorf("loading struct %s: %w", field.Name, err)
		}
	case reflect.Pointer:
		set, err = loadPtr(field, val)
		if err != nil {
			return false, fmt.Errorf("loading pointer %s: %w", field.Name, err)
		}
	case reflect.String:
		if env == "" {
			return false, nil
		}
		val.SetString(env)
		set = true
	case reflect.Int:
		if env == "" {
			return false, nil
		}
		i, err := strconv.ParseInt(env, 10, 0)
		if err != nil {
			return false, fmt.Errorf("loading \"%s\" as int: %w", env, err)
		}
		val.SetInt(i)
		set = true
	case reflect.Bool:
		if env == "" {
			return false, nil
		}
		b, err := strconv.ParseBool(env)
		if err != nil {
			return false, fmt.Errorf("loading \"%s\" as bool: %w", env, err)
		}
		val.SetBool(b)
		set = true
	case reflect.Slice:
		if env == "" {
			return false, nil
		}
		set, err = loadSlice(field, val, env)
		if err != nil {
			return false, fmt.Errorf("loading \"%s\" as slice: %w", env, err)
		}
	default:
		if field.Type.Name() == "Duration" {
			d, err := time.ParseDuration(env)
			if err != nil {
				return false, fmt.Errorf("loading \"%s\" as time.Duration: %w", env, err)
			}
			val.Set(reflect.ValueOf(d))
			return true, nil
		}
		return false, nil
	}
	return set, nil
}

func loadPtr(field reflect.StructField, val reflect.Value) (set bool, err error) {
	newval := reflect.New(field.Type.Elem()).Elem()
	if newval.Kind() == reflect.Struct {
		set, err = loadStruct(newval)
	} else {
		set, err = loadField(field, newval)
	}
	if err != nil {
		return false, err
	}
	if set {
		val.Set(newval.Addr())
	}
	return set, nil
}

func loadSlice(field reflect.StructField, val reflect.Value, env string) (set bool, err error) {
	t := field.Tag.Get("env")
	delim := ""
	idx := strings.Index(t, "delim=")
	if idx > -1 {
		delim = string(t[idx+6])
	}
	if delim == "" {
		delim = string(os.PathListSeparator)
	}
	sl := strings.Split(env, delim)
	rs := reflect.MakeSlice(field.Type, len(sl), cap(sl))
	switch field.Type.Elem().Kind() {
	case reflect.String:
		for i, str := range sl {
			rs.Index(i).SetString(str)
		}
		val.Set(rs)
		set = true
	case reflect.Int:
		for i, str := range sl {
			conv, err := strconv.ParseInt(str, 10, 0)
			if err != nil {
				return false, err
			}
			rs.Index(i).SetInt(conv)
		}
		val.Set(rs)
		set = true
	case reflect.Bool:
		for i, str := range sl {
			conv, err := strconv.ParseBool(str)
			if err != nil {
				return false, err
			}
			rs.Index(i).SetBool(conv)
		}
		val.Set(rs)
		set = true
	default:
		return false, fmt.Errorf("slice of unsupported type %s", field.Type.Elem().Kind())
	}
	return set, nil
}

func loadUnmarshalEnv(field reflect.StructField, val reflect.Value, env string) (set bool, err error) {
	unmarshal := val.Addr().MethodByName("UnmarshalEnv")
	envVal := reflect.ValueOf(env)
	ret := unmarshal.Call([]reflect.Value{envVal})
	if !ret[0].IsNil() {
		err, ok := ret[0].Interface().(error)
		if !ok {
			return false, fmt.Errorf("unknown error in UnmarshalEnv call on field \"%s\"", field.Name)
		}
		if err != nil {
			return false, fmt.Errorf("failed to UnmarshalEnv on value \"%s\": %w", env, err)
		}
	}
	return true, nil
}

func loadUnmarshalText(field reflect.StructField, val reflect.Value, env string) (set bool, err error) {
	unmarshal := val.Addr().MethodByName("UnmarshalText")
	byteVal := reflect.ValueOf([]byte(env))
	ret := unmarshal.Call([]reflect.Value{byteVal})
	if !ret[0].IsNil() {
		err, ok := ret[0].Interface().(error)
		if !ok {
			return false, fmt.Errorf("unknown error in UnmarshalText call on field \"%s\"", field.Name)
		}
		if err != nil {
			return false, fmt.Errorf("failed to UnmarshalText on value \"%s\": %w", env, err)
		}
	}
	return true, nil
}

// Takes a camel case string and converts it to upper case snake string
// which is what environment variables are usually named
//
// Stolen and augmented for my needs from Seth Bunce at
// https://groups.google.com/g/golang-nuts/c/MmerkVS9ke0?pli=1
func stringToEnvVar(st string) string {
	var parts []string
	start := 0
	upperCount := 0
	for end, r := range st {
		if end != 0 && unicode.IsLower(r) {
			if upperCount > 0 {
				if end == upperCount+1 {
					parts = append(parts, st[start:end-1])
				} else {
					parts = append(parts, st[start:end-upperCount])
					if upperCount > 1 {
						parts = append(parts, st[start+upperCount:end-1])
					}
				}
				start = end - 1
			}
			upperCount = 0
		}
		if end != 0 && unicode.IsUpper(r) {
			upperCount++
			// parts = append(parts, st[start:end])
			// start = end
		}
	}

	if start != len(st) {
		if upperCount > 0 {
			parts = append(parts, st[start:len(st)-upperCount])
			start = len(st) - upperCount
		}
		parts = append(parts, st[start:])
	}
	return strings.ToUpper(strings.Join(parts, "_"))
}
