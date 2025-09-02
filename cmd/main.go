package main

import (
	"fmt"

	"github.com/vn-go/wx"
)

type Oauth struct {
}

func (auth *Oauth) Login(ctx wx.Handler, body wx.Form[struct {
	Username string
	Password string
}]) {
	fmt.Println(body)
}
func main() {
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
