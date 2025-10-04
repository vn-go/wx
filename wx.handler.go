package wx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

type IdentifierInfo[TIdentifier any] struct {
	Id any
}
type HttpContext[TIdentifier any] struct {
	RootRouter string
	Req        *http.Request
	Res        http.ResponseWriter
	Uri        string
	Identifier any
	//fix so need'nt clear when put to sync.pool
	rootAbsUrl string
	//fix so need'nt clear when put to sync.pool
	schema   string
	uriRegex string
	//parse all path in Req.Uri
	//example "/api/hellocontroller/test/hello4/movie/a001.mp4?filePath=loaad-dir"
	//for uriRegex is "^hellocontroller/test/hello4/([^/]+)/([^/]+)\\.mp4$"
	//and aslo know RootRouter = "/api"
	uriParams     map[string]string
	uriParamsList []string
	//parse all query to example "?a=1&b=2" queryParams={"a":1,"b":2}
	queryParams map[string]string
}

func (c *HttpContext[TIdentifier]) Reset() {
	c.Req = nil
	c.Res = nil
	c.Uri = ""
	c.Identifier = nil
	//c.isRegex = false
	c.queryParams = map[string]string{}
	c.uriParams = map[string]string{}
	c.uriParamsList = []string{}

}

// loadAllParams parses both the path parameters (uriParams)
// and query parameters (queryParams) from the request URL.
//
// Example:
//
//	Req.URL.Path  = "/api/hellocontroller/test/hello4/movie/a001.mp4"
//	RootRouter    = "/api"
//	uriRegex      = "^hellocontroller/test/hello4/([^/]+)/([^/]+)\\.mp4$"
//	→ uriParams = {"0": "movie", "1": "a001"}
//
//	Req.URL.RawQuery = "filePath=load-dir&lang=en"
//	→ queryParams = {"filePath": "load-dir", "lang": "en"}
func (c *HttpContext[TIdentifier]) loadAllParams() {
	if c.Req == nil {
		return
	}

	// --- 1. Parse query parameters ---
	queryValues := c.Req.URL.Query()
	if len(queryValues) > 0 {
		if c.queryParams == nil {
			c.queryParams = make(map[string]string, len(queryValues))
		}
		for key, vals := range queryValues {
			if len(vals) > 0 {
				c.queryParams[key] = vals[0]
			}
		}
	}

	// --- 2. Parse URI path parameters (if regex is available) ---
	if c.uriRegex == "" {
		return
	}

	reg, err := regexp.Compile(c.uriRegex)
	if err != nil {
		// Invalid regex – skip path parsing
		return
	}

	path := strings.TrimPrefix(c.Req.URL.Path, c.RootRouter)
	path = strings.TrimPrefix(path, "/")

	matches := reg.FindStringSubmatch(path)
	if len(matches) == 0 {
		return
	}

	if c.uriParams == nil {
		c.uriParams = make(map[string]string, len(matches)-1)
	}

	// Note: If the regex does not have named groups, we use numeric keys: "0", "1", etc.
	for i, val := range matches[1:] {
		if i >= 0 && i < len(c.uriParamsList) {
			c.uriParams[c.uriParamsList[i]] = val
		}

	}
}

func (h *HttpContext[TIdentifier]) GetIdentifier() *TIdentifier {
	if ret, ok := h.Identifier.(*TIdentifier); ok {
		return ret
	} else {
		return nil
	}
}

type ContextPool[T any] struct {
	pool sync.Pool
}

func newContextPool[T any]() *ContextPool[T] {
	return &ContextPool[T]{
		pool: sync.Pool{
			New: func() any {
				return new(HttpContext[T])
			},
		},
	}
}

func (p *ContextPool[T]) get() *HttpContext[T] {
	return p.pool.Get().(*HttpContext[T])
}

func (p *ContextPool[T]) put(ctx *HttpContext[T]) {
	p.pool.Put(ctx)
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
	isRequireAuth                    bool
	uriInfo                          inspectUriInfo
}

var handlerMapping = map[string]func(w http.ResponseWriter, r *http.Request){}
var swaggerInfo = map[string]swaggerInfoItem{}
var jsonBufferPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}
var jsonLib = jsoniter.ConfigFastest

func writeJSONResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	stream := jsoniter.NewStream(jsoniter.ConfigFastest, w, 512)
	defer jsonLib.ReturnStream(stream)
	stream.WriteVal(v)
	stream.Flush()
	// stream := jsonLib.BorrowStream(w)

	// stream.WriteVal(v)
	// stream.Flush() // ghi ra trực tiếp ResponseWriter
}

func parseJsonBody[TData any](w http.ResponseWriter, r *http.Request, uri string) (TData, error) {
	var data TData

	// 1️⃣ Giới hạn kích thước body
	r.Body = http.MaxBytesReader(w, r.Body, int64(currentServer.MaxBodySize))

	// 2️⃣ Đọc toàn bộ body vào buffer (nhỏ, tránh syscall nhiều)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var ret TData
		return ret, &HttpError{
			Code: http.StatusBadRequest,
			Data: map[string]string{"error": "Cannot read request body"},
		}
	}

	// 3️⃣ Parse JSON bằng json-iterator (nhanh hơn 2–3 lần stdlib)
	if err := jsonLib.Unmarshal(body, &data); err != nil {
		log.Printf("JSON decode error for URI %s: %v", uri, err)
		var errMsg string
		switch {
		case err == io.EOF:
			errMsg = "Request body must not be empty"

		default:
			errMsg = fmt.Sprintf("Invalid request format: %v", err)
		}

		var ret TData
		return ret, &HttpError{
			Code: http.StatusBadRequest,
			Data: map[string]string{"error": errMsg},
		}
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
		formTag := field.Tag.Get("json")

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
func parseMuitPartFormData[TData any](r *http.Request, maxUploadMemory int64) (TData, error) {
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

type initNewControllerInstance[TController any] struct {
	val  *TController
	err  error
	once sync.Once
}

var initNewControllerInstanceCache sync.Map

func newControllerInstance[TController any]() (*TController, error) {
	a, _ := initNewControllerInstanceCache.LoadOrStore(reflect.TypeFor[TController](), &initNewControllerInstance[TController]{})
	init := a.(*initNewControllerInstance[TController])
	init.once.Do(func() {
		retIntance := reflect.New(reflect.TypeFor[TController]()).Interface().(*TController)
		//find New function of *TTController
		for i := 0; i < reflect.TypeFor[*TController]().NumMethod(); i++ {
			if reflect.TypeFor[*TController]().Method(i).Name == "New" {
				out := reflect.TypeFor[*TController]().Method(i).Func.Call([]reflect.Value{reflect.ValueOf(retIntance)})
				if len(out) > 0 {
					finalOut := out[len(out)-1]
					if finalOut.Interface() != nil {
						if err, ok := finalOut.Interface().(error); ok {

							init.err = err
							return
						}
					}
				}
			}
		}
		init.val = retIntance

	})
	if init.err != nil {
		initNewControllerInstanceCache.Delete(reflect.TypeFor[TController]())
		return nil, init.err
	}
	return init.val, nil
}

// type HttpGet[TResponse any] struct {
// 	Data TResponse
// }

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

// check TTypeCheck is a certain type
// TTypeCheck is a certain struct not any
// check TTypeCheck is a struct (used as an identifier/context), but not 'any'
// Returns: true if TTypeCheck is a struct, false if TTypeCheck is 'any' or another basic type.
func isSpecificType[TTypeCheck any]() bool {
	t := reflect.TypeFor[TTypeCheck]()

	// Nếu kiểu là 'interface' (tức là 'any'), chúng ta coi là không có yêu cầu Auth cụ thể.
	if t.Kind() == reflect.Interface {
		return false // isRequieAuth[any] -> false
	}

	// Nếu kiểu là 'struct', chúng ta coi là có yêu cầu Auth cụ thể.
	if t.Kind() == reflect.Struct {
		return true // isRequieAuth[struct{...}] -> true
	}

	// Các kiểu khác (string, int, v.v.)
	return false
}

type cacheSetAuthKey struct {
	ControllerType reflect.Type
	IdentifierType reflect.Type
}
type OKUser struct {
}

var cacheSetAuth = map[any]any{}

func SetAuth[TController any, TIdentifier any](controller *TController, fn func(ctx *HttpContext[TIdentifier]) error) {
	cacheSetAuth[*controller] = fn
	//cacheSetAuth[*controller] = fn
}
func callAuth[TController any, TIdentifier any](c *TController, ctx *HttpContext[TIdentifier]) error {
	fn, ok := cacheSetAuth[*c]
	if !ok {
		return &HttpError{
			Code: http.StatusUnauthorized,
			Data: map[string]string{"error": "require login"},
		}
	}

	if fx, ok := fn.(func(ctx *HttpContext[TIdentifier]) error); ok {
		err := fx(ctx)
		if err != nil {
			return err
		}
		return nil
	} else {
		return &HttpError{
			Code: http.StatusUnauthorized,
			Data: map[string]string{"error": "require login"},
		}
	}

}

var mapControllerInstance = map[string]any{}

func HandlerPost[TController any, TIdentifier any, TData any, TResponse any](uriHandler string, fn func(controllerInstance *TController, ctx *HttpContext[TIdentifier], data TData) (TResponse, error)) {

	uri := fmt.Sprintf("%s/%s", strings.ToLower(reflect.TypeFor[TController]().Name()), uriHandler)
	if _, ok := handlerMapping[uri]; ok {
		panic(fmt.Sprintf("'%s is ready", uri))
	}
	allFieldsHasIsUpload := getAllFieldsIsFileUpload[TData]()
	//multipart/form-data
	requestContentType := "application/json"
	if allFieldsHasIsUpload != nil {
		requestContentType = "multipart/form-data"
	}
	uriInfo := inspectUri(uri)
	swaggerItem := swaggerInfoItem{
		requestContentType:               requestContentType,
		responseContentType:              "application/json",
		requestBodyType:                  reflect.TypeFor[TData](),
		responseBodyType:                 reflect.TypeFor[TResponse](),
		controllerType:                   reflect.TypeFor[TController](),
		HttpMethod:                       "POST",
		IsHasFileUpload:                  allFieldsHasIsUpload != nil,
		listOfIndexFieldIsFormUploadFile: allFieldsHasIsUpload,
		isRequireAuth:                    isSpecificType[TIdentifier](),
		uriInfo:                          uriInfo,
	}
	swaggerInfo[uriInfo.mainUri] = swaggerItem
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	requireAuth := isSpecificType[TIdentifier]()
	hasBody := isSpecificType[TData]()

	var initControllerError error
	controllerInstance, initControllerError := newControllerInstance[TController]()
	if initControllerError != nil {
		mapControllerInstance[uri] = controllerInstance
	}
	isHasFileUpload := allFieldsHasIsUpload != nil

	handlerMapping[uriInfo.mainUri] = func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {

				fmt.Println(string(debug.Stack()))
			}
		}()
		if initControllerError != nil {
			if currentServer.IsReleaseMode {
				defer log.Panicln(initControllerError)
				writeJSONResponse(w, http.StatusInternalServerError, "server error")
				return
			} else {
				panic(initControllerError)
			}

		}
		if r.Method != "POST" {
			writeJSONResponse(w, http.StatusBadGateway, map[string]string{"error": "Method is not allow"})
			return
		}
		ctxPool := newContextPool[TIdentifier]()
		ctx := ctxPool.get()
		ctx.uriParamsList = uriInfo.uriParams
		ctx.RootRouter = currentServer.BaseUrl
		ctx.Req = r
		ctx.Res = w
		ctx.uriRegex = uriInfo.uriRegex
		ctx.Uri = uriInfo.mainUri
		if uriInfo.isRegex {
			ctx.loadAllParams()
		}
		defer func() {
			ctx.Reset()
			ctxPool.put(ctx)
		}()

		if requireAuth {
			err := callAuth(controllerInstance, ctx)
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
		}
		if isHasFileUpload {
			data, err := parseMuitPartFormData[TData](r, int64(currentServer.MaxUploadSize))
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
func HandlerGet[TController any, TIdentifier any, TResponse any](uriHandler string, fn func(controllerInstance *TController, ctx *HttpContext[TIdentifier]) (TResponse, error)) {
	uri := fmt.Sprintf("%s/%s", strings.ToLower(reflect.TypeFor[TController]().Name()), uriHandler)
	if _, ok := handlerMapping[uri]; ok {
		panic(fmt.Sprintf("'%s is ready", uri))
	}
	uriInfo := inspectUri(uri)
	swaggerItem := swaggerInfoItem{
		requestContentType:  "application/json",
		responseContentType: "application/json",

		responseBodyType: reflect.TypeFor[TResponse](),
		controllerType:   reflect.TypeFor[TController](),
		HttpMethod:       "GET",
		isRequireAuth:    isSpecificType[TIdentifier](),
		uriInfo:          uriInfo,
	}
	swaggerInfo[uriInfo.mainUri] = swaggerItem
	var nilrequestBodyType reflect.Type
	if swaggerItem.requestBodyType != nilrequestBodyType && swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	var initControllerError error
	controllerInstance, initControllerError := newControllerInstance[TController]()
	requireAuth := isSpecificType[TIdentifier]()
	handlerMapping[uri] = func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {

				fmt.Println(string(debug.Stack()))
			}
		}()
		if initControllerError != nil {
			if currentServer.IsReleaseMode {
				log.Panicln(initControllerError)
				writeJSONResponse(w, http.StatusInternalServerError, "server error")
				return
			} else {
				panic(initControllerError)
			}
		}
		if r.Method != "GET" {
			writeJSONResponse(w, http.StatusBadGateway, map[string]string{"error": "Method is not allow"})
			return
		}

		ctxPool := newContextPool[TIdentifier]()
		ctx := ctxPool.get()
		ctx.uriParamsList = uriInfo.uriParams
		ctx.RootRouter = currentServer.BaseUrl
		ctx.Req = r
		ctx.Res = w
		ctx.uriRegex = uriInfo.uriRegex
		ctx.Uri = uriInfo.mainUri
		if uriInfo.isRegex {
			ctx.loadAllParams()
		}
		defer func() {
			ctx.Reset()
			ctxPool.put(ctx)
		}()
		if requireAuth {
			err := callAuth(controllerInstance, ctx)
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

func HandlerForm[TController any, TIdentifier any, TData any, TResponse any](
	uriHandler string,
	fn func(controllerInstance *TController, ctx *HttpContext[TIdentifier], data TData) (TResponse, error)) {

	uri := fmt.Sprintf("%s/%s", strings.ToLower(reflect.TypeFor[TController]().Name()), uriHandler)
	uriInfo := inspectUri(uri)
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
		isRequireAuth:       isSpecificType[TIdentifier](),
		uriInfo:             uriInfo,
	}
	swaggerInfo[uriInfo.mainUri] = swaggerItem
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	if swaggerItem.requestBodyType.Kind() == reflect.Ptr {
		swaggerItem.requestBodyType = swaggerItem.requestBodyType.Elem()
	}
	hasBody := isSpecificType[TData]()
	var initControllerError error
	controllerInstance, initControllerError := newControllerInstance[TController]()
	requireAuth := isSpecificType[TIdentifier]()
	handlerMapping[uri] = func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {

				fmt.Println(string(debug.Stack())) // debug.Stack() trả về []byte chứa toàn bộ stack trace
			}
		}()

		if initControllerError != nil {
			if currentServer.IsReleaseMode {
				log.Panicln(initControllerError)
				writeJSONResponse(w, http.StatusInternalServerError, "server error")
				return
			} else {
				panic(initControllerError)
			}
		}
		if r.Method != "POST" {
			writeJSONResponse(w, http.StatusBadGateway, map[string]string{"error": "Method is not allow"})
			return
		}
		ctxPool := newContextPool[TIdentifier]()
		ctx := ctxPool.get()
		ctx.uriParamsList = uriInfo.uriParams
		ctx.RootRouter = currentServer.BaseUrl
		ctx.Req = r
		ctx.Res = w
		ctx.uriRegex = uriInfo.uriRegex
		ctx.Uri = uriInfo.mainUri
		if uriInfo.isRegex {
			ctx.loadAllParams()
		}
		defer func() {
			ctx.Reset()
			ctxPool.put(ctx)
		}()
		if requireAuth {
			err := callAuth(controllerInstance, ctx)
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
func GetHandler[TController any](uriHandler string) (http.HandlerFunc, error) {
	var ret http.HandlerFunc
	uri := fmt.Sprintf("%s/%s", strings.ToLower(reflect.TypeFor[TController]().Name()), uriHandler)
	if fx, ok := handlerMapping[uri]; ok {
		return fx, nil
	}
	return ret, fmt.Errorf("%s not found in %s", uriHandler, reflect.TypeFor[TController]().String())

}
