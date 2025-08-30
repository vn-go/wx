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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{}, info.routeTags)
	assert.Equal(t, "test-controller/hello", info.uriHandler)
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test\\-controller\\/hello", info.regexUri)
	assert.Equal(t, false, info.isRegexHandler)
	assert.Equal(t, -1, info.indexOfArgIsRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"method:get"}, info.routeTags)
	assert.Equal(t, "test-controller/hello2", info.uriHandler)
	assert.Equal(t, "GET", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test\\-controller\\/hello2", info.regexUri)
	assert.Equal(t, false, info.isRegexHandler)
	assert.Equal(t, -1, info.indexOfArgIsRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"@/files;method:get"}, info.routeTags)
	assert.Equal(t, "test-controller/hello3/files", info.uriHandler)
	assert.Equal(t, "GET", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test\\-controller\\/hello3\\/files", info.regexUri)
	assert.Equal(t, false, info.isRegexHandler)
	assert.Equal(t, -1, info.indexOfArgIsRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"@/files/{Path};method:get"}, info.routeTags)
	assert.Equal(t, "test-controller/hello4/files/", info.uriHandler)
	assert.Equal(t, "GET", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "^test\\-controller/hello4/files/([^/]+)$", info.regexUri)
	assert.Equal(t, true, info.isRegexHandler)
	assert.Equal(t, -1, info.indexOfArgIsRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"@/files/{*Path};method:get"}, info.routeTags)
	assert.Equal(t, "test-controller/hello5/files/", info.uriHandler)
	assert.Equal(t, "GET", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test-controller/hello5/files/(.*)", info.regexUri)
	assert.Equal(t, true, info.isRegexHandler)
	assert.Equal(t, -1, info.indexOfArgIsRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{}, info.routeTags)
	assert.Equal(t, "test-controller/json-body", info.uriHandler)
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test\\-controller\\/json\\-body", info.regexUri)
	assert.Equal(t, false, info.isRegexHandler)
	assert.Equal(t, 2, info.indexOfArgIsRequestBody)
	assert.Equal(t, reflect.TypeOf(JsonBody{}), info.typeOfRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"method:post"}, info.routeTags)
	assert.Equal(t, "test-controller/json-body2", info.uriHandler)
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test\\-controller\\/json\\-body2", info.regexUri)
	assert.Equal(t, false, info.isRegexHandler)
	assert.Equal(t, 2, info.indexOfArgIsRequestBody)

	assert.Equal(t, reflect.TypeOf(JsonBody{}), info.typeOfRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"@/{Tenant};method:post"}, info.routeTags)
	assert.Equal(t, "test-controller/json-body3/", info.uriHandler)
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "^test\\-controller/json\\-body3/([^/]+)$", info.regexUri)
	assert.Equal(t, true, info.isRegexHandler)
	assert.Equal(t, 2, info.indexOfArgIsRequestBody)

	assert.Equal(t, reflect.TypeOf(JsonBody{}), info.typeOfRequestBody)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{}, info.routeTags)
	assert.Equal(t, "test-controller/file-body", info.uriHandler)
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test\\-controller\\/file\\-body", info.regexUri)
	assert.Equal(t, false, info.isRegexHandler)
	assert.Equal(t, 2, info.indexOfArgIsRequestBody)
	assert.Equal(t, reflect.TypeOf(FileUploadBody{}), info.typeOfRequestBody)

	assert.Equal(t, true, info.isFormPost)
	assert.Equal(t, []int{2}, info.listOfIndexFieldIsFormUploadFile)

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
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"@/{Tenant}/files/{*Path}?name={Name}&age={Age};method:post"}, info.routeTags)
	assert.Equal(t, "test-controller/file-body-uri-params/", info.uriHandler)
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test-controller/file-body-uri-params/([^/]+)/files/(.*)", info.regexUri)
	assert.Equal(t, true, info.isRegexHandler)
	assert.Equal(t, 2, info.indexOfArgIsRequestBody)
	assert.Equal(t, reflect.TypeOf(FileUploadBody{}), info.typeOfRequestBody)

	assert.Equal(t, true, info.isFormPost)
	assert.Equal(t, []int{2}, info.listOfIndexFieldIsFormUploadFile)
	assert.Equal(t, "Tenant", info.uriParams[0].Name)
	assert.Equal(t, []int{1}, info.uriParams[0].FieldIndex)
	assert.Equal(t, "Path", info.uriParams[1].Name)
	assert.Equal(t, []int{2}, info.uriParams[1].FieldIndex)
	assert.Equal(t, true, info.isQueryUri)
	assert.Equal(t, "Name", info.queryParams[0].Name)
	assert.Equal(t, []int{3}, info.queryParams[0].FieldIndex)
	assert.Equal(t, "Age", info.queryParams[1].Name)
	assert.Equal(t, []int{4}, info.queryParams[1].FieldIndex)

}
func (tst *TestController) FileBodySimple(ctx *struct {
	HttpContext `route:"@/{Tenant}/files/{*Path}?name={Name}&age={Age};method:post"`
	Tenant      string
	Path        string
	Name        string
	Age         int
}, file multipart.FileHeader) {

}
func TestGetHandlerInfo_TestController_FileBodySimple(t *testing.T) {
	mt, ok := utils.GetMethodByName(reflect.TypeOf(TestController{}), "FileBodySimple")
	assert.True(t, ok)
	info, err := utils.GetHandlerInfo(mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	utils.ExtractUriInfo(info)
	utils.ExtractBodyInfo(info)
	assert.Equal(t, []string{"@/{Tenant}/files/{*Path}?name={Name}&age={Age};method:post"}, info.routeTags)
	assert.Equal(t, "test-controller/file-body-simple/", info.uriHandler)
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, false, info.isAbsUri)
	assert.Equal(t, info.isAbsUri, false)
	assert.Equal(t, "test-controller/file-body-simple/([^/]+)/files/(.*)", info.regexUri)
	assert.Equal(t, true, info.isRegexHandler)
	assert.Equal(t, 2, info.indexOfArgIsRequestBody)
	assert.Equal(t, reflect.TypeOf(multipart.FileHeader{}), info.typeOfRequestBody)

	assert.Equal(t, true, info.isFormPost)
	assert.Equal(t, []int([]int(nil)), info.listOfIndexFieldIsFormUploadFile)
	assert.Equal(t, "Tenant", info.uriParams[0].Name)
	assert.Equal(t, []int{1}, info.uriParams[0].FieldIndex)
	assert.Equal(t, "Path", info.uriParams[1].Name)
	assert.Equal(t, []int{2}, info.uriParams[1].FieldIndex)
	assert.Equal(t, true, info.isQueryUri)
	assert.Equal(t, "Name", info.queryParams[0].Name)
	assert.Equal(t, []int{3}, info.queryParams[0].FieldIndex)
	assert.Equal(t, "Age", info.queryParams[1].Name)
	assert.Equal(t, []int{4}, info.queryParams[1].FieldIndex)

}
