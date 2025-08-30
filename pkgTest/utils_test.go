package pkgtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/wx"
	_ "github.com/vn-go/wx"
	"github.com/vn-go/wx/mock"
)

type ObjTest001 struct {
}

func (obj *ObjTest001) TestMethod001(ctx *wx.HttpContext) {

}
func TestUtilsFindMethod(t *testing.T) {
	mt, ok := mock.FindMethod[ObjTest001]("TestMethod001")

	assert.True(t, ok)
	assert.Equal(t, "TestMethod001", mt.Name)

}
func BenchmarkUtilsFindMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mt, ok := mock.FindMethod[ObjTest001]("TestMethod001")
		assert.True(b, ok)
		assert.Equal(b, "TestMethod001", mt.Name)
	}

}
func TestGetHandlerInfo(t *testing.T) {
	mt, ok := mock.FindMethod[ObjTest001]("TestMethod001")
	assert.True(t, ok)
	assert.Equal(t, "TestMethod001", mt.Name)
	info, err := mock.GetHandlerInfo[ObjTest001](mt)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, 1, info.IndexOfArg)
	assert.Equal(t, []int{0}, info.ReqFieldIndex)
	assert.Equal(t, []int{1}, info.ResFieldIndex)

}
