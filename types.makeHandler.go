package wx

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/vn-go/wx/internal"
)

// helper for consistent JSON payload
func toJSON(code, message string) []byte {
	resp := map[string]string{
		"error":   code,
		"message": message,
	}
	b, _ := json.Marshal(resp)
	return b
}
func (h *handlerInfo) catchError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var status int
	var body []byte

	switch e := err.(type) {
	// ------------------- 4xx Client Errors -------------------
	case *BadRequestError, *ParamMissMatchError, *UriParamParseError,
		*UriParamConvertError, *RequireError, *BodyParseError, *FileParseError:
		status = http.StatusBadRequest // 400
		body = toJSON("bad_request", e.Error())

	case *MethodNotAllowError:
		status = http.StatusMethodNotAllowed // 405
		body = toJSON("method_not_allowed", e.Error())

	case *NewMethodOfAuthNotFoundError:
		status = http.StatusUnauthorized // 401
		body = toJSON("unauthorized", e.Error())

	case *RegexUriNotMatchError:
		status = http.StatusNotFound // 404
		body = toJSON("not_found", e.Error())

	case *UnSupportError, *UnacceptableContentError:
		status = http.StatusUnsupportedMediaType // 415
		body = []byte(e.Error())                 // already JSON in UnacceptableContentError

	// ------------------- 5xx Server Errors -------------------
	case *ServiceInitError, *ServerError:
		status = http.StatusInternalServerError // 500
		body = toJSON("server_error", e.Error())
	case *UnauthorizedError:
		status = http.StatusUnauthorized // 500
		body = toJSON("Unauthorized", e.Error())

	default:
		// fallback for unexpected error
		status = http.StatusInternalServerError
		body = toJSON("server_error", err.Error())
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// func (h *handlerInfo) getVerifyMethodOfAuth() {

//		var found bool
//		h.newMethodOfAuth, found = authUtils.GetNewMethod(h.typeOfFiedAuth)
//		if !found {
//			return nil, fmt.Errorf("%s.Verify was not call, please call%s.Verify for setting up auth ", ret.typeOfFiedAuth.String(), ret.typeOfFiedAuth.String())
//		}
//	}
func (h *handlerInfo) getAuth(valueOfHandler reflect.Value, w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
	newMethodOfAuth, found := authUtils.GetNewMethod(h.typeOfFiedAuth)
	if !found {
		err := fmt.Errorf("%s.Verify was not call, please call%s.Verify for setting up auth ", h.typeOfFiedAuth.String(), h.typeOfFiedAuth.String())
		return reflect.Value{}, NewServerError("server error", err)
	}
	ret := reflect.New(h.typeOfFiedAuth)

	retRun := newMethodOfAuth.Call([]reflect.Value{valueOfHandler})
	last := retRun[len(retRun)-1]        // last return value
	if last.IsValid() && !last.IsNil() { // safe checks
		if err, ok := last.Interface().(error); ok {
			return reflect.Value{}, err
		}
	}
	ret.Elem().FieldByName("Data").Set(retRun[0])
	if h.typeOfFiedAuth.Kind() == reflect.Struct {
		return ret.Elem(), nil
	}
	return ret, nil

}
func (h *handlerInfo) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != h.httpMethod {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		//contentType := r.Header.Get("Content-Type")

		ret, err := h.Invoke(w, r)
		if err != nil {
			h.catchError(w, err)
			return
		}
		if ret != nil && !h.isNoOutPut {
			if len(ret) > 1 {
				retData := []any{}
				for _, x := range ret {
					if x.Kind() == reflect.Ptr {
						x = x.Elem()
					}
					retData = append(retData, x.Interface())

				}
				w.Header().Set("Content-Type", "application/json")
				//write ret to
				if err := json.NewEncoder(w).Encode(retData); err != nil {
					h.catchError(w, NewServerError("Internal server error", err))
					//http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				//write ret to
				if ret[0].Kind() == reflect.Ptr {
					ret[0] = ret[0].Elem()

				}
				if ret[0].IsValid() {
					if err := json.NewEncoder(w).Encode(ret[0].Interface()); err != nil {
						h.catchError(w, NewServerError("Internal server error", err))
					}
				} else {
					if err := json.NewEncoder(w).Encode(nil); err != nil {
						h.catchError(w, NewServerError("Internal server error", err))
					}
				}

			}

		}

	}
}
func (h *handlerInfo) GetUriHandler() string {
	return h.uriHandler
}
func (h *handlerInfo) GetHttpMethod() string {
	return h.httpMethod
}
func (info *handlerInfo) createHandler(w http.ResponseWriter, r *http.Request) Handler {

	ret := func() *httpContext {
		return &httpContext{
			Req: r,
			Res: w,
		}
	}
	return ret
}
func (info *handlerInfo) Invoke(w http.ResponseWriter, r *http.Request) ([]reflect.Value, error) {

	contentType := r.Header.Get("Content-Type")
	valueOfArgsIsHandler, valueOfHandlerFunction := info.CreateHandlerValue(r, w)
	controller, err := info.CreateController(valueOfHandlerFunction)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	var authValue reflect.Value
	if info.isAuth {
		var err error
		authValue, err = info.getAuth(valueOfHandlerFunction, w, r)
		if err != nil {
			return nil, err
		}
	}
	// valueOfReq := reflect.ValueOf(r)
	// valueOfRes := reflect.ValueOf(w)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	err = info.applyUri(valueOfArgsIsHandler, r)
	if err != nil {
		return nil, err
	}
	if info.hasHttpContextInController {
		fiedlOfhttpContextInController := controller.Elem().FieldByIndex(info.indexFieldIsHandlerInController)
		if fiedlOfhttpContextInController.Kind() == reflect.Ptr {
			handler := info.createHandler(w, r)
			fiedlOfhttpContextInController.Set(reflect.ValueOf(&handler))
		} else {
			//"reflect: Call using *wx.OK[github.com/vn-go/wx.Handler] as type wx.OK[github.com/vn-go/wx.Handler]"
			handler := info.createHandler(w, r)
			fiedlOfhttpContextInController.Set(reflect.ValueOf(handler))
			//fiedlOfhttpContextInController = fiedlOfhttpContextInController.Addr()
		}
	}

	args := make([]reflect.Value, info.method.Type.NumIn())
	args[0] = *controller
	args[info.indexOfArgIsHandler] = valueOfArgsIsHandler
	if info.indexOfArgIsAuth > -1 {
		if args[info.indexOfArgIsAuth].Kind() == reflect.Ptr {
			if len(info.fieldIndexOfAuth) > 0 {
				args[info.indexOfArgIsAuth].Elem().FieldByIndex(info.fieldIndexOfAuth).Set(authValue)
			} else {
				args[info.indexOfArgIsAuth] = authValue
			}

		} else {
			if len(info.fieldIndexOfAuth) > 0 {
				args[info.indexOfArgIsAuth].Elem().FieldByIndex(info.fieldIndexOfAuth).Set(authValue)
			} else {
				args[info.indexOfArgIsAuth] = authValue
			}

		}
	}

	if info.indexOfArgIsRequestBody != -1 {
		bodyValue, err := info.GetBodyValue(r, contentType)
		if err != nil {
			return nil, err
		}
		args[info.indexOfArgIsRequestBody] = bodyValue

	}
	retRun := info.method.Func.Call(args)
	if len(retRun) == 0 {
		return nil, nil
	}
	last := retRun[len(retRun)-1]        // last return value
	if last.IsValid() && !last.IsNil() { // safe checks
		if err, ok := last.Interface().(error); ok {
			return nil, err
		}
	}
	return retRun[0 : len(retRun)-1], nil
}

type initCreateControllerOnce struct {
	controller *reflect.Value
	err        error
	once       sync.Once
}

var cacheCreateControllerOnce sync.Map

// func (info *handlerInfo) CreateControllerOnce() (*reflect.Value, error) {
// 	actally, _ := cacheCreateControllerOnce.LoadOrStore(info.controllerTypeElem, &initCreateControllerOnce{})
// 	item := actally.(*initCreateControllerOnce)
// 	item.once.Do(func() {
// 		item.controller, item.err = info.CreateController()
// 	})
// 	return item.controller, item.err

// }
func (info *handlerInfo) CreateController(valueOfHandlerFunction reflect.Value) (*reflect.Value, error) {
	controllerValue := reflect.New(info.controllerType.Elem())
	if info.conrollerNewMethod != nil {
		if info.indexFieldIsHandlerInController != nil {
			controllerValue.Elem().FieldByIndex(info.indexFieldIsHandlerInController).Set(valueOfHandlerFunction)
		}

		ret := info.conrollerNewMethod.Func.Call([]reflect.Value{controllerValue})
		if !ret[0].IsZero() {
			if ret[0].Elem().Interface() != nil {
				return nil, ret[0].Elem().Interface().(error)
			}
		}
		return &controllerValue, nil
	}
	return &controllerValue, nil

}
func (info *handlerInfo) CreateHandlerValue(r *http.Request, w http.ResponseWriter) (reflect.Value, reflect.Value) {
	if utils.controllers.isHandler(info.typeOfArgIsIsHandlerElem) {
		if info.typeOfArgIsIsHandler.Kind() == reflect.Ptr {
			// type Hnx func ()
			// var h Hnx = func() {
			// 	fmt.Println("Hello Hnx")
			// }

			// var fx *Hnx = &h // fx là *Hnx

			var retVale Handler = func() *httpContext {
				return &httpContext{
					Req: r,
					Res: w,
				}
			}
			var retValePtr *Handler = &retVale

			retVal := reflect.ValueOf(retValePtr)
			return retVal, retVal.Elem()
		} else {
			ret := func() *httpContext {
				return &httpContext{
					Req: r,
					Res: w,
				}
			}
			retVal := reflect.ValueOf(ret)
			return retVal, retVal
		}

	}
	retValOfHandler := reflect.New(info.typeOfArgIsIsHandlerElem)

	ret := func() *httpContext {
		return &httpContext{
			Req: r,
			Res: w,
		}
	}
	retValOfHandlerFn := reflect.ValueOf(ret)
	retValOfHandler.Elem().FieldByIndex(info.indexFieldIshandler).Set(retValOfHandlerFn)
	if info.typeOfArgIsIsHandler.Kind() == reflect.Struct {
		retValOfHandler = retValOfHandler.Elem()
	}

	return retValOfHandler, retValOfHandlerFn
	// httpContextValue := reflect.New(info.typeOfArgIshttpContextElem)
	// httpContextValue.Elem().FieldByIndex(info.reqFieldIndex).Set(reqValue)
	// httpContextValue.Elem().FieldByIndex(info.resFieldIndex).Set(resValue)
	// return &httpContextValue, nil
}
func (info *handlerInfo) CreatehttpContextDelete(reqValue, resValue reflect.Value) (*reflect.Value, error) {
	// httpContextValue := reflect.New(info.typeOfArgIshttpContextElem)
	// httpContextValue.Elem().FieldByIndex(info.reqFieldIndex).Set(reqValue)
	// httpContextValue.Elem().FieldByIndex(info.resFieldIndex).Set(resValue)
	// return &httpContextValue, nil
	panic("depreciate")
}
func (info *handlerInfo) GetBodyValue(r *http.Request, contentType string) (reflect.Value, error) {
	//"multipart/form-data; boundary=bc93ed97d895d9ff5f8eb8f994205bc3f8184e4f5d0668de8791448fe447"

	if strings.HasPrefix(contentType, "multipart/form-data; ") {
		return info.GetMultipartFormDataValue(r)
	}
	if contentType == "application/x-www-form-urlencoded" {
		return info.GetXWwwFormUrlencoded(r)
	}
	if contentType == "application/json" && info.isFormPost {
		/*
						{
			  "error": "unsupported_media_type",
			  "message": "Content-Type application/json is not supported. Please use multipart/form-data."
			}
		*/
		return reflect.Value{}, NewUnacceptableContentError(
			"unsupported_media_type",
			"Content-Type application/json is not supported. Please use multipart/form-data.",
		)
	}

	bodyData := reflect.New(info.typeOfRequestBodyElem)
	if r.Body != nil && r.Body != http.NoBody {
		defer r.Body.Close() // <-- auto close after read body of request
		if err := json.NewDecoder(r.Body).Decode(bodyData.Interface()); err != nil {

			return reflect.Value{}, err
		}
	} else if info.typeOfRequestBody.Kind() == reflect.Struct {

		return reflect.Value{}, NewBadRequestError("request body is required")
	}

	return bodyData, nil
}
func (info *handlerInfo) getFieldByName(typ reflect.Type, fieldName string) *reflect.StructField {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	key := typ.String() + "/RequestExecutor/GetFieldByName/" + fieldName
	ret, err := internal.OnceCall(key, func() (*reflect.StructField, error) {
		ret, ok := typ.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, fieldName)
		})
		if !ok {
			return nil, nil
		}
		return &ret, nil
	})
	if err != nil {
		return nil
	}
	return ret
}
func (info *handlerInfo) getMaxMemory() int64 {
	return 20 << 20
}
func (info *handlerInfo) GetMultipartFormDataValue(r *http.Request) (reflect.Value, error) {
	if isFormType(info.typeOfRequestBodyElem) {
		ret := reflect.New(info.typeOfRequestBodyElem)
		if fieldData, ok := info.typeOfRequestBodyElem.FieldByName("Data"); ok {
			dataVal, err := info.getMultipartFormDataValueByType(fieldData.Type, r)
			if err != nil {
				return reflect.Value{}, NewServerError("Internal server error", err)
			}
			fieldSet := ret.Elem().FieldByIndex(fieldData.Index)
			if fieldSet.CanConvert(dataVal.Type()) {
				fieldSet.Set(dataVal)
			} else if dataVal.Kind() == reflect.Ptr {
				dataVal = dataVal.Elem()
				if fieldSet.CanConvert(dataVal.Type()) {
					fieldSet.Set(dataVal)
				}

			}

			if info.typeOfRequestBody.Kind() == reflect.Struct {
				return ret.Elem(), nil
			}
			return ret, nil
		} else {
			return reflect.Value{}, NewServerError("Internal server error", fmt.Errorf("%s do not have Data Field", info.typeOfRequestBodyElem.String()))
		}

	} else {

		return info.getMultipartFormDataValueByType(info.typeOfRequestBody, r)
	}
}
func (info *handlerInfo) getXWwwFormUrlencoded(bodyType reflect.Type, r *http.Request) (reflect.Value, error) {
	var target reflect.Value
	var targetType reflect.Type

	ret := reflect.New(bodyType)

	target = ret.Elem()
	targetType = bodyType
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	err := r.ParseForm()
	if err != nil {
		return reflect.Value{}, NewBodyParseError("can not read data", err)
	}
	for key, values := range r.Form {
		field := info.getFieldByName(targetType, key)
		if field == nil || len(values) == 0 {
			continue
		}
		fv := target.FieldByIndex(field.Index)

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(values[0])
		case reflect.Slice:
			if fv.Type().Elem().Kind() == reflect.String {
				fv.Set(reflect.ValueOf(values))
			}
		case reflect.Ptr:
			elemKind := fv.Type().Elem().Kind()
			if elemKind == reflect.String {
				ptr := reflect.New(fv.Type().Elem())
				ptr.Elem().SetString(values[0])
				fv.Set(ptr)
			} else if elemKind == reflect.Struct {
				ptr := reflect.New(fv.Type().Elem())
				if err := json.Unmarshal([]byte(values[0]), ptr.Interface()); err != nil {
					return reflect.Value{}, err
				}
				fv.Set(ptr)
			}
		case reflect.Struct:
			if err := json.Unmarshal([]byte(values[0]), fv.Addr().Interface()); err != nil {
				return reflect.Value{}, err
			}
		}
	}

	// set files

	if bodyType.Kind() == reflect.Struct {
		return ret.Elem(), nil
	}

	return ret, nil
}
func (info *handlerInfo) GetXWwwFormUrlencoded(r *http.Request) (reflect.Value, error) {
	if isFormType(info.typeOfRequestBodyElem) {
		ret := reflect.New(info.typeOfRequestBodyElem)
		if fieldData, ok := info.typeOfRequestBodyElem.FieldByName("Data"); ok {
			dataVal, err := info.getXWwwFormUrlencoded(fieldData.Type, r)
			if err != nil {
				return reflect.Value{}, NewServerError("Internal server error", err)
			}
			fieldSet := ret.Elem().FieldByIndex(fieldData.Index)
			if fieldSet.CanConvert(dataVal.Type()) {
				fieldSet.Set(dataVal)
			} else if dataVal.Kind() == reflect.Ptr {
				dataVal = dataVal.Elem()
				if fieldSet.CanConvert(dataVal.Type()) {
					fieldSet.Set(dataVal)
				}

			}

			if info.typeOfRequestBody.Kind() == reflect.Struct {
				return ret.Elem(), nil
			}
			return ret, nil
		} else {
			return reflect.Value{}, NewServerError("Internal server error", fmt.Errorf("%s do not have Data Field", info.typeOfRequestBodyElem.String()))
		}

	} else {

		return info.getMultipartFormDataValueByType(info.typeOfRequestBody, r)
	}
}
func (info *handlerInfo) getMultipartFormDataValueByType(bodyType reflect.Type, r *http.Request) (reflect.Value, error) {
	var target reflect.Value
	var targetType reflect.Type

	ret := reflect.New(bodyType)

	target = ret.Elem()
	targetType = bodyType
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	if bodyType == reflect.TypeFor[multipart.FileHeader]() {
		if err := r.ParseMultipartForm(info.getMaxMemory()); err != nil {
			return reflect.Value{}, NewFileParseError("error parsing multipart form", err)
		}
		for _, v := range r.MultipartForm.File {
			if len(v) > 0 {
				return reflect.ValueOf(*v[0]), nil
			}

		}
		return reflect.Value{}, NewRequireError([]string{}, "file upload is missing")
	}
	// 2- none-require file upload
	if bodyType == reflect.TypeFor[*multipart.FileHeader]() {
		if err := r.ParseMultipartForm(info.getMaxMemory()); err != nil {
			return reflect.Value{}, NewFileParseError("error parsing multipart form", err)
		}
		for _, v := range r.MultipartForm.File {
			if len(v) > 0 {
				return reflect.ValueOf(v[0]), nil
			}

		}
		return ret, nil
	}
	// 3- multifile upload
	if bodyType == reflect.TypeFor[[]multipart.FileHeader]() {
		if err := r.ParseMultipartForm(info.getMaxMemory()); err != nil {
			return reflect.Value{}, NewFileParseError("error parsing multipart form", err)
		}
		for _, v := range r.MultipartForm.File {
			retFiles := []multipart.FileHeader{}

			for _, x := range v {
				retFiles = append(retFiles, *x)
			}
			return reflect.ValueOf(retFiles), nil

		}
		return ret, nil
	}
	// 4- multifile upload with nullable of element
	if bodyType == reflect.TypeFor[[]*multipart.FileHeader]() {
		if err := r.ParseMultipartForm(info.getMaxMemory()); err != nil {
			return reflect.Value{}, NewFileParseError("error parsing multipart form", err)
		}
		for _, v := range r.MultipartForm.File {

			return reflect.ValueOf(v), nil

		}
		return ret, nil
	}
	// 5- multifile upload with nullable of element and nullable of array
	if bodyType == reflect.TypeFor[*[]*multipart.FileHeader]() {
		if err := r.ParseMultipartForm(info.getMaxMemory()); err != nil {
			return reflect.Value{}, NewFileParseError("error parsing multipart form", err)
		}
		for _, v := range r.MultipartForm.File {
			retFiles := []*multipart.FileHeader{}
			retFiles = append(retFiles, v...)
			return reflect.ValueOf(&retFiles), nil

		}
		return ret, nil
	}
	var formData map[string][]string
	var files map[string][]*multipart.FileHeader
	//fieldsIsFile := reflect.StructField{}
	//contentType := r.Header.Get("Content-Type")
	if err := r.ParseMultipartForm(info.getMaxMemory()); err != nil {
		return reflect.Value{}, NewFileParseError("error parsing multipart form", err)
	}
	formData = r.MultipartForm.Value
	files = r.MultipartForm.File //<-- wx tu dong lay file theo kieu nay

	for key, values := range formData {
		field := info.getFieldByName(targetType, key)
		if field == nil || len(values) == 0 {
			continue
		}
		fv := target.FieldByIndex(field.Index)

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(values[0])
		case reflect.Slice:
			if fv.Type().Elem().Kind() == reflect.String {
				fv.Set(reflect.ValueOf(values))
			}
		case reflect.Ptr:
			elemKind := fv.Type().Elem().Kind()
			if elemKind == reflect.String {
				ptr := reflect.New(fv.Type().Elem())
				ptr.Elem().SetString(values[0])
				fv.Set(ptr)
			} else if elemKind == reflect.Struct {
				ptr := reflect.New(fv.Type().Elem())
				if err := json.Unmarshal([]byte(values[0]), ptr.Interface()); err != nil {
					return reflect.Value{}, err
				}
				fv.Set(ptr)
			}
		case reflect.Struct:
			if err := json.Unmarshal([]byte(values[0]), fv.Addr().Interface()); err != nil {
				return reflect.Value{}, err
			}
		}
	}

	// set files
	for key, fhArr := range files {
		field := info.getFieldByName(targetType, key)
		if field == nil || len(fhArr) == 0 {
			continue
		}
		fv := target.FieldByIndex(field.Index)
		ft := fv.Type()

		switch {
		case ft == reflect.TypeOf(&multipart.FileHeader{}):
			fv.Set(reflect.ValueOf(fhArr[0]))
		case ft == reflect.TypeOf([]*multipart.FileHeader{}):
			fv.Set(reflect.ValueOf(fhArr))
		case ft == reflect.TypeOf([]multipart.FileHeader{}):
			slice := make([]multipart.FileHeader, len(fhArr))
			for i, f := range fhArr {
				slice[i] = *f
			}
			fv.Set(reflect.ValueOf(slice))
		}
	}
	if bodyType.Kind() == reflect.Struct {
		return ret.Elem(), nil
	}

	return ret, nil
}
func (info *handlerInfo) GetParamFieldOfHandlerContext(typ reflect.Type, fieldName string) (reflect.StructField, bool) {
	key := typ.String() + "/" + fieldName
	ret, err := internal.OnceCall(key, func() (*reflect.StructField, error) {
		field, ok := typ.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, fieldName)
		})
		if !ok {
			return nil, nil
		}
		return &field, nil
	})
	if err != nil {
		return reflect.StructField{}, false
	}
	if ret == nil {
		return reflect.StructField{}, false
	}
	return *ret, true

}
func (info *handlerInfo) applyUri(contextValue reflect.Value, r *http.Request) error {
	if info.isRegexHandler {
		placeHolders := info.regexUriFind.FindAllStringSubmatch(r.URL.Path, -1)
		if len(placeHolders) == 0 {

			return NewRegexUriNotMatchError("regex uri not match")
		}
		for i := 1; i < len(placeHolders[0]); i++ {
			fieldIndex := info.uriParams[i-1].FieldIndex
			fieldSet := contextValue.Elem().FieldByIndex(fieldIndex)
			fieldSet.Set(reflect.ValueOf(placeHolders[0][i]))
			//field := contextValue.Type().Elem().FieldByIndex(fieldIndex)

		}
		if info.isQueryUri {
			// url, err := r.URL.Parse(r.URL.Path)
			// if err != nil {
			// 	return NewServerError("can not read url", err)
			// }

			query := r.URL.Query()

			typeOfContextValue := contextValue.Type().Elem()
			for k, x := range query {
				field, ok := info.GetParamFieldOfHandlerContext(typeOfContextValue, k)
				if ok {
					fieldSet := contextValue.Elem().FieldByIndex(field.Index)
					if fieldSet.IsValid() {
						if fieldSet.Kind() == reflect.String {
							fieldSet.SetString(x[0])
						} else if fieldSet.Kind() == reflect.Ptr {
							if fieldSet.Type().Elem().Kind() == reflect.String {
								fieldSet.Set(reflect.ValueOf(&x[0]))
							}
						} else if fieldSet.Kind() == reflect.Slice {
							if fieldSet.Type().Elem().Kind() == reflect.String {
								fieldSet.Set(reflect.ValueOf(x))
							} else if fieldSet.Type().Elem().Kind() == reflect.Ptr {
								if fieldSet.Type().Elem().Elem().Kind() == reflect.String {
									vals := make([]*string, len(x))
									for i, v := range x {
										vals[i] = &v
									}

									fieldSet.Set(reflect.ValueOf(vals))
								}
							}
						}

					}
				}

			}

		}
	}
	return nil
}
