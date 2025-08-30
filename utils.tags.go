package wx

import (
	"reflect"
	"strings"
)

type tagsHelperType struct {
}

func (h *tagsHelperType) ExtractTags(typ reflect.Type) ([]string, bool) {
	var found bool
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return []string{}, false
	}
	ret := []string{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldType := field.Type
		if typ.ConvertibleTo(utils.HttpContextType) || typ.ConvertibleTo(utils.HttpContextPtrType) {

			return ret, true
		}
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			subRet, found := h.ExtractTags(fieldType)
			if found {
				if field.Tag.Get("route") != "" {

					ret = append(ret, field.Tag.Get("route"))
				}

			}
			ret = append(ret, subRet...)
		}

	}

	return ret, found

}
func (h *tagsHelperType) ExtractUriFromTags(tags []string) string {
	ret := ""
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]
		if tag == "" {
			continue
		}
		items := strings.Split(tag, ";")

		for _, item := range items {
			uriVal := ""
			if strings.HasPrefix(item, "uri:") {
				uriVal = item[4:]

			} else if item != "" && !strings.Contains(item, ":") {
				uriVal = item
			}
			if uriVal != "" {
				if strings.Contains(ret, "@") {
					ret = strings.Replace(ret, "@", uriVal, 1)
				} else {
					ret += "/" + uriVal
				}
			}
		}

	}
	ret = strings.TrimPrefix(strings.TrimSuffix(ret, "/"), "/")
	return ret

}
func (h *tagsHelperType) ExtractHttpMethodFromTags(tags []string) string {
	ret := ""
	for i := len(tags) - 1; i >= 0; i-- {
		tag := tags[i]
		if tag == "" {
			continue
		}
		items := strings.Split(tag, ";")
		for _, item := range items {
			if strings.HasPrefix(item, "method:") {
				ret = strings.ToUpper(item[7:])

			}
		}
	}

	return ret
}
