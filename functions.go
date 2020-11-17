package bcsv

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
)

func fileGetContents(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	return string(data), err
}

func filePutContents(filename string, data string, mode os.FileMode) error {
	return ioutil.WriteFile(filename, []byte(data), mode)
}

// reflectColumnsAsStringArrayForBcsvTags takes a struct and returns the computed columns names for each field with a "bcsv" tag
func reflectColumnsAsStringArrayForBcsvTags(t reflect.Type) []string {
	l := t.NumField()

	out := []string{}
	for x := 0; x < l; x++ {
		f := t.Field(x)
		h, ok := fieldHeaderNameForBcsvTag(f)
		if ok {
			out = append(out, h)
		}
	}

	return out
}

// reflectValuesAsStringArrayForBcsvTags takes a struct and returns the computed string values for each field with a "bcsv" tag
func reflectValuesAsStringArrayForBcsvTags(val reflect.Value, typ reflect.Type) []string {
	l := val.NumField()

	out := []string{}
	for x := 0; x < l; x++ {
		field := val.Field(x)
		bcsvTag := typ.Field(x).Tag.Get("bcsv")

		if bcsvTag == "-" {
			continue
		}

		out = append(out, reflectValueAsString(reflect.Value(field)))
	}

	return out
}

// reflectValueAsString returns the string representation of the reflect value
func reflectValueAsString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return fmt.Sprintf("%v", val.Int())
	case reflect.Float32:
		return strconv.FormatFloat(val.Float(), 'g', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return fmt.Sprintf("%v", val.Uint())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%+.3g", val.Complex())
	default:
		panic(fmt.Sprintf("Unsupported type %s", val.Kind()))
	}
}

// fieldHeaderNameForBcsvTag returns the header name for field with a "bcsv" tag
func fieldHeaderNameForBcsvTag(f reflect.StructField) (string, bool) {
	bcsvTag := f.Tag.Get("bcsv")

	if bcsvTag == "-" {
		return "", false
	}

	// If there is no tag set, use a default name
	if bcsvTag == "" {
		return f.Name, true
	}

	return bcsvTag, true
}
