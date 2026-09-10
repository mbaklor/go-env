package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

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
	return loadStruct(rv)
}

func loadStruct(rv reflect.Value) error {
	for field, val := range rv.Fields() {
		if !val.CanSet() {
			continue
		}
		if val.Kind() == reflect.Struct {
			err := loadStruct(val)
			if err != nil {
				return err
			}
			continue

		}
		err := loadField(field, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadField(field reflect.StructField, val reflect.Value) error {
	t := field.Tag.Get("env")
	if t == "" {
		t = stringToEnvVar(field.Name)
	}
	env := os.Getenv(t)
	switch val.Kind() {
	case reflect.String:
		val.SetString(env)
	case reflect.Int:
		i, err := strconv.ParseInt(env, 10, 0)
		if err != nil {
			return err
		}
		val.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(env)
		if err != nil {
			return err
		}
		val.SetBool(b)
	}
	return nil
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
