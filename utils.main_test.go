package wx

import (
	"fmt"
	"mime/multipart"
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

var count = 0

func (ex *Example) SimpleUpload(ctx *HttpContext, file *multipart.File) {
	//count++
	//fmt.Println(file)
}
func TestSimpleUpload(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUpload")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), Mock.CreateMockFile("test", "OK"))
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
}
func BenchmarkSimpleUpload(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUpload")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), Mock.CreateMockFile("test", "OK"))
	assert.NoError(t, err)
	res := Mock.NewRes()
	t.ReportAllocs()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
	}
	fmt.Println(count)
}
func (ex *Example) SimpleUploadFiles(ctx *HttpContext, file []multipart.FileHeader) {
	//fmt.Println(file)
}
func TestSimpleUploadFiles(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUploadFiles")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	mockFiles := []mockFile{}
	for i := 0; i < 10; i++ {
		file := Mock.CreateMockFile("test", "OK")
		mockFiles = append(mockFiles, *file)
	}

	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), mockFiles)
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
}
func BenchmarkSimpleUploadFiles(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUploadFiles")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	mockFiles := []mockFile{}
	for i := 0; i < 10; i++ {
		file := Mock.CreateMockFile("test", "OK")
		mockFiles = append(mockFiles, *file)
	}

	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), mockFiles)
	assert.NoError(t, err)
	res := Mock.NewRes()
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
	}

}

type User struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (ex *Example) SimpleUploadFiles2(ctx *HttpContext, data struct {
	Files []multipart.FileHeader
	User  User
}) {

}
func TestSimpleUploadFiles2(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUploadFiles2")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	mockFiles := []mockFile{}
	for i := 0; i < 10; i++ {
		file := Mock.CreateMockFile("Files", "OK")
		mockFiles = append(mockFiles, *file)
	}
	postdata := struct {
		Files []mockFile
		User  User
	}{
		Files: mockFiles,
		User: User{
			Code: "adadas",
			Name: "dsadasdasdadad",
		},
	}
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), postdata)
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
}
func BenchmarkSimpleUploadFiles2(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUploadFiles2")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	mockFiles := []mockFile{}
	for i := 0; i < 10; i++ {
		file := Mock.CreateMockFile("Files", "OK")
		mockFiles = append(mockFiles, *file)
	}
	postdata := struct {
		Files []mockFile
		User  User
	}{
		Files: mockFiles,
		User: User{
			Code: "adadas",
			Name: "dsadasdasdadad",
		},
	}
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), postdata)
	assert.NoError(t, err)
	res := Mock.NewRes()
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
	}

}
func (ex *Example) FormPost(ctx *HttpContext, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.FormRequest("POST", handler.GetUriHandler(), User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
}
func BenchmarkFormPost(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.FormRequest("POST", handler.GetUriHandler(), User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
		//assert.Equal(t, 200, res.Code)
	}

}
func (ex *Example) FormPost2(ctx *struct {
	HttpContext `route:"@/files/{FileName}"`
	FileName    string
}, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost2(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost2")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.FormRequest("POST", handler.GetUriHandler()+"Test.txt", User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
}
func (ex *Example) FormPost3(ctx *struct {
	HttpContext `route:"@/files/{*filePath}"`
	FilePath    string
}, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost3(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost3")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.FormRequest("POST", handler.GetUriHandler()+"dasda/dsad/sad/das/Test.txt", User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
}
func (ex *Example) FormPost4(ctx *struct {
	HttpContext `route:"@/files/{*filePath}?dirPath={dirPath}"`
	FilePath    string
	DirPath     string
}, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost4(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost4")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.FormRequest("POST", handler.GetUriHandler()+"dasda/dsad/sad/das/Test.txt?dirPath=/abd/cde", User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
}
