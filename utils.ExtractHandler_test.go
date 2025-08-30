package wx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example struct {
	Ctx *HttpContext
}

func (ex *Example) PostNoBodyMethod(ctx *HttpContext) {
	ex.Ctx.Res.Write([]byte("ok"))

}
func TestPostNoBodyMethod(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("PostNoBodyMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), nil)
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)

	//handler.ConrollerNewMethod

}
func BenchmarkPostNoBodyMethod(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("PostNoBodyMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), nil)
	assert.NoError(t, err)
	res := Mock.NewRes()
	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
		//assert.Equal(t, 200, res.Code)
	}

	//handler.ConrollerNewMethod

}
func (ex *Example) GetMethod(ctx *struct {
	HttpContext `route:"method:get"`
}) {
	//ex.Res.Write([]byte("ok"))

}
func TestGetMethod(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("GetMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("GET", handler.GetUriHandler(), nil)
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)

	//handler.ConrollerNewMethod

}
func BenchmarkGetMethod(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("GetMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("GET", handler.GetUriHandler(), nil)
	assert.NoError(t, err)
	res := Mock.NewRes()
	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
		//assert.Equal(t, 200, res.Code)
	}

	//handler.ConrollerNewMethod

}
func (ex *Example) PostBodyMethod(ctx *HttpContext, data *struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}) {
	ex.Ctx.Res.Write([]byte("ok"))

}
func TestPostBodyMethod(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("PostBodyMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "test",
		Age:  10,
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)

	//handler.ConrollerNewMethod

}
func BenchmarkPostBodyMethod(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("PostBodyMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "test",
		Age:  10,
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
		//assert.Equal(t, 200, res.Code)
	}

	//handler.ConrollerNewMethod

}
