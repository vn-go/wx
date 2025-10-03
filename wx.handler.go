package wx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type HttpContext struct {
	Req *http.Request
	Res http.ResponseWriter
	Uri string
}
type HttpError struct {
	Code HTTPErrorCode
	Data any
}
type swaggerInfoItem struct {
	requestContentType               string
	responseContentType              string
	requestBodyType                  reflect.Type
	responseBodyType                 reflect.Type
	controllerType                   reflect.Type
	HttpMethod                       string
	IsHasFileUpload                  bool
	listOfIndexFieldIsFormUploadFile []int
}

var handlerMapping = map[string]func(w http.ResponseWriter, r *http.Request){}
var swaggerInfo = map[string]swaggerInfoItem{}

// writeJSONResponse is a utility function to send a JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error writing JSON response: %v", err)
		}
	}
}
func parseJsonBody[TData any](w http.ResponseWriter, r *http.Request, uri string) (TData, error) {
	var data TData
	// 2. Deserialize request body into `data`
	// Limit the request body size (e.g., 1MB)
	r.Body = http.MaxBytesReader(w, r.Body, int64(currentServer.MaxBodySize))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Reject unknown fields in the JSON body
	if err := decoder.Decode(&data); err != nil {
		// Handle common JSON parsing errors
		log.Printf("JSON decode error for URI %s: %v", uri, err)
		var errMsg string
		if err == io.EOF {
			errMsg = "Request body must not be empty"
		} else if syntaxErr, ok := err.(*json.SyntaxError); ok {
			errMsg = fmt.Sprintf("Invalid JSON at offset %d", syntaxErr.Offset)
		} else if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			errMsg = fmt.Sprintf("Invalid type for field '%s'", unmarshalErr.Field)
		} else {
			errMsg = "Invalid request format"
		}
		var ret TData
		return ret, &HttpError{
			Code: http.StatusBadRequest,
			Data: map[string]string{"error": errMsg},
		}
		// Return 400 Bad Request

	}
	return data, nil

}

// parseFormBody parses 'application/x-www-form-urlencoded' body into struct TData.
// TData must use `form:"field_name"` tags for mapping.
func parseFormBody[TData any](w http.ResponseWriter, r *http.Request, uri string) (TData, error) {
	var data TData
	zeroValue := *new(TData) // Zero value of TData

	// 1. Check Content-Type (Ensure it's form data)
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return zeroValue, &HttpError{
			Code: http.StatusBadRequest,
			Data: map[string]string{"error": "Unsupported Content-Type. Expected application/x-www-form-urlencoded"},
		}
	}

	// 2. Limit the request body size (e.g., 1MB) and Parse form
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Call ParseForm to parse the body of the request
	if err := r.ParseForm(); err != nil {
		return zeroValue, &HttpError{
			Code: http.StatusBadRequest,
			Data: map[string]string{"error": "Invalid form data format or body too large"},
		}
	}

	// 3. Map form values (r.PostForm) to TData using Reflection

	// Get the value and type of TData
	dataValue := reflect.ValueOf(&data).Elem()
	dataType := dataValue.Type()

	// Ensure TData is a struct
	if dataType.Kind() != reflect.Struct {
		return zeroValue, &HttpError{
			Code: http.StatusInternalServerError,
			Data: map[string]string{"error": "TData must be a struct for form parsing"},
		}
	}

	// Iterate over struct fields
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		formTag := field.Tag.Get("form")

		// Use field name if 'form' tag is missing
		if formTag == "" {
			formTag = field.Name
		}

		formValues := r.PostForm[formTag]

		// Skip if no value found in form data
		if len(formValues) == 0 {
			continue
		}

		// Get the field value in the struct
		structField := dataValue.Field(i)

		// Try to set the value based on the field type (simple types only)
		value := formValues[0] // We only handle the first value for simplicity

		switch structField.Kind() {
		case reflect.String:
			structField.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
				structField.SetInt(intVal)
			} else {
				return zeroValue, &HttpError{
					Code: http.StatusBadRequest,
					Data: map[string]string{"error": fmt.Sprintf("Invalid integer format for field '%s'", formTag)},
				}
			}
		case reflect.Bool:
			// Handle common boolean values
			if boolVal, err := strconv.ParseBool(value); err == nil {
				structField.SetBool(boolVal)
			} else if value == "on" || value == "1" {
				structField.SetBool(true)
			} else if value == "off" || value == "0" {
				structField.SetBool(false)
			} else {
				return zeroValue, &HttpError{
					Code: http.StatusBadRequest,
					Data: map[string]string{"error": fmt.Sprintf("Invalid boolean format for field '%s'", formTag)},
				}
			}
		// Thêm các loại dữ liệu khác (float, slice, ...) nếu cần
		default:
			// Ignore unsupported types
		}
	}

	return data, nil
}

