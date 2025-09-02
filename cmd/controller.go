package main

import (
	"mime/multipart"

	"github.com/vn-go/wx"
)

type User struct {
	Username string
}
type Media struct {
	wx.Authenticate[User]
}

func (m *Media) Upload(ctx *wx.Handler, file multipart.FileHeader) (string, error) {
	return "heelo", nil
}
func init() {
	(&wx.Authenticate[User]{}).Verify(func(ctx wx.Handler) (*User, error) {

		authHeader := ctx().Req.Header.Get("Authorization")
		if authHeader == "" {
			return nil, wx.NewUnauthorizedError()
		}

		return &User{}, nil
	})
}
