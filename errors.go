/*
this file decalre custom errors of wx
*/
package wx

/*
Package errors provides custom error types used throughout the application.
This file defines detailed error structs to distinguish and handle errors
related to URI parameter parsing, value conversion, and service initialization.

Author: [Your Name or Team]
Created: 2025-08-16
License: MIT License (or your project's license)

Examples of errors in this package:
- UriParamParseError: error when a URI parameter is not found.
- UriParamConvertError: error when converting parameter value types.
- BadRequestError: invalid request error.
- ParamMissMatchError: parameter mismatch error.
- ServiceInitError: error during service initialization.
*/

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type errFactory struct {
}
type UriParamParseError struct {
	ParamName    string
	TypeOfStruct reflect.Type
}

func (f *errFactory) newUriParamParseError(paramName string, typeOfStruct reflect.Type) error {
	return &UriParamParseError{
		ParamName:    paramName,
		TypeOfStruct: typeOfStruct,
	}
}
func (e *UriParamParseError) Error() string {
	return fmt.Sprintf("%s was not found in %s", e.ParamName, e.TypeOfStruct.String())
}

type UriParamConvertError struct {
	ParamName        string
	ValueSetType     reflect.Type
	fielValueSetType reflect.Type
}

func (f *errFactory) newUriParamConvertError(paramName string, valueSetType reflect.Type, fielValueSetType reflect.Type) error {
	return &UriParamConvertError{
		ParamName:        paramName,
		ValueSetType:     valueSetType,
		fielValueSetType: fielValueSetType,
	}
}
func (e *UriParamConvertError) Error() string {
	return fmt.Sprintf("error converting from %s to %s", e.ValueSetType.String(), e.fielValueSetType.String())
}

type BadRequestError struct {
	Message string
}

func (f *errFactory) newBadRequestError(message string) error {
	return &BadRequestError{
		Message: message,
	}
}

func (e *BadRequestError) Error() string {
	return e.Message

}

type ParamMissMatchError struct {
	Message string
}

func (f *errFactory) newParamMissMatchError(message string) error {
	return &ParamMissMatchError{
		Message: message,
	}
}

func (e *ParamMissMatchError) Error() string {
	return e.Message

}

type ServiceInitError struct {
	Message string
}

func (f *errFactory) newServiceInitError(message string) error {
	return &ServiceInitError{
		Message: message,
	}
}

func (e *ServiceInitError) Error() string {
	return e.Message

}

type RequireError struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}

func (err *RequireError) Error() string {
	bff, ex := json.Marshal(err)
	if ex != nil {
		return err.Message
	}
	return string(bff)
}
func (f *errFactory) NewRequireError(fields []string, message string) error {
	return &RequireError{
		Fields:  fields,
		Message: message,
	}
}

type BodyParseError struct {
	Message    string
	InnerError error
}

func (f *errFactory) newBodyParseError(message string, InnerError error) error {
	return &BodyParseError{
		Message:    message,
		InnerError: InnerError,
	}
}

func (e *BodyParseError) Error() string {
	return e.Message

}

type FileParseError struct {
	Message    string
	InnerError error
}

func (e *FileParseError) Error() string {
	return e.Message
}
func (f *errFactory) newFileParseError(message string, InnerError error) error {
	return &FileParseError{
		Message:    message,
		InnerError: InnerError,
	}
}

type MethodNotAllowError struct {
	Message string
}

func (e *MethodNotAllowError) Error() string {
	return e.Message
}
func (f *errFactory) newMethodNotAllowError(message string) error {
	return &MethodNotAllowError{
		Message: message,
	}

}

type NewMethodOfAuthNotFoundError struct {
	Message string
}

func (e *NewMethodOfAuthNotFoundError) Error() string {
	return e.Message
}
func (f *errFactory) newMethodOfAuthNotFoundError(message string) error {
	return &NewMethodOfAuthNotFoundError{
		Message: message,
	}
}

type RegexUriNotMatchError struct {
	Message string
}

func (e *RegexUriNotMatchError) Error() string {
	return e.Message
}
func (f *errFactory) newRegexUriNotMatchError(message string) error {
	return &RegexUriNotMatchError{
		Message: message,
	}

}

type UnSupportError struct {
	Message string `json:"message"`
}

func (e *UnSupportError) Error() string {
	return e.Message
}
func (f *errFactory) newUnSupportError(message string) error {
	return &UnSupportError{
		Message: message,
	}
}

type ServerError struct {
	// error form system, nere show at front-end
	Err     error
	Message string
}

func (e *ServerError) Error() string {
	return e.Message
}
func (f *errFactory) NewServerError(msg string, err error) error {
	return &ServerError{
		Err:     err,
		Message: msg,
	}
}

/*
	{
	  "error": "unsupported_media_type",
	  "message": "Content-Type application/json is not supported. Please use multipart/form-data."
	}
*/
type unacceptableContentErrorData struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
type UnacceptableContentError struct {
	Data unacceptableContentErrorData
}

func (e *UnacceptableContentError) Error() string {
	bff, err := json.Marshal(e.Data)
	if err != nil {
		return Errors.NewServerError("Internal server error", err).Error()
	}
	return string(bff)

}
func (f *errFactory) newUnacceptableContentError(code, message string) error {
	return &UnacceptableContentError{
		Data: unacceptableContentErrorData{
			Error:   code,
			Message: message,
		},
	}
}

type ForbiddenError struct {
	data any
}

func (e *ForbiddenError) Error() string {
	bff, err := json.Marshal(e.data)
	if err != nil {
		return "Forbiden"
	} else {
		return string(bff)
	}

}
func (f *errFactory) NewForbidenError(data any) error {
	return &ForbiddenError{
		data: data,
	}
}

type UnauthorizedError struct {
}

func (e *UnauthorizedError) Error() string {
	return "Unauthorized"
}
func (f *errFactory) NewUnauthorizedError() error {
	return &UnauthorizedError{}
}

type HTTPErrorCode int

const (
	ErrOK                  HTTPErrorCode = 200
	ErrBadRequest          HTTPErrorCode = 400
	ErrUnauthorized        HTTPErrorCode = 401
	ErrForbidden           HTTPErrorCode = 403
	ErrNotFound            HTTPErrorCode = 404
	ErrConflict            HTTPErrorCode = 409
	ErrUnprocessableEntity HTTPErrorCode = 422
	ErrInternalServerError HTTPErrorCode = 500
	// có thể bổ sung thêm...
)

// String trả về chuỗi mô tả tương ứng với error code
func (c HTTPErrorCode) String() string {
	switch c {
	case ErrOK:
		return "OK"
	case ErrBadRequest:
		return "Bad Request"
	case ErrUnauthorized:
		return "Unauthorized"
	case ErrForbidden:
		return "Forbidden"
	case ErrNotFound:
		return "Not Found"
	case ErrConflict:
		return "Conflict"
	case ErrUnprocessableEntity:
		return "Unprocessable Entity"
	case ErrInternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Error"
	}
}

type HttpError struct {
	Code HTTPErrorCode
	Data any
}

func (e *HttpError) Error() string {
	v, err := json.Marshal(e.Data)
	if err != nil {
		return e.Code.String()
	}
	return string(v)
}
func (f *errFactory) NewHttpError(Code HTTPErrorCode, data any) error {
	return &HttpError{
		Code: Code,
		Data: data,
	}
}

var Errors = &errFactory{}