// parseMuitPartFormData parses 'multipart/form-data' body into struct TData,
// handling both regular fields and file uploads (*multipart.FileHeader).
// TData fields must use `form:"field_name"` tags for mapping.
func parseMuitPartFormData[TData any](r *http.Request, uri string, maxUploadMemory int64) (TData, error) {
	var data TData
	zeroValue := *new(TData) // Zero value of TData

	// 1. Check Content-Type (Ensure it's multipart/form-data)
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		return zeroValue, &HttpError{
			Code: http.StatusBadRequest,
			Data: map[string]string{"error": "Unsupported Content-Type. Expected multipart/form-data"},
		}
	}

	// 2. Parse Multipart Form
	// r.ParseMultipartForm limits the total size of the request body (including files) to maxUploadMemory
	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		// This handles errors like body too large, missing boundary, etc.
		return zeroValue, &HttpError{
			Code: http.StatusBadRequest,
			Data: map[string]string{"error": "Failed to parse multipart form (e.g., file too large, invalid format)"},
		}
	}

	// 3. Map form values (r.MultipartForm.Value) and Files (r.MultipartForm.File) to TData using Reflection

	dataValue := reflect.ValueOf(&data).Elem()
	dataType := dataValue.Type()

	if dataType.Kind() != reflect.Struct {
		return zeroValue, &HttpError{
			Code: http.StatusInternalServerError,
			Data: map[string]string{"error": "TData must be a struct for form parsing"},
		}
	}

	// Define the expected type for file headers
	fileHeaderType := reflect.TypeOf(&multipart.FileHeader{})

	// Iterate over struct fields
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		formTag := field.Tag.Get("form")

		// Use field name if 'form' tag is missing
		if formTag == "" {
			formTag = field.Name
		}

		structField := dataValue.Field(i)

		// A. Handle File Uploads (*multipart.FileHeader)
		if structField.Type() == fileHeaderType {
			if files := r.MultipartForm.File[formTag]; len(files) > 0 {
				// Assign the *multipart.FileHeader to the struct field
				structField.Set(reflect.ValueOf(files[0]))
			}
			continue // Done with this field, move to the next
		}

		// B. Handle Regular Form Values
		formValues := r.MultipartForm.Value[formTag]

		// Skip if no value found in form data
		if len(formValues) == 0 {
			continue
		}

		// Try to set the value based on the field type (simple types only)
		value := formValues[0] // We use the first value for simplicity

		switch structField.Kind() {
		case reflect.String:
			structField.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
				structField.SetInt(intVal)
			} else {
				return zeroValue, &HttpError{
					Code: http.StatusBadRequest,
					Data: map[string]string{"error": fmt.Sprintf("Invalid integer format for field '%s'", formTag)},
				}
			}
		case reflect.Bool:
			// Handle common boolean values
			if boolVal, err := strconv.ParseBool(value); err == nil {
				structField.SetBool(boolVal)
			} else if value == "on" || value == "1" || strings.EqualFold(value, "true") {
				structField.SetBool(true)
			} else if value == "off" || value == "0" || strings.EqualFold(value, "false") {
				structField.SetBool(false)
			} else {
				return zeroValue, &HttpError{
					Code: http.StatusBadRequest,
					Data: map[string]string{"error": fmt.Sprintf("Invalid boolean format for field '%s'", formTag)},
				}
			}
		// Add support for other types like float, slice, etc. here if needed
		default:
			// Optionally log or ignore unsupported types
		}
	}

	return data, nil
}

