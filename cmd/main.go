package main

import (
	"github.com/vn-go/wx"
)

type Oauth struct {
}
type LoginForm struct {
	Username string
	Password string
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
