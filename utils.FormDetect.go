package wx

import (
	"reflect"
	"sync"
)

type formDetectType struct {
	/*
		multipart.FileHeader
	*/
	fileHeaderType reflect.Type
	/*
	*multipart.FileHeader
	 */
	fileHeaderTypePtr reflect.Type
	/*
		[]multipart.FileHeader
	*/
	fileHeaderTypes reflect.Type
	/*
	*[]multipart.FileHeader
	 */
	fileHeaderTypesPtr reflect.Type
	/*
		[]*multipart.FileHeader
	*/
	fileHeaderPtrTypes    reflect.Type
	fileHeaderPtrTypesPtr reflect.Type
}
type initFindFormUploadField struct {
	fieldIndexes [][]int
	found        bool
	once         sync.Once
}

var cacheFindFormUploadField sync.Map

func (fd *formDetectType) FindFormUploadField(typ reflect.Type) ([][]int, bool) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {

		if typ == fd.fileHeaderType ||
			typ == fd.fileHeaderTypePtr ||
			typ == fd.fileHeaderPtrTypes ||
			typ == fd.fileHeaderPtrTypesPtr {
			return nil, true
		}
		return nil, false
	}
	acctually, _ := cacheFindFormUploadField.LoadOrStore(typ.String(), &initFindFormUploadField{})
	item := acctually.(*initFindFormUploadField)
	item.once.Do(func() {
		item.fieldIndexes, item.found = fd.findFormUploadField(typ, map[reflect.Type]bool{})
	})
	return item.fieldIndexes, item.found
}
func (fd *formDetectType) findFormUploadField(typ reflect.Type, visted map[reflect.Type]bool) ([][]int, bool) {
	ret := [][]int{}
	var found bool
	if visted[typ] {
		return nil, false
	}
	visted[typ] = true
	if typ == fd.fileHeaderType ||
		typ == fd.fileHeaderTypePtr ||
		typ == fd.fileHeaderTypes ||
		typ == fd.fileHeaderTypesPtr ||
		typ == fd.fileHeaderPtrTypes ||
		typ == fd.fileHeaderPtrTypesPtr {
		return nil, true
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, false
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if subIdex, ok := fd.findFormUploadField(fieldType, visted); ok {
			if subIdex == nil {
				ret = append(ret, []int{i})

			} else {
				for _, v := range subIdex {
					ret = append(ret, append([]int{i}, v...))
				}

			}
			found = true
		}

	}
	if !found {
		return nil, false
	}
	return ret, found

}
