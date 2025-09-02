package wx

import (
	"fmt"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example struct {
	Ctx *Handler
}

func (ex *Example) PostNoBodyMethod(ctx *Handler) {

	fmt.Println(ex.Ctx)

}
func TestPostNoBodyMethod(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("PostNoBodyMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
	Handler `route:"method:get"`
}) {
	//ex.Res.Write([]byte("ok"))

}
func TestGetMethod(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("GetMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
func (ex *Example) PostBodyMethod(ctx Handler, data *struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}) {
	(*ex.Ctx)().Res.Write([]byte("OK"))
	//ex.Ctx().Res.Write([]byte("OK"))
}
func TestPostBodyMethod(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("PostBodyMethod")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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

func (ex *Example) SimpleUpload(ctx Handler, file *multipart.FileHeader) {
	//count++
	//fmt.Println(file)
}
func TestSimpleUpload(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUpload")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
func (ex *Example) SimpleUploadFiles(ctx Handler, file []multipart.FileHeader) {
	//fmt.Println(file)
}
func TestSimpleUploadFiles(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUploadFiles")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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

func (ex *Example) SimpleUploadFiles2(handler Handler, data struct {
	Files []multipart.FileHeader
	User  User
}) {
	ctx := handler()
	ctx.Res.Write([]byte{})
}
func TestSimpleUploadFiles2(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("SimpleUploadFiles2")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
func (ex *Example) FormPost(ctx Handler, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
	Handler  `route:"@/files/{FileName}"`
	FileName string
}, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost2(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost2")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
	Handler  `route:"@/files/{*filePath}"`
	FilePath string
}, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost3(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost3")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
func BenchmarkFormPost3(t *testing.B) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost3")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := Mock.FormRequest("POST", handler.GetUriHandler()+"dasda/dsad/sad/das/Test.txt", User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := Mock.NewRes()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
	}

	//assert.Equal(t, 200, res.Code)
}
func (ex *Example) FormPost4(ctx *struct {
	Handler  `route:"@/files/{*filePath}?dirPath={dirPath}"`
	FilePath string
	DirPath  string
}, data *Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost4(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost4")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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
func (ex *Example) FormPost5(ctx *struct {
	Handler  `route:"@/files/{*filePath}?dirPath={dirPath}"`
	FilePath string
	DirPath  string
}, data *Form[*multipart.FileHeader]) (*User, error) {

	return nil, nil
}
func TestFormPost5(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example]("FormPost5")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
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

type Example2 struct {
	Handler
}

func (ex *Example2) New() error {
	//ex.Handler().Res.Write([]byte{})
	return nil
}
func (ex *Example2) Post(h *Handler) {

}
func TestExample2Post(t *testing.T) {
	handler, err := MakeHandlerFromMethod[Example2]("Post")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, false, handler.isAuth)
	assert.Nil(t, handler.fieldIndexOfAuth)
	assert.Equal(t, -1, handler.indexOfArgIsAuth)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), nil)
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
}

type OK[T any] struct {
	V T
}
type AuthExample struct {
	Handler
	OAuth2[User]
	FX *AuthExample
}

func (ex *AuthExample) New() error {
	ex.Verify(func(ctx *httpContext) (*User, error) {
		return &User{}, nil
	})
	return nil
}
func (ex *AuthExample) Post(ctx OK[Handler]) {
	fmt.Println(ctx)

}

type User001 struct {
}

func TestAuthFind(t *testing.T) {

	handler, err := MakeHandlerFromMethod[AuthExample]("Post")
	assert.NoError(t, err)
	assert.Equal(t, true, handler.isAuth)
	assert.Equal(t, []int{1}, handler.fieldIndexOfAuth)
	assert.Equal(t, 0, handler.indexOfArgIsAuth)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)
	req, err := Mock.JsonRequest("POST", handler.GetUriHandler(), nil)
	assert.NoError(t, err)
	res := Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
}
