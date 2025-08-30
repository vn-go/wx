/*
"This file declares all types used internally in the wx package
*/
package wx

import "reflect"

type queryParam struct {
	Name       string
	FieldIndex []int
}
type uriParam struct {
	Position   int
	Name       string
	IsSlug     bool
	FieldIndex []int
}
type handlerInfo struct {
	/*
		The argument index in the method is a Handler, or the struct has an embedded Handler field.
		Example:
			type MyStrcut struct {
			}
			func(m*MyStruct) Handler1( ctx *HttpContext) //<-- is hanler method, default Http is "POST"
			func(m*MyStruct) Handler1( ctx * struct {
				*HttpContext `route:"method:get"` //<-- specify method is "GET"
			}) //<-- is hanler method, default is post
		Note:
			Some DEV can convert publish method of struct to API handler by put an argument HttpContext any where in args of method,
			so need a field to keep position of Hanlder arg

	*/
	IndexOfArgIsHttpContext    int
	TypeOfArgIsHttpContext     reflect.Type
	TypeOfArgIsHttpContextElem reflect.Type
	/*
		HttpContext is a struct that has two important fields, Req and Res, with corresponding types *http.Request and http.ResponseWriter.
		This field represents the FieldIndex of Res
	*/
	ResFieldIndex []int
	/*
		HttpContext is a struct that has two important fields, Req and Res, with corresponding types *http.Request and http.ResponseWriter.
		This field represents the FieldIndex of Req
	*/
	ReqFieldIndex                    []int
	Method                           reflect.Method
	IsAbsUri                         bool
	Uri                              string
	IsQueryUri                       bool
	UriQuery                         string
	ControllerTypeElem               reflect.Type
	ControllerType                   reflect.Type
	UriParams                        []uriParam
	ListOfIndexFieldIsFormUploadFile []int
	TypeOfRequestBodyElem            reflect.Type
	TypeOfRequestBody                reflect.Type
	IndexOfArhIsAuthClaims           int
	IndexOfArgIsRequestBody          int
	IsFormPost                       bool
	FormPostTypeEle                  reflect.Type
	FormPostType                     reflect.Type
	HttpMethod                       string

	RouteTags      []string
	QueryParams    []queryParam
	RegexUri       string
	UriHandler     string
	IsRegexHandler bool
}
