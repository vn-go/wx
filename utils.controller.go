package wx

import (
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
