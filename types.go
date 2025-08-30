/*
this file declare all types need to be for wx using
*/
package wx

import "net/http"

/*
Any published method of a struct that has exactly one argument which is a Handler or an embedded HttpContext is called a HttpContext method.
*/
type HttpContext struct {
	Req *http.Request
	Res http.ResponseWriter
}
