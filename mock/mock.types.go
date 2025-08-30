package mock

import "reflect"

type HandlerInfo struct {
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
	IndexOfArg int
	/*
		HttpContext is a struct that has two important fields, Req and Res, with corresponding types *http.Request and http.ResponseWriter.
		This field represents the FieldIndex of Res
	*/
	ResFieldIndex []int
	/*
		HttpContext is a struct that has two important fields, Req and Res, with corresponding types *http.Request and http.ResponseWriter.
		This field represents the FieldIndex of Req
	*/
	ReqFieldIndex []int
	Method        reflect.Method
}
