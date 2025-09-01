/*
this file declare all types need to be for wx using
*/
package wx

import (
	"net/http"
	"reflect"
	"strings"
)

/*
Any published method of a struct that has exactly one argument which is a Handler or an embedded HttpContext is called a HttpContext method.
*/
type HttpContext struct {
	Req *http.Request
	Res http.ResponseWriter
}
type Handler func() *HttpContext
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
