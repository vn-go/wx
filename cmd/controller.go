package main

import (
	"mime/multipart"

	"github.com/vn-go/wx"
)

type User struct {
	Username string
}

/*
 */
type Media struct {
	wx.Authenticate[User]
}

func (m *Media) Upload(ctx *wx.Handler, file multipart.FileHeader) (string, error) {
	return "heelo", nil
}
func (m *Media) Auth(ctx *wx.Handler, data wx.Form[struct {
}]) (string, error) {
	return "heelo", nil
}
func (media *Media) Files(ctx *struct {
	wx.Handler `route:"@/{*FilePath};method:get"`
	FilePath   string
}) error {
	// uriOfViewFileHandler, err := wx.GetUriOfHandler[Media]("Files")
	// if err != nil {
	// 	return "", err
	// }
	// rootUrl := ctx.Handler().GetAbsRootUri() + uriOfViewFileHandler
	// filePath := media.RootDir + "/" + ctx.FilePath
	return ctx.Handler().StreamingFile(`D:\code\go\wx\wx\cmd\controller.go`)
	//return ctx.FilePath, nil
}
func (media *Media) Files2(ctx *struct {
	wx.Handler `route:"@/{FilePath}/{filename};method:get"`
	FilePath   string
}) error {
	// uriOfViewFileHandler, err := wx.GetUriOfHandler[Media]("Files")
	// if err != nil {
	// 	return "", err
	// }
	// rootUrl := ctx.Handler().GetAbsRootUri() + uriOfViewFileHandler
	// filePath := media.RootDir + "/" + ctx.FilePath
	return ctx.Handler().StreamingFile(`D:\code\go\wx\wx\cmd\controller.go`)
	//return ctx.FilePath, nil
}
func (media *Media) Files3(ctx *struct {
	wx.Handler `route:"@/{FilePath}/{filename}?test={code};method:get"`
	FilePath   string
}) error {
	// uriOfViewFileHandler, err := wx.GetUriOfHandler[Media]("Files")
	// if err != nil {
	// 	return "", err
	// }
	// rootUrl := ctx.Handler().GetAbsRootUri() + uriOfViewFileHandler
	// filePath := media.RootDir + "/" + ctx.FilePath
	return ctx.Handler().StreamingFile(`D:\code\go\wx\wx\cmd\controller.go`)
	//return ctx.FilePath, nil
}
func init() {
	(&wx.Authenticate[User]{}).Verify(func(ctx wx.Handler) (*User, error) {

		authHeader := ctx().Req.Header.Get("Authorization")

		if authHeader == "" {
			return nil, wx.Errors.NewUnauthorizedError()
		}

		return &User{}, nil
	})
}