func newControllerInstance[TController any]() *TController {
	return reflect.New(reflect.TypeFor[TController]()).Interface().(*TController)
}

type EmptyBody struct{}
type HttpGet[TResponse any] struct {
	Data TResponse
}

func getAllFieldsIsFileUpload[TRquest any]() []int {
	ret := []int{}
	typ := reflect.TypeFor[TRquest]()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == reflect.TypeFor[multipart.FileHeader]() || typ.Field(i).Type == reflect.TypeFor[*multipart.FileHeader]() {
			ret = append(ret, i)
		}
	}
	if len(ret) > 0 {
		return ret
	} else {
		return nil
	}
}
func HandlerPost[TController any, TData any, TResponse any](uriHandler string, fn func(controllerInstance *TController, ctx *HttpContext, data TData) (TResponse, error)) {
	uri := fmt.Sprintf("%s/%s", reflect.TypeFor[TController]().Name(), uriHandler)
	if _, ok := handlerMapping[uri]; ok {
		panic(fmt.Sprintf("'%s is ready", uri))
	}
	allFieldsHasIsUpload := getAllFieldsIsFileUpload[TData]()
	//multipart/form-data
	requestContentType := "application/json"
	if allFieldsHasIsUpload != nil {
		requestContentType = "multipart/form-data"
	}
	swaggerItem := swaggerInfoItem{
		requestContentType:               requestContentType,
		responseContentType:              "application/json",
		requestBodyType:                  reflect.TypeFor[TData](),
		responseBodyType:                 reflect.TypeFor[TResponse](),
		controllerType:                   reflect.TypeFor[TController](),
		HttpMethod:                       "POST",
		IsHasFileUpload:                  allFieldsHasIsUpload != nil,
		listOfIndexFieldIsFormUploadFile: allFieldsHasIsUpload,
	}
	swaggerInfo[uri] = swaggerItem
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	hasBody := true
	if reflect.TypeFor[TData]() == reflect.TypeFor[EmptyBody]() {
		hasBody = false

	}
	controllerInstance := newControllerInstance[TController]()
	isHasFileUpload := allFieldsHasIsUpload != nil
	handlerMapping[uri] = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSONResponse(w, http.StatusBadGateway, map[string]string{"error": "Method is not allow"})
			return
		}
		ctx := &HttpContext{
			Req: r,
			Res: w,
			Uri: uri,
		}
		if isHasFileUpload {
			data, err := parseMuitPartFormData[TData](r, uri, int64(currentServer.MaxUploadSize))
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
				return
			}
			res, err := fn(controllerInstance, ctx, data)
			// 4. Handle response and errors
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
			} else {
				// Return successful response (200 OK)
				writeJSONResponse(w, http.StatusOK, res)
			}
		} else if hasBody {
			data, err := parseJsonBody[TData](w, r, uri)
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
				return
			}
			res, err := fn(controllerInstance, ctx, data)
			// 4. Handle response and errors
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
			} else {
				// Return successful response (200 OK)
				writeJSONResponse(w, http.StatusOK, res)
			}
		} else {
			var nilData TData
			res, err := fn(controllerInstance, ctx, nilData)
			// 4. Handle response and errors
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
			} else {
				// Return successful response (200 OK)
				writeJSONResponse(w, http.StatusOK, res)
			}
		}
	}

}
func HandlerGet[TController any, TResponse any](uriHandler string, fn func(controllerInstance *TController, ctx *HttpContext) (TResponse, error)) {
	uri := fmt.Sprintf("%s/%s", reflect.TypeFor[TController]().Name(), uriHandler)
	if _, ok := handlerMapping[uri]; ok {
		panic(fmt.Sprintf("'%s is ready", uri))
	}

	swaggerItem := swaggerInfoItem{
		requestContentType:  "application/json",
		responseContentType: "application/json",

		responseBodyType: reflect.TypeFor[TResponse](),
		controllerType:   reflect.TypeFor[TController](),
		HttpMethod:       "GET",
	}
	swaggerInfo[uri] = swaggerItem
	var nilrequestBodyType reflect.Type
	if swaggerItem.requestBodyType != nilrequestBodyType && swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}

	controllerInstance := newControllerInstance[TController]()
	handlerMapping[uri] = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeJSONResponse(w, http.StatusBadGateway, map[string]string{"error": "Method is not allow"})
			return
		}
		ctx := &HttpContext{
			Req: r,
			Res: w,
			Uri: uri,
		}
		res, err := fn(controllerInstance, ctx)
		if err != nil {
			if httpErr, ok := err.(*HttpError); ok {
				// Return a custom HTTP error (e.g., 401, 400)
				// The Data field is used as the JSON body
				log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
				writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
			} else {
				// Return an unhandled server error (500 Internal Server Error)
				log.Printf("Unhandled Server Error for %s: %v", uri, err)

				// Return a generic error message for security reasons
				writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		} else {
			// Return successful response (200 OK)
			writeJSONResponse(w, http.StatusOK, res)
		}
	}

}

