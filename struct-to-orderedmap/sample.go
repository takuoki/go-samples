package structtoorderedmap

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func StructToOrderedmap(value interface{}, publicOnly bool) (*orderedmap.OrderedMap, error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return nil, nil
		}
		return convertStruct(v.Elem(), publicOnly), nil
	case reflect.Struct:
		return convertStruct(v, publicOnly), nil
	default:
		return nil, fmt.Errorf("unsupported data type: %v", v.Kind())
	}
}

func convertStruct(v reflect.Value, publicOnly bool) *orderedmap.OrderedMap {
	if v.Kind() != reflect.Struct {
		panic("argument must be struct")
	}

	o := orderedmap.New()
	o.SetEscapeHTML(false)

	for i := 0; i < v.NumField(); i++ {
		if unsupportedType(v.Field(i)) {
			continue
		}
		if publicOnly && !v.Type().Field(i).IsExported() {
			continue
		}

		key := v.Type().Field(i).Tag.Get("json")
		if key == "" {
			key = v.Type().Field(i).Name
		} else {
			s := strings.Split(key, ",")
			key = s[0]
		}

		o.Set(key, convertValue(v.Field(i), publicOnly))
	}

	return o
}

func convertMap(v reflect.Value, publicOnly bool) *orderedmap.OrderedMap {
	if v.Kind() != reflect.Map {
		panic("argument must be map")
	}

	o := orderedmap.New()
	o.SetEscapeHTML(false)

	keys := v.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Kind() != reflect.String {
			panic(fmt.Sprintf("map keys only support string type: %v", keys[i].Kind()))
		}
		if keys[j].Kind() != reflect.String {
			panic(fmt.Sprintf("map keys only support string type: %v", keys[j].Kind()))
		}
		return keys[i].String() < keys[j].String()
	})

	for _, key := range keys {
		o.Set(key.String(), convertValue(v.MapIndex(key), publicOnly))
	}

	return o
}

func convertArraySlice(v reflect.Value, publicOnly bool) []interface{} {
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		panic("argument must be array or slice")
	}

	a := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		a[i] = convertValue(v.Index(i), publicOnly)
	}
	return a
}

func convertValue(v reflect.Value, publicOnly bool) interface{} {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.String:
		return v.String()
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return convertValue(v.Elem(), publicOnly)
	case reflect.Array:
		return convertArraySlice(v, publicOnly)
	case reflect.Slice:
		if v.IsNil() {
			return nil
		}
		return convertArraySlice(v, publicOnly)
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		return convertMap(v, publicOnly)
	case reflect.Struct:
		return convertStruct(v, publicOnly)
	default:
		panic(fmt.Sprintf("unsupported type: %v", v.Kind()))
	}
}

func unsupportedType(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String,
		reflect.Ptr, reflect.Interface,
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.Struct:
		return false
	default:
		return true
	}
}
