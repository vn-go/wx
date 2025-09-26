package wx

import (
	"reflect"
	"sync"
)

type poolValue struct{}

type initPoolNew struct {
	once sync.Once
}

var (
	initPoolNewCache sync.Map
)
var mapTypePool map[reflect.Type]*sync.Pool = make(map[reflect.Type]*sync.Pool)

// Ngưỡng tối thiểu để dùng pool (bytes)
const minPoolSize = 128

func (p *poolValue) Get(typ reflect.Type) reflect.Value {
	t := typ
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Lấy size của type
	size := int(t.Size())

	// Nếu nhỏ hơn ngưỡng thì bỏ qua pool (tạo mới trực tiếp)
	if size < minPoolSize {
		return reflect.New(typ)
	}

	// Lớn hơn ngưỡng → dùng sync.Pool
	a, _ := initPoolNewCache.LoadOrStore(t, &initPoolNew{})
	i := a.(*initPoolNew)

	i.once.Do(func() {
		mapTypePool[t] = &sync.Pool{
			New: func() any {
				return reflect.New(typ).Interface()
			},
		}

	})

	return reflect.ValueOf(mapTypePool[t].Get())
}

func (p *poolValue) Put(val reflect.Value) {
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Reset value
	val.Elem().Set(reflect.Zero(typ))

	// Kiểm tra size
	size := int(typ.Size())
	if size < minPoolSize {
		// nhỏ thì không pool lại
		return
	}

	if i, ok := mapTypePool[typ]; ok {
		i.Put(val.Interface())
	}
}

var poolValues = &poolValue{}
var PoolValues = poolValues