func HandlerForm[TController any, TData any, TResponse any](uriHandler string, fn func(controllerInstance *TController, ctx *HttpContext, data TData) (TResponse, error)) {
	uri := fmt.Sprintf("%s/%s", reflect.TypeFor[TController]().Name(), uriHandler)
	if _, ok := handlerMapping[uri]; ok {
		panic(fmt.Sprintf("'%s is ready", uri))
	}

	swaggerItem := swaggerInfoItem{
		requestContentType:  "application/x-www-form-urlencoded",
		responseContentType: "application/json",
		requestBodyType:     reflect.TypeFor[TData](),
		responseBodyType:    reflect.TypeFor[TResponse](),
		controllerType:      reflect.TypeFor[TController](),
		HttpMethod:          "FORM_POST",
	}
	swaggerInfo[uri] = swaggerItem
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	hasBody := true
	if reflect.TypeFor[TData]() == reflect.TypeFor[EmptyBody]() {
		hasBody = false

	}
	controllerInstance := newControllerInstance[TController]()
	handlerMapping[uri] = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSONResponse(w, http.StatusBadGateway, map[string]string{"error": "Method is not allow"})
			return
		}
		ctx := &HttpContext{
			Req: r,
			Res: w,
			Uri: uri,
		}
		if hasBody {
			data, err := parseFormBody[TData](w, r, uri)
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
				return
			}
			res, err := fn(controllerInstance, ctx, data)
			// 4. Handle response and errors
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
			} else {
				// Return successful response (200 OK)
				writeJSONResponse(w, http.StatusOK, res)
			}
		} else {
			var nilData TData
			res, err := fn(controllerInstance, ctx, nilData)
			// 4. Handle response and errors
			if err != nil {
				if httpErr, ok := err.(*HttpError); ok {
					// Return a custom HTTP error (e.g., 401, 400)
					// The Data field is used as the JSON body
					log.Printf("Handled HTTP Error for %s: %v", uri, httpErr)
					writeJSONResponse(w, int(httpErr.Code), httpErr.Data)
				} else {
					// Return an unhandled server error (500 Internal Server Error)
					log.Printf("Unhandled Server Error for %s: %v", uri, err)

					// Return a generic error message for security reasons
					writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				}
			} else {
				// Return successful response (200 OK)
				writeJSONResponse(w, http.StatusOK, res)
			}
		}
	}

}
