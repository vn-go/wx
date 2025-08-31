package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

type oneCallItem[TResult any] struct {
	Val *TResult
	Err error
}
type initOnceCall[Tkey any, TResult any] struct {
	Once sync.Once
	Item *oneCallItem[TResult]
}

var cacheOnceCall sync.Map

func OnceCall[Tkey any, TResult any](key Tkey, fn func() (*TResult, error)) (*TResult, error) {
	actual, _ := cacheOnceCall.LoadOrStore(key, &initOnceCall[Tkey, TResult]{})
	onceCall := actual.(*initOnceCall[Tkey, TResult])
	onceCall.Once.Do(func() {
		onceCall.Item = &oneCallItem[TResult]{}
		onceCall.Item.Val, onceCall.Item.Err = fn()
	})
	return onceCall.Item.Val, onceCall.Item.Err
}
func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func FirstElements[T any](input [][]T) []T {
	result := make([]T, 0, len(input))
	for _, inner := range input {
		if len(inner) > 0 {
			result = append(result, inner[0])
		}
	}
	return result
}
func IsNumber(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}
func ValToString(val interface{}) string {
	if val == nil {
		return ""
	}

	v := reflect.ValueOf(val)

	switch v.Kind() {
	case reflect.String:
		return v.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)

	case reflect.Bool:
		return strconv.FormatBool(v.Bool())

	// special-case []byte → string
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes())
		}
	}

	// fallback: fmt.Sprint
	return fmt.Sprint(val)
}
