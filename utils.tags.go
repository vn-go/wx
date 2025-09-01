package wx

import (
	"reflect"
	"strings"
)

type tagsHelperType struct {
}

func (h *tagsHelperType) ExtractTags(typ reflect.Type, visited map[reflect.Type]bool) ([]string, bool) {
	ret := []string{}
	if typ == utils.HttpContextType || typ == utils.HttpContextPtrType {
		return ret, true
	}
	var found bool
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if _, ok := visited[typ]; ok {
		return ret, false
	}
	visited[typ] = true
	if typ.Kind() != reflect.Struct {
		return []string{}, false
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		if fieldType == utils.HttpContextType || fieldType == utils.HttpContextPtrType {
			tags := field.Tag.Get("route")
			if tags != "" {
				ret = append(ret, tags)
			}
			return ret, true
		}
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			subRet, found := h.ExtractTags(fieldType, visited)
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
