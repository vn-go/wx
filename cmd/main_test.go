package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/wx"
)

func TestLogin(t *testing.T) {
	h, err := wx.MakeHandlerFromMethod[Oauth]("Login")
	assert.NoError(t, err)
	req, err := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), struct {
		Username string
		Password string
	}{
		Username: "admin",
		Password: "admin",
	})
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	h.Handler().ServeHTTP(res, req)
	t.Log(res)

}
func TestLogin2(t *testing.T) {
	h, err := wx.MakeHandlerFromMethod[Oauth]("Login2")
	assert.NoError(t, err)
	req, err := wx.Mock.JsonRequest(h.GetHttpMethod(), h.GetUriHandler(), struct {
		Username string
		Password string
	}{
		Username: "admin",
		Password: "admin",
	})
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	h.Handler().ServeHTTP(res, req)
	t.Log(res)

}
func TestUpload(t *testing.T) {
	h, err := wx.MakeHandlerFromMethod[Oauth]("Upload")
	assert.NoError(t, err)
	req, err := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), struct {
		Username string
		Password string
	}{
		Username: "admin",
		Password: "admin",
	})
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	h.Handler().ServeHTTP(res, req)
	t.Log(res)

}
func BenchmarkLogin(b *testing.B) {
	b.Run("without sync.pool", func(t *testing.B) {
		wx.Options.UsePool = false
		h, _ := wx.MakeHandlerFromMethod[Oauth]("Login")
		//assert.NoError(t, err)
		req, _ := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), struct {
			Username string
			Password string
		}{
			Username: "admin",
			Password: "admin",
		})
		res := wx.Mock.NewRes()
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			//assert.NoError(t, err)

			h.Handler().ServeHTTP(res, req)
			//t.Log(res)
		}

	})
	b.Run("with sync.pool", func(t *testing.B) {
		wx.Options.UsePool = true
		h, _ := wx.MakeHandlerFromMethod[Oauth]("Login")
		//assert.NoError(t, err)
		req, _ := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), LoginForm{
			Username: "admin",
			Password: "admin",
		})
		res := wx.Mock.NewRes()
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			//assert.NoError(t, err)

			h.Handler().ServeHTTP(res, req)
			//t.Log(res)
		}

	})
}
func BenchmarkLoginParallel(b *testing.B) {
	b.RunParallel(func(t *testing.PB) {
		wx.Options.UsePool = false
		h, _ := wx.MakeHandlerFromMethod[Oauth]("Login")
		//assert.NoError(t, err)
		req, _ := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), struct {
			Username string
			Password string
		}{
			Username: "admin",
			Password: "admin",
		})
		res := wx.Mock.NewRes()

		for t.Next() {
			//assert.NoError(t, err)

			h.Handler().ServeHTTP(res, req)

			//b.Log(res)
		}

	})
	// b.RunParallel(func(t *testing.PB) {
	// 	wx.Options.UsePool = true
	// 	h, _ := wx.MakeHandlerFromMethod[Oauth]("Login")
	// 	//assert.NoError(t, err)
	// 	req, _ := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), LoginForm{
	// 		Username: "admin",
	// 		Password: "admin",
	// 	})
	// 	res := wx.Mock.NewRes()

	// 	for t.Next() {
	// 		//assert.NoError(t, err)

	// 		h.Handler().ServeHTTP(res, req)
	// 		//t.Log(res)
	// 	}

	// })
}

/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginParallel$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginParallel-16    	 2435820	       746.0 ns/op	    3548 B/op	      35 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	6.365s
---
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginParallel$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginParallel-16    	 1726064	       618.5 ns/op	    3084 B/op	      35 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	3.414s
*/
func BenchmarkLogin2(b *testing.B) {
	b.Run("without sync.pool", func(t *testing.B) {
		wx.Options.UsePool = false
		h, _ := wx.MakeHandlerFromMethod[Oauth]("Login2")
		//assert.NoError(t, err)
		req, _ := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), struct {
			Username string
			Password string
		}{
			Username: "admin",
			Password: "admin",
		})
		res := wx.Mock.NewRes()
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			//assert.NoError(t, err)

			h.Handler().ServeHTTP(res, req)
			//t.Log(res)
		}

	})
	b.Run("with sync.pool", func(t *testing.B) {
		wx.Options.UsePool = true
		h, _ := wx.MakeHandlerFromMethod[Oauth]("Login2")
		//assert.NoError(t, err)
		req, _ := wx.Mock.FormRequest(h.GetHttpMethod(), h.GetUriHandler(), LoginForm{
			Username: "admin",
			Password: "admin",
		})
		res := wx.Mock.NewRes()
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			//assert.NoError(t, err)

			h.Handler().ServeHTTP(res, req)
			//t.Log(res)
		}

	})
}

/*
--- No 1------
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  613933	      1643 ns/op	     781 B/op	      16 allocs/op
BenchmarkLogin/with_sync.pool-16            	  731961	      1780 ns/op	     911 B/op	      18 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	4.794s
--- 003----
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  609432	      1742 ns/op	     718 B/op	      16 allocs/op
BenchmarkLogin/without_sync.pool-16         	  581112	      1768 ns/op	     779 B/op	      17 allocs/op

BenchmarkLogin/with_sync.pool-16            	  643310	      1758 ns/op	     752 B/op	      18 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	4.531s
---004---
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  669085	      1693 ns/op	     708 B/op	      16 allocs/op
BenchmarkLogin/with_sync.pool-16            	  642627	      2034 ns/op	     808 B/op	      19 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	5.053s
----
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  721431	      1700 ns/op	     794 B/op	      16 allocs/op
BenchmarkLogin/with_sync.pool-16            	  613471	      1888 ns/op	     765 B/op	      18 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	4.672s
---------

Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  697820	      1658 ns/op	     704 B/op	      16 allocs/op
BenchmarkLogin/with_sync.pool-16            	  612224	      1750 ns/op	     717 B/op	      16 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	5.885s
---
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  686419	      1626 ns/op	     705 B/op	      16 allocs/op
BenchmarkLogin/with_sync.pool-16            	  625784	      1727 ns/op	     715 B/op	      16 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	3.866s
*/
/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  254341	      4440 ns/op	    2895 B/op	      35 allocs/op
BenchmarkLogin/with_sync.pool-16            	  282067	      3873 ns/op	    2792 B/op	      35 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	4.853s
---
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin/without_sync.pool-16         	  272132	      3900 ns/op	    2826 B/op	      35 allocs/op
BenchmarkLogin/with_sync.pool-16            	  229406	      4867 ns/op	    3010 B/op	      35 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	4.436s
*/
/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLogin2$ github.com/vn-go/wx/cmd

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx/cmd
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLogin2/without_sync.pool-16         	  310581	      3861 ns/op	    3569 B/op	      35 allocs/op
BenchmarkLogin2/with_sync.pool-16            	  320131	      3759 ns/op	    3517 B/op	      35 allocs/op
PASS
ok  	github.com/vn-go/wx/cmd	4.974s
*/
