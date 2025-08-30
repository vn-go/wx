package wx

type invokerType struct {
}

// func (invoker *invokerType) Invoke(info handlerInfo, w http.ResponseWriter, r *http.Request) {
// 	controller, err := utils.controllers.Create(info)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	reqValue := reflect.ValueOf(r)
// 	resValue := reflect.ValueOf(w)

// 	httpContextValue, err := utils.controllers.CreateHttpContext(info, reqValue, resValue)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	info.method.Func.Call([]reflect.Value{*controller, *httpContextValue})
// }

// var invoker = &invokerType{}
