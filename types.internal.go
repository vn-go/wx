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
	indexOfArgIsHttpContext    int
	typeOfArgIsHttpContext     reflect.Type
	typeOfArgIsHttpContextElem reflect.Type
	/*
		HttpContext is a struct that has two important fields, Req and Res, with corresponding types *http.Request and http.ResponseWriter.
		This field represents the Fieldindex of Res
	*/
	resFieldIndex []int
	/*
		HttpContext is a struct that has two important fields, Req and Res, with corresponding types *http.Request and http.ResponseWriter.
		This field represents the Fieldindex of Req
	*/
	reqFieldIndex                    []int
	method                           reflect.Method
	isNoOutPut                       bool
	isAbsUri                         bool
	uri                              string
	isQueryUri                       bool
	uriQuery                         string
	controllerTypeElem               reflect.Type
	controllerType                   reflect.Type
	conrollerNewMethod               *reflect.Method
	uriParams                        []uriParam
	listOfIndexFieldIsFormUploadFile []int
	typeOfRequestBodyElem            reflect.Type
	typeOfRequestBody                reflect.Type
	indexOfArhIsAuthClaims           int
	indexOfArgIsRequestBody          int
	isFormPost                       bool
	// formPostTypeEle                  reflect.Type
	// formPostType                     reflect.Type
	httpMethod                       string

	routeTags                          []string
	queryParams                        []queryParam
	regexUri                           string
	uriHandler                         string
	isRegexHandler                     bool
	fieldIndexOfResController          []int
	fieldIndexOfReqController          []int
	fiedIndexOfHttpContextInController []int
	hasHttpContextInController         bool
}
