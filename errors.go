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

//	func (f *errFactory) newUriParamParseError(paramName string, typeOfStruct reflect.Type) error {
//		return &UriParamParseError{
//			ParamName:    paramName,
//			TypeOfStruct: typeOfStruct,
//		}
//	}
func (e *UriParamParseError) Error() string {
	return fmt.Sprintf("%s was not found in %s", e.ParamName, e.TypeOfStruct.String())
}

type UriParamConvertError struct {
	ParamName        string
	ValueSetType     reflect.Type
	fielValueSetType reflect.Type
}

//	func (f *errFactory) newUriParamConvertError(paramName string, valueSetType reflect.Type, fielValueSetType reflect.Type) error {
//		return &UriParamConvertError{
//			ParamName:        paramName,
//			ValueSetType:     valueSetType,
//			fielValueSetType: fielValueSetType,
//		}
//	}
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

// func (f *errFactory) newParamMissMatchError(message string) error {
// 	return &ParamMissMatchError{
// 		Message: message,
// 	}
// }

func (e *ParamMissMatchError) Error() string {
	return e.Message

}

type ServiceInitError struct {
	Message string
}

// func (f *errFactory) newServiceInitError(message string) error {
// 	return &ServiceInitError{
// 		Message: message,
// 	}
// }

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

// func (f *errFactory) newMethodNotAllowError(message string) error {
// 	return &MethodNotAllowError{
// 		Message: message,
// 	}

// }

type NewMethodOfAuthNotFoundError struct {
	Message string
}

func (e *NewMethodOfAuthNotFoundError) Error() string {
	return e.Message
}

// func (f *errFactory) newMethodOfAuthNotFoundError(message string) error {
// 	return &NewMethodOfAuthNotFoundError{
// 		Message: message,
// 	}
// }

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

// func (f *errFactory) newUnSupportError(message string) error {
// 	return &UnSupportError{
// 		Message: message,
// 	}
// }

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

	ErrPaymentRequired = 402 // RFC 9110, 15.5.3

	ErrMethodNotAllowed  = 405 // RFC 9110, 15.5.6
	ErrNotAcceptable     = 406 // RFC 9110, 15.5.7
	ErrProxyAuthRequired = 407 // RFC 9110, 15.5.8
	ErrRequestTimeout    = 408 // RFC 9110, 15.5.9

	ErrGone                         = 410 // RFC 9110, 15.5.11
	ErrLengthRequired               = 411 // RFC 9110, 15.5.12
	ErrPreconditionFailed           = 412 // RFC 9110, 15.5.13
	ErrRequestEntityTooLarge        = 413 // RFC 9110, 15.5.14
	ErrRequestURITooLong            = 414 // RFC 9110, 15.5.15
	ErrUnsupportedMediaType         = 415 // RFC 9110, 15.5.16
	ErrRequestedRangeNotSatisfiable = 416 // RFC 9110, 15.5.17
	ErrExpectationFailed            = 417 // RFC 9110, 15.5.18
	ErrTeapot                       = 418 // RFC 9110, 15.5.19 (Unused)
	ErrMisdirectedRequest           = 421 // RFC 9110, 15.5.20

	ErrLocked                      = 423 // RFC 4918, 11.3
	ErrFailedDependency            = 424 // RFC 4918, 11.4
	ErrTooEarly                    = 425 // RFC 8470, 5.2.
	ErrUpgradeRequired             = 426 // RFC 9110, 15.5.22
	ErrPreconditionRequired        = 428 // RFC 6585, 3
	ErrTooManyRequests             = 429 // RFC 6585, 4
	ErrRequestHeaderFieldsTooLarge = 431 // RFC 6585, 5
	ErrUnavailableForLegalReasons  = 451 // RFC 7725, 3

	StatusContinue           = 100 // RFC 9110, 15.2.1
	StatusSwitchingProtocols = 101 // RFC 9110, 15.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1
	StatusEarlyHints         = 103 // RFC 8297

	StatusOK                   = 200 // RFC 9110, 15.3.1
	StatusCreated              = 201 // RFC 9110, 15.3.2
	StatusAccepted             = 202 // RFC 9110, 15.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 9110, 15.3.4
	StatusNoContent            = 204 // RFC 9110, 15.3.5
	StatusResetContent         = 205 // RFC 9110, 15.3.6
	StatusPartialContent       = 206 // RFC 9110, 15.3.7
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 9110, 15.4.1
	StatusMovedPermanently  = 301 // RFC 9110, 15.4.2
	StatusFound             = 302 // RFC 9110, 15.4.3
	StatusSeeOther          = 303 // RFC 9110, 15.4.4
	StatusNotModified       = 304 // RFC 9110, 15.4.5
	StatusUseProxy          = 305 // RFC 9110, 15.4.6
	_                       = 306 // RFC 9110, 15.4.7 (Unused)
	StatusTemporaryRedirect = 307 // RFC 9110, 15.4.8
	StatusPermanentRedirect = 308 // RFC 9110, 15.4.9

	ErrNotImplemented                = 501 // RFC 9110, 15.6.2
	ErrBadGateway                    = 502 // RFC 9110, 15.6.3
	ErrServiceUnavailable            = 503 // RFC 9110, 15.6.4
	ErrGatewayTimeout                = 504 // RFC 9110, 15.6.5
	ErrHTTPVersionNotSupported       = 505 // RFC 9110, 15.6.6
	ErrVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	ErrInsufficientStorage           = 507 // RFC 4918, 11.5
	ErrLoopDetected                  = 508 // RFC 5842, 7.2
	ErrNotExtended                   = 510 // RFC 2774, 7
	ErrNetworkAuthenticationRequired = 511 // RFC 6585, 6
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
