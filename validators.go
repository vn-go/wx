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
	opMax     string
	minValue  reflect.Value
	opMin     string
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
type InfoCheck struct {
	Message       string `json:"message"`
	MessageLayout string `json:"messageLayout"`
	FieldName     string `json:"fieldName"`
	MinValue      any    `json:"minValue"`
	MaxValue      any    `json:"maxValue"`
	Value         any    `json:"value"`
}

var initInitCheckCache sync.Map

func (v *validators) revertOp(op string) string {
	switch op {
	case "<":
		return ">="
	case "<=":
		return ">"
	case ">":
		return "<="
	case ">=":
		return "<"
	}
	return ""
}
func (v *validators) compareFloat64(a, b float64, op string) bool {
	if op == "<" {
		return a < b
	}
	if op == "<=" {
		return a <= b
	}
	if op == ">" {
		return a > b
	}
	if op == ">=" {
		return a >= b
	}
	return false
}
func (v *validators) compareInt64(a, b int64, op string) bool {
	if op == "<" {
		return a < b
	}
	if op == "<=" {
		return a <= b
	}
	if op == ">" {
		return a > b
	}
	if op == ">=" {
		return a >= b
	}
	return false
}
func (v *validators) compareUint64(a, b uint64, op string) bool {
	if op == "<" {
		return a < b
	}
	if op == "<=" {
		return a <= b
	}
	if op == ">" {
		return a > b
	}
	if op == ">=" {
		return a >= b
	}
	return false
}
func (v *validators) compareTime(a, b time.Time, op string) bool {
	switch op {
	case "<":
		return a.Before(b)
	case "<=":
		return a.Before(b) || a.Equal(b)
	case ">":
		return a.After(b)
	case ">=":
		return a.After(b) || a.Equal(b)
	case "==":
		return a.Equal(b)
	case "!=":
		return !a.Equal(b)
	default:
		return false
	}
}
func (v *validators) getValue(ft reflect.Type, val string) reflect.Value {
	if ft == reflect.TypeFor[time.Time]() {
		if m, err := time.Parse("2006-01-02 15:04:05Z", val); err == nil {
			return reflect.ValueOf(m)
		}
	}
	switch ft.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if r, err := strconv.Atoi(val); err == nil {
			return reflect.ValueOf(r)
		}
	case reflect.Float32:
		if r, err := strconv.ParseFloat(val, 32); err == nil {
			return reflect.ValueOf(float32(r))
		}

	case reflect.Float64:
		if r, err := strconv.ParseFloat(val, 64); err == nil {
			return reflect.ValueOf(r)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if r, err := strconv.Atoi(val); err == nil {
			return reflect.ValueOf(r)
		}
	case reflect.String:
		if r, err := strconv.Atoi(val); err == nil {
			return reflect.ValueOf(r)
		}

	default:
		var retVale reflect.Value
		return retVale
	}
	var retVale reflect.Value
	return retVale
}
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
				continue
			}

		}

		tag := typ.Field(i).Tag.Get("check")
		if tag == "" {
			continue
		}
		tag = ";" + tag + ";"

		// parse range
		if strings.Contains(tag, ";range:") {
			txtTags := strings.Split(strings.Split(tag, ";range:")[1], ";")[0]
			if strings.Contains(txtTags, ":") {
				items := strings.Split(txtTags, ":")
				opMax := ">="
				if items[1][len(items[1])-1] == ']' {
					opMax = ">"
				}
				opMin := "<="
				if items[0][0] == '[' {
					opMin = "<"
				}
				chks = append(chks, checks{
					opMax:     opMax,
					opMin:     opMin,
					checkType: rangeType,
					minValue:  v.getValue(ft, items[0][1:]),
					maxValue:  v.getValue(ft, items[1][0:len(items[1])-1]),
				})
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
func (v *validators) checkIsInRangeTime(fieldName string, val, min, max reflect.Value, opMin, opMax string) *InfoCheck {
	dVal := val.Interface().(time.Time)
	if min.IsValid() && max.IsValid() {
		from := min.Interface().(time.Time)
		to := max.Interface().(time.Time)
		if !(v.compareTime(from, dVal, opMin) && v.compareTime(dVal, to, opMax)) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s %s and %s %s", fieldName, v.revertOp(opMin), from, v.revertOp(opMax), to),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' be %s {MinValue} and %s {MaxValue}", v.revertOp(opMin), v.revertOp(opMax)),
				MinValue:      from,
				MaxValue:      to,
				Value:         dVal,
			}
		}

	}
	if min.IsValid() && !max.IsValid() {
		from := min.Interface().(time.Time)
		if !v.compareTime(from, dVal, opMin) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s  %s", fieldName, v.revertOp(opMin), from),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MinValue}", v.revertOp(opMin)),
				MinValue:      from,
				Value:         dVal,
			}
		}

	}
	if !min.IsValid() && max.IsValid() {
		to := max.Interface().(time.Time)
		if !v.compareTime(dVal, to, opMax) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s  %s", fieldName, v.revertOp(opMax), to),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MaxValue}", v.revertOp(opMax)),
				MaxValue:      to,
				Value:         dVal,
			}
		}

	}
	return nil
}
func (v *validators) checkIsInRangeInt(fieldName string, val, min, max reflect.Value, opMin, opMax string) *InfoCheck {
	if min.IsValid() && max.IsValid() {
		if !(v.compareInt64(min.Int(), val.Int(), opMin) && v.compareInt64(val.Int(), max.Int(), opMax)) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of  '%s' must be %s %d and %s %d", fieldName, v.revertOp(opMin), min.Int(), v.revertOp(opMax), max.Int()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MinValue} and %s {MaxValue}", v.revertOp(opMin), v.revertOp(opMax)),
				MinValue:      min.Int(),
				MaxValue:      max.Int(),
				Value:         val.Int(),
			}
		}
	}
	if min.IsValid() && !max.IsValid() {
		if !v.compareInt64(min.Int(), val.Int(), opMin) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s %d", fieldName, v.revertOp(opMin), min.Int()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MaxValue}", v.revertOp(opMin)),
				MinValue:      min.Int(),
				Value:         val.Int(),
			}
		}
	}
	if !min.IsValid() && max.IsValid() {
		if !v.compareInt64(val.Int(), max.Int(), opMax) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s %d", fieldName, v.revertOp(opMax), max.Int()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MaxValue}", v.revertOp(opMax)),

				MaxValue: max.Int(),
				Value:    val.Int(),
			}
		}
	}
	return nil
}
func (v *validators) checkIsInRangeUInt(fieldName string, val, min, max reflect.Value, opMin, opMax string) *InfoCheck {
	if min.IsValid() && max.IsValid() {
		if !(v.compareUint64(min.Uint(), val.Uint(), opMin) && v.compareUint64(val.Uint(), max.Uint(), opMax)) {
			return &InfoCheck{
				Message:       fmt.Sprintf("value of '%s' must be %s %d and %s %d", fieldName, v.revertOp(opMin), min.Int(), v.revertOp(opMax), max.Int()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MinValue} and %s {MaxValue}", v.revertOp(opMin), v.revertOp(opMax)),
				MinValue:      min.Uint(),
				MaxValue:      max.Uint(),
				Value:         val.Uint(),
			}
		}
	}
	if min.IsValid() && !max.IsValid() {
		if !v.compareUint64(min.Uint(), val.Uint(), opMin) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s %d", fieldName, v.revertOp(opMin), min.Uint()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MinValue}", v.revertOp(opMin)),
				MinValue:      min.Uint(),
				Value:         val.Uint(),
			}
		}
	}
	if !min.IsValid() && max.IsValid() {
		if !v.compareUint64(val.Uint(), max.Uint(), opMax) {
			return &InfoCheck{
				Message:       fmt.Sprintf("'%s' must be %s %d", fieldName, v.revertOp(opMax), max.Uint()),
				MessageLayout: fmt.Sprintf("'{FieldName}' must be %s {MaxValue}", v.revertOp(opMax)),

				MaxValue: max.Uint(),
				Value:    val.Uint(),
			}
		}
	}
	return nil
}

