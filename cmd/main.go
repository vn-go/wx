package main

import (
	"mime/multipart"

	"github.com/vn-go/wx"
)

type Oauth struct {
}
type LoginForm struct {
	Username string `check:"range:[4:20]"`
	Password string `check:"range:[4:20]"`
}

func (auth *Oauth) Login2(ctx wx.Handler, body LoginForm) (any, error) {
	if body.Username == "admin" && body.Password == "admin" {
		return &struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
		}{
			AccessToken: "123456",
			TokenType:   "bearer",
		}, nil
	}
	return nil, wx.Errors.NewUnauthorizedError()
}
func (auth *Oauth) Login(ctx wx.Handler, body wx.Form[LoginForm]) (any, error) {
	if body.Data.Username == "admin" && body.Data.Password == "admin" {
		return &struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
		}{
			AccessToken: "123456",
			TokenType:   "bearer",
		}, nil
	}
	return nil, wx.Errors.NewUnauthorizedError()
}
func (auth *Oauth) Upload(ctx wx.Handler, body struct {
	File multipart.FileHeader
	Data struct {
		Code string `check:"range:[3:10]"`
	}
	Name string `check:"range:[3:10]"`
}) (any, error) {
	return body.Data, nil
}
func main() {
	wx.Options.UsePool = true
	if err := wx.Routes("/api", &Media{}, &Oauth{}); err != nil {
		panic(err)
	}
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/swagger")
	swagger.Info(wx.SwaggerInfo{
		Title:       "Swagger Example API",
		Description: "This is a sample server Petstore server.",
		Version:     "1.0.0",
	})
	swagger.OAuth2Password("/api/oauth/login")
	swagger.Build()

	server.Middleware(wx.MiddlWares.Cors)
	server.Middleware(wx.MiddlWares.Zip)
	server.Start()
}
