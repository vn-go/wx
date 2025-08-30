package wx

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/vn-go/wx/internal"
)

type controllerHelperType struct {
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

func (h *controllerHelperType) Create(info *handlerInfo) (*reflect.Value, error) {
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
