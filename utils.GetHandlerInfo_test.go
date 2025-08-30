package wx

import (
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestController struct {
}

func (tst *TestController) Hello(ctx *HttpContext) {

}
func TestGetHandlerInfo_TestController_Hello(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "Hello")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{}, info.RouteTags)
	assert.Equal(t, "test-controller/hello", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test\\-controller\\/hello", info.RegexUri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, -1, info.IndexOfArgIsRequestBody)

}
func (tst *TestController) Hello2(ctx *struct {
	HttpContext `route:"method:get"`
}) {

}
func TestGetHandlerInfo_TestController_Hello2(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "Hello2")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{"method:get"}, info.RouteTags)
	assert.Equal(t, "test-controller/hello2", info.UriHandler)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test\\-controller\\/hello2", info.RegexUri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, -1, info.IndexOfArgIsRequestBody)

}
func (tst *TestController) Hello3(ctx *struct {
	HttpContext `route:"@/files;method:get"`
}) {

}
func TestGetHandlerInfo_TestController_Hello3(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "Hello3")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{"@/files;method:get"}, info.RouteTags)
	assert.Equal(t, "test-controller/hello3/files", info.UriHandler)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test\\-controller\\/hello3\\/files", info.RegexUri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, -1, info.IndexOfArgIsRequestBody)

}
func (tst *TestController) Hello4(ctx *struct {
	HttpContext `route:"@/files/{Path};method:get"`
	Path        string
}) {

}
func TestGetHandlerInfo_TestController_Hello4(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "Hello4")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{"@/files/{Path};method:get"}, info.RouteTags)
	assert.Equal(t, "test-controller/hello4/files/", info.UriHandler)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "^test\\-controller/hello4/files/([^/]+)$", info.RegexUri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, -1, info.IndexOfArgIsRequestBody)

}
func (tst *TestController) Hello5(ctx *struct {
	HttpContext `route:"@/files/{*Path};method:get"`
	Path        string
}) {

}
func TestGetHandlerInfo_TestController_Hello5(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "Hello5")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{"@/files/{*Path};method:get"}, info.RouteTags)
	assert.Equal(t, "test-controller/hello5/files/", info.UriHandler)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test-controller/hello5/files/(.*)", info.RegexUri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, -1, info.IndexOfArgIsRequestBody)

}

type JsonBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (tst *TestController) JsonBody(ctx *HttpContext, body JsonBody) {

}
func TestGetHandlerInfo_TestController_JsonBody(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "JsonBody")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{}, info.RouteTags)
	assert.Equal(t, "test-controller/json-body", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test\\-controller\\/json\\-body", info.RegexUri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, 2, info.IndexOfArgIsRequestBody)
	assert.Equal(t, reflect.TypeOf(JsonBody{}), info.TypeOfRequestBody)

}
func (tst *TestController) JsonBody2(ctx *struct {
	HttpContext `route:"method:post"`
}, body JsonBody) {

}
func TestGetHandlerInfo_TestController_JsonBody2(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "JsonBody2")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{"method:post"}, info.RouteTags)
	assert.Equal(t, "test-controller/json-body2", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test\\-controller\\/json\\-body2", info.RegexUri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, 2, info.IndexOfArgIsRequestBody)

	assert.Equal(t, reflect.TypeOf(JsonBody{}), info.TypeOfRequestBody)

}
func (tst *TestController) JsonBody3(ctx *struct {
	HttpContext `route:"@/{Tenant};method:post"`
	Tenant      string
}, body JsonBody) {

}
func TestGetHandlerInfo_TestController_JsonBody3(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "JsonBody3")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{"@/{Tenant};method:post"}, info.RouteTags)
	assert.Equal(t, "test-controller/json-body3/", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "^test\\-controller/json\\-body3/([^/]+)$", info.RegexUri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, 2, info.IndexOfArgIsRequestBody)

	assert.Equal(t, reflect.TypeOf(JsonBody{}), info.TypeOfRequestBody)

}

type FileUploadBody struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Files multipart.FileHeader
}

func (tst *TestController) FileBody(ctx *HttpContext, body FileUploadBody) {

}
func TestGetHandlerInfo_TestController_FileBody(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "FileBody")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{}, info.RouteTags)
	assert.Equal(t, "test-controller/file-body", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test\\-controller\\/file\\-body", info.RegexUri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, 2, info.IndexOfArgIsRequestBody)
	assert.Equal(t, reflect.TypeOf(FileUploadBody{}), info.TypeOfRequestBody)

	assert.Equal(t, true, info.IsFormPost)
	assert.Equal(t, []int{2}, info.ListOfIndexFieldIsFormUploadFile)

}
func (tst *TestController) FileBodyUriParams(ctx *struct {
	HttpContext `route:"@/{Tenant}/files/{*Path}?name={Name}&age={Age};method:post"`
	Tenant      string
	Path        string
	Name        string
	Age         int
}, body FileUploadBody) {

}
func TestGetHandlerInfo_TestController_FileBodyUriParams(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "FileBodyUriParams")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodInfo(info)
	assert.Equal(t, []string{"@/{Tenant}/files/{*Path}?name={Name}&age={Age};method:post"}, info.RouteTags)
	assert.Equal(t, "test-controller/file-body-uri-params/", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, info.IsAbsUri, false)
	assert.Equal(t, "test-controller/file-body-uri-params/([^/]+)/files/(.*)", info.RegexUri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, 2, info.IndexOfArgIsRequestBody)
	assert.Equal(t, reflect.TypeOf(FileUploadBody{}), info.TypeOfRequestBody)

	assert.Equal(t, true, info.IsFormPost)
	assert.Equal(t, []int{2}, info.ListOfIndexFieldIsFormUploadFile)
	assert.Equal(t, "Tenant", info.UriParams[0].Name)
	assert.Equal(t, []int{1}, info.UriParams[0].FieldIndex)
	assert.Equal(t, "Path", info.UriParams[1].Name)
	assert.Equal(t, []int{2}, info.UriParams[1].FieldIndex)
	assert.Equal(t, true, info.IsQueryUri)
	assert.Equal(t, "Name", info.QueryParams[0].Name)
	assert.Equal(t, []int{3}, info.QueryParams[0].FieldIndex)
	assert.Equal(t, "Age", info.QueryParams[1].Name)
	assert.Equal(t, []int{4}, info.QueryParams[1].FieldIndex)

}
