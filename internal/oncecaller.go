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

func findFirstField(typ reflect.Type, match func(t reflect.Type) bool, visited map[string]bool) ([]int, bool) {
	if _, ok := visited[typ.String()]; ok {
		return nil, false
	}
	visited[typ.String()] = true
	if match(typ) {
		return []int{}, true
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, false
	}
	for i := 0; i < typ.NumField(); i++ {
		if indexes, found := findFirstField(typ.Field(i).Type, match, visited); found {
			return append(typ.Field(i).Index, indexes...), true
		}
	}
	return nil, false
}

func FindFirstField(typ reflect.Type, match func(t reflect.Type) bool) ([]int, bool) {
	return findFirstField(typ, match, map[string]bool{})
}

type findFieldIndexByTypeResult struct {
	index []int
	found bool
}
type initFindFieldIndexByType struct {
	val  findFieldIndexByTypeResult
	once sync.Once
}

var cacheFindFieldIndexByType sync.Map

func FindFirstFieldOnce(typ reflect.Type, match func(t reflect.Type) (string, bool)) ([]int, bool) {
	key, ok := match(typ)
	if !ok {
		return nil, false
	}

	cacheKey := typ.PkgPath() + "/" + typ.String() + "::" + key
	actually, _ := cacheFindFieldIndexByType.LoadOrStore(cacheKey, &initFindFieldIndexByType{})
	item := actually.(*initFindFieldIndexByType)

	item.once.Do(func() {
		index, found := FindFirstField(typ, func(t reflect.Type) bool {
			_, ok := match(t)
			return ok
		})
		item.val = findFieldIndexByTypeResult{index, found}
	})

	return item.val.index, item.val.found
}

func FindFirstArg(method reflect.Method, match func(t reflect.Type) bool) (int, []int, bool) {
	for i := 0; i < method.Type.NumIn(); i++ {
		if index, found := FindFirstField(method.Type.In(i), match); found {
			return i, index, true
		}
	}
	return 0, nil, false
}
