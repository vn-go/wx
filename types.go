/*
this file declare all types need to be for wx using
*/
package wx

import (
	"net/http"
	"reflect"
	"strings"
	"sync"
)

type authUtilsType struct {
	fn      map[reflect.Type]reflect.Value
	pkgPath string
	prefix  string
}

var authUtils = &authUtilsType{
	fn:      map[reflect.Type]reflect.Value{},
	pkgPath: reflect.TypeFor[authUtilsType]().PkgPath(),
	prefix:  strings.Split(reflect.TypeFor[Authenticate[any]]().Name(), "[")[0] + "[",
}

func (auth *authUtilsType) GetNewMethod(typ reflect.Type) (reflect.Value, bool) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if fieldData, ok := typ.FieldByName("Data"); ok {

		if ret, ok := auth.fn[fieldData.Type.Elem()]; ok {
			return ret, ok
		} else {
			return reflect.Value{}, false
		}
	}
	return reflect.Value{}, false
}

type Authenticate[T any] struct {
	Data *T
	Err  error
}
type initVerify struct {
	once sync.Once
}

var cacheVerify sync.Map

func (oauth *Authenticate[T]) Verify(fn func(ctx Handler) (*T, error)) {
	typ := reflect.TypeFor[T]()
	actually, _ := cacheVerify.LoadOrStore(typ, &initVerify{})
	actually.(*initVerify).once.Do(func() {
		authUtils.fn[reflect.TypeFor[T]()] = reflect.ValueOf(fn)
	})

}

/*
Any published method of a struct that has exactly one argument which is a Handler or an embedded HttpContext is called a HttpContext method.
*/
type httpContext struct {
	Req        *http.Request
	Res        http.ResponseWriter
	rootAbsUrl string
	schema     string
}
type Handler func() *httpContext
type Form[T any] struct {
	Data T
}

var pkgPath = reflect.TypeOf(Form[any]{}).PkgPath()
var formPrefixName = strings.Split(reflect.TypeOf(Form[any]{}).Name(), "[")[0] + "["

func isFormType(typ reflect.Type) bool {
	if typ.Kind() == reflect.Struct {
		if typ.PkgPath() == pkgPath && strings.HasPrefix(typ.Name(), formPrefixName) {
			return true
		}
	}
	return false
}
