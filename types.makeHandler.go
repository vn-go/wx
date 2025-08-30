package wx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

func (h *handlerInfo) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != h.httpMethod {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		contentType := r.Header.Get("Content-Type")

		ret, err := h.Invoke(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if contentType == "application/json" && !h.isNoOutPut {
			w.Header().Set("Content-Type", "application/json")
			//write ret to
			if err := json.NewEncoder(w).Encode(ret); err != nil {
				http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
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
func (info *handlerInfo) Invoke(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	contentType := r.Header.Get("Content-Type")
	controller, err := info.CreateController()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	valueOfReq := reflect.ValueOf(r)
	valueOfRes := reflect.ValueOf(w)
	httpContextValue, err := info.CreateHttpContext(valueOfReq, valueOfRes)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	if info.conrollerNewMethod != nil {
		retRun := info.conrollerNewMethod.Func.Call([]reflect.Value{*controller})
		if retRun[0].Elem().Interface() != nil {
			// http.Error(w, retRun[0].Elem().Interface().(error).Error(), http.StatusInternalServerError)
			return nil, retRun[0].Elem().Interface().(error)

		}

	}
	if info.hasHttpContextInController {
		fiedlOfHttpContextInController := controller.Elem().FieldByIndex(info.fiedIndexOfHttpContextInController)
		if fiedlOfHttpContextInController.Kind() == reflect.Ptr {
			fiedlOfHttpContextInController.Set(reflect.ValueOf(&HttpContext{
				Req: r,
				Res: w,
			}))
		} else {
			fiedlOfHttpContextInController.Set(reflect.ValueOf(HttpContext{
				Req: r,
				Res: w,
			}))
		}
	}
	args := make([]reflect.Value, info.method.Type.NumIn())
	args[0] = *controller
	args[info.indexOfArgIsHttpContext] = *httpContextValue
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
	if retRun[len(retRun)-1].Elem().Interface() != nil {
		if err, ok := retRun[len(retRun)-1].Elem().Interface().(error); ok {
			return nil, err
		}
		return retRun, nil
	}
	return retRun[0 : len(retRun)-1], nil
}

type initCreateControllerOnce struct {
	controller *reflect.Value
	err        error
	once       sync.Once
}

var cacheCreateControllerOnce sync.Map

func (info *handlerInfo) CreateControllerOnce() (*reflect.Value, error) {
	actally, _ := cacheCreateControllerOnce.LoadOrStore(info.controllerTypeElem, &initCreateControllerOnce{})
	item := actally.(*initCreateControllerOnce)
	item.once.Do(func() {
		item.controller, item.err = info.CreateController()
	})
	return item.controller, item.err

}
func (info *handlerInfo) CreateController() (*reflect.Value, error) {
	controllerValue := reflect.New(info.controllerType.Elem())
	if info.conrollerNewMethod != nil {
		ret := info.conrollerNewMethod.Func.Call([]reflect.Value{controllerValue})
		if ret[0].Elem().Interface() != nil {
			return nil, ret[0].Elem().Interface().(error)
		}
		return &controllerValue, nil
	}
	return &controllerValue, nil

}
func (info *handlerInfo) CreateHttpContext(reqValue, resValue reflect.Value) (*reflect.Value, error) {
	httpContextValue := reflect.New(info.typeOfArgIsHttpContextElem)
	httpContextValue.Elem().FieldByIndex(info.reqFieldIndex).Set(reqValue)
	httpContextValue.Elem().FieldByIndex(info.resFieldIndex).Set(resValue)
	return &httpContextValue, nil
}
func (info *handlerInfo) GetBodyValue(r *http.Request, contentType string) (reflect.Value, error) {
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