func (v *validators) checkIsInRangeFloat(fieldName string, val, min, max reflect.Value, opMin, opMax string) *InfoCheck {
	if min.IsValid() && max.IsValid() {
		if !v.compareFloat64(min.Float(), val.Float(), opMin) && v.compareFloat64(val.Float(), max.Float(), opMax) {
			return &InfoCheck{
				Message:       fmt.Sprintf("value of '%s' must be %s %f and %s %f", fieldName, v.revertOp(opMin), min.Float(), v.revertOp(opMax), max.Float()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s  {MinValue} and %s {MaxValue}", v.revertOp(opMin), v.revertOp(opMax)),
				MinValue:      min.Float(),
				MaxValue:      max.Float(),
				Value:         val.Float(),
			}
		}
	}
	if min.IsValid() && !max.IsValid() {
		if !v.compareFloat64(min.Float(), val.Float(), opMin) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s %f", fieldName, v.revertOp(opMin), min.Float()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MinValue}", v.revertOp(opMin)),
				MinValue:      min.Float(),
				Value:         val.Float(),
			}
		}
	}
	if !min.IsValid() && max.IsValid() {
		if !v.compareFloat64(val.Float(), max.Float(), opMax) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Value of '%s' must be %s %f", fieldName, v.revertOp(opMax), max.Float()),
				MessageLayout: fmt.Sprintf("Value of '{FieldName}' must be %s {MaxValue}", v.revertOp(opMax)),

				MaxValue: max.Float(),
				Value:    val.Float(),
			}
		}
	}
	return nil
}
func (v *validators) checkLenIsInRange(fieldName string, val, min, max reflect.Value, opMin, opMax string) *InfoCheck {
	if min.IsValid() && max.IsValid() {
		if !(v.compareInt64(min.Int(), int64(val.Len()), opMin) && v.compareInt64(int64(val.Len()), max.Int(), opMax)) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Len of '%s' must be %s %d and %s %d", fieldName, v.revertOp(opMin), min.Int(), v.revertOp(opMax), max.Int()),
				MessageLayout: fmt.Sprintf("Len of '{FieldName}' must be %s  {MinValue} and %s {MaxValue}", v.revertOp(opMin), v.revertOp(opMax)),
				MinValue:      min.Int(),
				MaxValue:      max.Int(),
				Value:         val.String(),
			}
		}
	}
	if min.IsValid() && !max.IsValid() {
		if !v.compareInt64(min.Int(), int64(val.Len()), opMin) {
			return &InfoCheck{
				Message:       fmt.Sprintf("Len of '%s' must be %s %d", fieldName, v.revertOp(opMin), min.Int()),
				MessageLayout: fmt.Sprintf("Len of  '{FieldName}' must be %s {MinValue}", v.revertOp(opMin)),
				MinValue:      min.Int(),
				Value:         val.String(),
			}
		}
	}
	if !min.IsValid() && max.IsValid() {
		if !v.compareInt64(int64(val.Len()), max.Int(), opMax) {
			return &InfoCheck{
				Message:       fmt.Sprintf("'%s' must be %s %d", fieldName, v.revertOp(opMax), max.Int()),
				MessageLayout: fmt.Sprintf("Len of  '{FieldName}' must be %s {MaxValue}", v.revertOp(opMax)),

				MaxValue: max.Int(),
				Value:    val.String(),
			}
		}
	}
	return nil
}
func (v *validators) checkIsInRange(fieldName string, val, min, max reflect.Value, opMin, opMax string) *InfoCheck {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Type() == reflect.TypeFor[time.Time]() {
		return v.checkIsInRangeTime(fieldName, val, min, max, opMin, opMax)
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.checkIsInRangeInt(fieldName, val, min, max, opMin, opMax)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.checkIsInRangeUInt(fieldName, val, min, max, opMin, opMax)
	case reflect.Float32, reflect.Float64:
		return v.checkIsInRangeFloat(fieldName, val, min, max, opMin, opMax)
	case reflect.String:
		return v.checkLenIsInRange(fieldName, val, min, max, opMin, opMax)
	default:
		return nil
	}

}
func (v *validators) compareInts(a, b reflect.Value) bool {
	if b.Kind() == reflect.Ptr {
		if b.IsNil() {
			return true
		}
		b = b.Elem()
	}
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}

	switch b.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int()
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() < b.Uint()
	default:
		return false
	}
}
func (v *validators) CheckValue(val reflect.Value) []InfoCheck {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil
	}
	t := val.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}

	info := v.init(t)
	ret := []InfoCheck{}
	for _, i := range info {
		f := val.FieldByIndex(i.field.Index)

		for _, c := range i.chks {

			if c.checkType == rangeType {
				if item := v.checkIsInRange(i.field.Name, f, c.minValue, c.maxValue, c.opMin, c.opMax); item != nil {
					item.FieldName = i.field.Name
					ret = append(ret, *item)
				}
			}
			if c.checkType == regexType {
				if !c.re.MatchString(f.String()) {
					msg := c.message
					if msg == "" {
						msg = fmt.Sprintf("'%s' does not match regex %s", i.field.Name, c.re.String())
					}
					ret = append(ret, InfoCheck{
						Message:       msg,
						MessageLayout: "Value of {FieldName} is invalid",
						FieldName:     i.field.Name,
						Value:         f.String(),
					})
				}
			}
		}
	}
	return ret
}
func (v *validators) Check(data any) []InfoCheck {
	// if data == nil {
	// 	return nil
	// }
	// t := reflect.TypeOf(data)
	// val := reflect.ValueOf(data)
	return v.CheckValue(reflect.ValueOf(data))

}

var Validators = &validators{}
