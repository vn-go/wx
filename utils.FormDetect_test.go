package wx

import (
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test001 struct {
	Name string
	Age  int
}

func Test_FormDetect_Test001(t *testing.T) {
	var fx [][]int
	var ok bool
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(Test001{}))
	assert.False(t, ok)
	assert.Nil(t, fx)
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(&Test001{}))
	assert.False(t, ok)
	assert.Nil(t, fx)
}

type Test002 struct {
	Name  string
	Age   int
	Files []multipart.FileHeader
}

func Test_FormDetect_Test002(t *testing.T) {
	var fx [][]int
	var ok bool
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(Test002{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(&Test002{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
}

type Test003 struct {
	Name  string
	Age   int
	Files []*multipart.FileHeader
}

func Test_FormDetect_Test003(t *testing.T) {
	var fx [][]int
	var ok bool
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(Test003{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(&Test003{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
}

type Test004 struct {
	Name  string
	Age   int
	Files *[]multipart.FileHeader
}

func Test_FormDetect_Test004(t *testing.T) {
	var fx [][]int
	var ok bool
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(Test004{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(&Test004{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
}

type Test005 struct {
	Name  string
	Age   int
	Files *[]*multipart.FileHeader
}

func Test_FormDetect_Test005(t *testing.T) {
	var fx [][]int
	var ok bool
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(Test004{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(&Test004{}))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{2}}, fx)
}

type Test006 struct {
	Name  string
	Age   int
	Files *[]*multipart.FileHeader
}

func Test_FormDetect_Test006(t *testing.T) {
	var fx [][]int
	var ok bool
	bodyTest := Form[Test006]{}
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(bodyTest))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{0, 2}}, fx)
	bodyTest2 := &Form[Test006]{}
	fx, ok = utils.formDetect.FindFormUploadField(reflect.TypeOf(bodyTest2))
	assert.True(t, ok)
	assert.Equal(t, [][]int{{0, 2}}, fx)

}
