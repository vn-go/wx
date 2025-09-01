package wx

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/vn-go/wx/internal"
)

type controllerHelperType struct {
	httpContextTypePtr reflect.Type
	httpContextType    reflect.Type
}
type initFindIndexOfFieldHttpContext struct {
	val   []int
	found bool
	once  sync.Once
}

var cacheFindIndexOfFieldHttpContext sync.Map

func (h *controllerHelperType) isHandler(typ reflect.Type) bool {
	// if field.Type == reflect.TypeOf(&HttpContext{}) || field.Type == reflect.TypeOf(HttpContext{}) {

	// 	return true
	// }
	return typ == h.httpContextType || typ == h.httpContextTypePtr

	//return false
}
func (h *controllerHelperType) findIndexOfFieldHttpContextInternal(receiverType reflect.Type, visited map[reflect.Type]bool) ([]int, bool) {

	if receiverType.Kind() == reflect.Ptr {
		receiverType = receiverType.Elem()
	}

	if visited[receiverType] {
		return nil, false
	}
	visited[receiverType] = true

	for i := 0; i < receiverType.NumField(); i++ {
		field := receiverType.Field(i)
		if h.isHandler(field.Type) {

			return field.Index, true
		}
		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
		}
		if field.Type.Kind() == reflect.Struct {
			if ret, found := h.findIndexOfFieldHttpContextInternal(field.Type, visited); found {
				return append(field.Index, ret...), true
			}
		}
	}
	return nil, false
}
func (h *controllerHelperType) FindIndexOfFieldHttpContext(receiverType reflect.Type) ([]int, bool) {
	acctually, _ := cacheFindIndexOfFieldHttpContext.LoadOrStore(receiverType.String(), &initFindIndexOfFieldHttpContext{})
	item := acctually.(*initFindIndexOfFieldHttpContext)
	item.once.Do(func() {
		item.val, item.found = h.findIndexOfFieldHttpContextInternal(receiverType, map[reflect.Type]bool{})

	})
	return item.val, item.found

}
func (h *controllerHelperType) FindReqResFieldIndexDelete(receiverType reflect.Type) ([]int, []int) {
	fiedIndexOfHttpContext, found := h.FindIndexOfFieldHttpContext(receiverType)
	if !found {
		return nil, nil
	}
	fiedlOfHttpContext := receiverType.FieldByIndex(fiedIndexOfHttpContext)
	fiedlOfHttpContextType := fiedlOfHttpContext.Type
	if fiedlOfHttpContextType.Kind() == reflect.Ptr {
		fiedlOfHttpContextType = fiedlOfHttpContextType.Elem()
	}
	fieldReq, okReq := fiedlOfHttpContextType.FieldByName("Req")
	fieldRes, okRes := fiedlOfHttpContextType.FieldByName("Res")
	if !okReq || !okRes {
		return nil, nil
	}
	return append(fiedIndexOfHttpContext, fieldReq.Index...), append(fiedIndexOfHttpContext, fieldRes.Index...)
	//return fieldReq.Index, fieldRes.Index

}

func (h *controllerHelperType) FindControllerName(reiverType reflect.Type) string {
	if reiverType.Kind() == reflect.Ptr {
		reiverType = reiverType.Elem()
	}
	if reiverType.Kind() != reflect.Struct {
		return ""
	}
	key := reiverType.String() + "/FindControllerName"
	ret, _ := internal.OnceCall(key, func() (*string, error) {

		for i := 0; i < reiverType.NumField(); i++ {
			field := reiverType.Field(i)
			tags := field.Tag.Get("controller")
			if tags != "" {
				tags = h.ToKebabCase(tags)

				return &tags, nil
			}
		}

		items := strings.Split(reiverType.String(), ".")
		ret := h.ToKebabCase(items[len(items)-1])
		return &ret, nil
		/* find first posistion of  "/controllers/" */

	})
	return *ret
}
func (h *controllerHelperType) ToKebabCase(s string) string {
	key := s + "/helperType/ToKebabCase"
	ret, _ := internal.OnceCall(key, func() (*string, error) {
		// Khớp các chữ cái viết hoa và thêm dấu gạch ngang trước đó.
		// Ví dụ: MyMethod -> -My-Method
		re := regexp.MustCompile("([A-Z])")
		snake := re.ReplaceAllString(s, "-$1")

		// Chuyển toàn bộ chuỗi sang chữ thường và loại bỏ dấu gạch ngang ở đầu nếu có.
		// Ví dụ: -My-Method -> -my-method -> my-method
		ret := strings.ToLower(strings.TrimPrefix(snake, "-"))
		return &ret, nil
	})
	return *ret
}
func (h *controllerHelperType) FindNewMeyhod(info *handlerInfo) (*reflect.Method, error) {
	for i := 0; i < info.controllerType.NumMethod(); i++ {
		method := info.controllerType.Method(i)
		if method.Name == "New" {
			if method.Type.NumOut() != 1 {
				return nil, fmt.Errorf("%s.New() must return 1 error", info.controllerType.String())
			}
			if !method.Type.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				return nil, fmt.Errorf("%s.New() must return error", info.controllerType.String())

			}

			return &method, nil
		}
	}
	return nil, nil
}
