package wx

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type checkTypeEnum int

const (
	rangeType checkTypeEnum = iota
	regexType
)

type checks struct {
	checkType checkTypeEnum
	minValue  reflect.Value
	maxValue  reflect.Value
	re        *regexp.Regexp
	message   string
}
type checkInfo struct {
	chks  []checks
	field reflect.StructField
}

type validators struct{}

type initInitCheck struct {
	val  []checkInfo
	once sync.Once
}

var initInitCheckCache sync.Map

func (v *validators) initCheck(typ reflect.Type, visited map[reflect.Type]bool) []checkInfo {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	if visited[typ] {
		return nil
	}
	visited[typ] = true

	var ret []checkInfo
	for i := 0; i < typ.NumField(); i++ {
		var chks []checks
		ft := typ.Field(i).Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
			if ft.Kind() == reflect.Struct {
				ret = append(ret, v.initCheck(ft, visited)...)
			}
			continue
		}

		tag := typ.Field(i).Tag.Get("check")
		if tag == "" {
			continue
		}
		tag = ";" + tag + ";"

		// parse range
		if strings.Contains(tag, ";range:(") {
			txtTags := strings.Split(strings.Split(tag, "(")[1], ")")[0]
			items := strings.Split(txtTags, ",")
			if len(items) == 2 {
				// int range
				var intMin, intMax *int
				if val, err := strconv.Atoi(items[0]); err == nil {
					intMin = &val
				}
				if val, err := strconv.Atoi(items[1]); err == nil {
					intMax = &val
				}
				if intMin != nil || intMax != nil {
					chk := checks{checkType: rangeType}
					if intMin != nil {
						chk.minValue = reflect.ValueOf(*intMin)
					}
					if intMax != nil {
						chk.maxValue = reflect.ValueOf(*intMax)
					}
					chks = append(chks, chk)
				}

				// time range
				var tMin, tMax *time.Time
				if m, err := time.Parse("2006-01-02 15:04:05Z", items[0]); err == nil {
					tMin = &m
				}
				if m, err := time.Parse("2006-01-02 15:04:05Z", items[1]); err == nil {
					tMax = &m
				}
				if tMin != nil || tMax != nil {
					chk := checks{checkType: rangeType}
					if tMin != nil {
						chk.minValue = reflect.ValueOf(*tMin)
					}
					if tMax != nil {
						chk.maxValue = reflect.ValueOf(*tMax)
					}
					chks = append(chks, chk)
				}
			}
		}

		// parse regex
		if strings.Contains(tag, ";regex:") {
			reTxt := strings.Split(strings.Split(tag, ";regex:")[1], ";")[0]
			if re, err := regexp.Compile(reTxt); err == nil {
				chks = append(chks, checks{
					checkType: regexType,
					re:        re,
				})
			}
		}

		// parse custom error
		if strings.Contains(tag, ";error:") {
			msg := strings.Split(strings.Split(tag, ";error:")[1], ";")[0]
			for i := range chks {
				chks[i].message = msg
			}
		}

		if len(chks) > 0 {
			ret = append(ret, checkInfo{
				chks:  chks,
				field: typ.Field(i),
			})
		}
	}
	return ret
}

func (v *validators) init(typ reflect.Type) []checkInfo {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	a, _ := initInitCheckCache.LoadOrStore(typ, &initInitCheck{})
	i := a.(*initInitCheck)
	i.once.Do(func() {
		i.val = v.initCheck(typ, make(map[reflect.Type]bool))
	})
	return i.val
}

func (v *validators) compareInts(a, b reflect.Value) bool {
	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() < b.Uint()
	default:
		return false
	}
}

func (v *validators) Check(data any) []string {
	if data == nil {
		return nil
	}
	t := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}

	info := v.init(t)
	var ret []string
	for _, i := range info {
		f := val.FieldByIndex(i.field.Index)

		for _, c := range i.chks {
			switch c.checkType {
			case rangeType:
				// time.Time
				if i.field.Type == reflect.TypeOf(time.Time{}) {
					if c.minValue.IsValid() {
						if !c.minValue.Interface().(time.Time).Before(f.Interface().(time.Time)) {
							msg := c.message
							if msg == "" {
								msg = fmt.Sprintf("'%s' must be after %v", i.field.Name, c.minValue.Interface())
							}
							ret = append(ret, msg)
						}
					}

					if c.maxValue.IsValid() {
						if !f.Interface().(time.Time).Before(c.maxValue.Interface().(time.Time)) {
							msg := c.message
							if msg == "" {
								msg = fmt.Sprintf("'%s' must be before %v", i.field.Name, c.maxValue.Interface())
							}
							ret = append(ret, msg)
						}
					}
				} else if i.field.Type.Kind() == reflect.String {
					// string length
					length := reflect.ValueOf(f.Len())
					if c.minValue.IsValid() && !v.compareInts(c.minValue, length) {
						msg := c.message
						if msg == "" {
							msg = fmt.Sprintf("'%s' length must be >= %v", i.field.Name, c.minValue.Interface())
						}
						ret = append(ret, msg)
					}

					if c.maxValue.IsValid() && !v.compareInts(length, c.maxValue) {
						msg := c.message
						if msg == "" {
							msg = fmt.Sprintf("'%s' length must be <= %v", i.field.Name, c.maxValue.Interface())
						}
						ret = append(ret, msg)
					}
				} else {
					// number
					if c.minValue.IsValid() && !v.compareInts(c.minValue, f) {
						msg := c.message
						if msg == "" {
							msg = fmt.Sprintf("'%s' must be >= %v", i.field.Name, c.minValue.Interface())
						}
						ret = append(ret, msg)
					}
					if c.maxValue.IsValid() && !v.compareInts(f, c.maxValue) {
						msg := c.message
						if msg == "" {
							msg = fmt.Sprintf("'%s' must be <= %v", i.field.Name, c.maxValue.Interface())
						}
						ret = append(ret, msg)
					}
				}
			case regexType:
				if !c.re.MatchString(f.String()) {
					msg := c.message
					if msg == "" {
						msg = fmt.Sprintf("'%s' does not match regex %s", i.field.Name, c.re.String())
					}
					ret = append(ret, msg)
				}
			}
		}
	}
	return ret
}

var Validators = &validators{}
