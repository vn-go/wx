package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

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

type UserClaims struct {
}
type BaseController struct {
}

func (c *BaseController) New() error {

	return nil
}

type HelloController struct {
}

func (c *HelloController) New() error {
	wx.SetAuth(c, func(ctx *wx.HttpContext[UserClaims]) error {
		return nil

	})
	return nil
}

type Auth struct {
}

func main() {
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// wx.Routes("api", &Media{})
	// wx.HandlerPost("test/hello", func(ctl *HelloController, ctx *wx.HttpContext, msg string) (string, error) {
	// 	return fmt.Sprintf("%s:Hello %s", ctx.Uri, msg), nil
	// })
	// wx.HandlerPost("test/hello2", func(ctl *HelloController, ctx *wx.HttpContext, msg string) (string, error) {
	// 	return fmt.Sprintf("%s:Hello %s", ctx.Uri, msg), nil
	// })
	// wx.HandlerGet("test/hello3", func(ctl *HelloController, ctx *wx.HttpContext) (string, error) {
	// 	return fmt.Sprintf("%s:Hello", ctx.Uri), nil
	// })
	wx.HandlerPost("test/hello4/{file}/{code}.mp4?filePath=&fileCode=", func(ctl *HelloController, ctx *wx.HttpContext[UserClaims], data struct {
		// Code     string
		// Name     string
		FileTest *multipart.FileHeader
	}) (string, error) {
		return fmt.Sprintf("%s:Hello", ctx.Uri), nil
	})
	// wx.HandlerForm("login", func(auth *Auth, ctx *wx.HttpContext[any], data struct {
	// 	Username string `json:"username"`
	// 	Password string `json:"password"`
	// }) (any, error) {
	// 	return data, nil
	// })
	//wx.Options.IsDebug = core.Services.Config.Debug
	//wx.Options.UsePool = true
	//routes.InitRoute()

	//wx.Routes("/api", &Hello{})
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	server.IsReleaseMode = true
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.OAuth2Password("/api/auth/login")
	swagger.Build()
	server.Start()
}
