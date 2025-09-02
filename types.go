/*
this file declare all types need to be for wx using
*/
package wx

import (
	"net/http"
	"reflect"
	"strings"
)

type authUtilsType struct {
	fn      map[reflect.Type]reflect.Value
	pkgPath string
	prefix  string
}

var authUtils = &authUtilsType{
	fn:      map[reflect.Type]reflect.Value{},
	pkgPath: reflect.TypeFor[authUtilsType]().PkgPath(),
	prefix:  strings.Split(reflect.TypeFor[OAuth2[any]]().Name(), "[")[0] + "[",
}

type OAuth2[T any] struct {
}

func (oauth *OAuth2[T]) Verify(fn func(ctx *httpContext) (*T, error)) {
	authUtils.fn[reflect.TypeFor[T]()] = reflect.ValueOf(fn)

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
