package main

import (
	"fmt"
	"regexp"
	"strings"
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
type ControllerTest struct {
}

// wx.Handler `route:"@/{*FilePath};method:get"`
func (c *ControllerTest) Hello(h struct {
	wx.Handler `route:"dx/{name} method:get"`
}) (string, error) {
	return "Hello", nil
}

func TestHellWorldV1(t *testing.T) {
	h, err := wx.MakeHandlerFromMethod[ControllerTest]("Hello")
	assert.NoError(t, err)
	fmt.Println(h.GetUriHandler())
	req, err := wx.Mock.NewGetRequest(h.GetUriHandler())
	res := wx.Mock.NewRes()
	h.Handler().ServeHTTP(res, req)
}
func TestHellWorldV2(t *testing.T) {
	wx.HandlerGet("hello2", func(controllerInstance *ControllerTest, ctx *wx.HttpContext[any]) (any, error) {
		return "Hello", nil
	})
	//h, err := wx.MakeHandlerFromMethod[ControllerTest]("Hello")
	//assert.NoError(t, err)
	req, err := wx.Mock.NewGetRequest("controllertest/hello2")
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	h, err := wx.GetHandler[ControllerTest]("hello2")
	assert.NoError(t, err)
	h.ServeHTTP(res, req)

}
func BenchmarkCompareV1AndV2(b *testing.B) {
	wx.HandlerGet("hello2", func(controllerInstance *ControllerTest, ctx *wx.HttpContext[any]) (any, error) {
		return "Hello", nil
	})
	b.Run("v1", func(b *testing.B) {
		h, err := wx.MakeHandlerFromMethod[ControllerTest]("Hello")
		assert.NoError(b, err)
		req, _ := wx.Mock.NewGetRequest(h.GetUriHandler())
		res := wx.Mock.NewRes()
		for i := 0; i < b.N; i++ {

			h.Handler().ServeHTTP(res, req)
		}

	})
	b.Run("V2", func(b *testing.B) {
		h, err := wx.GetHandler[ControllerTest]("hello2")
		assert.NoError(b, err)
		req, _ := wx.Mock.NewGetRequest("controllertest/hello2") //<-- dat o day de trang tinh vao ket qua test

		res := wx.Mock.NewRes() //<-- dat o day de trang tinh vao ket qua test
		for i := 0; i < b.N; i++ {
			//req, _ := wx.Mock.NewGetRequest("controllertest/hello2") //<-- cac ban test truoc cua v2 dat o day

			//res := wx.Mock.NewRes() //<-- cac ban test truoc cua v2 dat o day
			h.ServeHTTP(res, req)
		}

	})

}

/*
PS D:\code\go\wx\wx\cmd> go tool pprof cpu_v1.out
File: cmd.test.exe
Build ID: C:\Users\MSICYB~1\AppData\Local\Temp\go-build868172620\b001\cmd.test.exe2025-10-04 14:34:45.4332469 +0700 +07
Type: cpu
Time: 2025-10-04 14:34:45 +07
Duration: 2.22s, Total samples = 2.59s (116.87%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 930ms, 35.91% of 2590ms total
Dropped 102 nodes (cum <= 12.95ms)
Showing top 10 nodes out of 175

	 flat  flat%   sum%        cum   cum%
	200ms  7.72%  7.72%      200ms  7.72%  runtime.stdcall2
	 90ms  3.47% 11.20%      100ms  3.86%  net/url.escape
	 90ms  3.47% 14.67%      210ms  8.11%  reflect.Value.call
	 90ms  3.47% 18.15%       90ms  3.47%  runtime.nextFreeFast (inline)
	 90ms  3.47% 21.62%      150ms  5.79%  runtime.scanobject
	 90ms  3.47% 25.10%       90ms  3.47%  runtime.stdcall1
	 80ms  3.09% 28.19%      170ms  6.56%  runtime.cgocall
	 70ms  2.70% 30.89%       70ms  2.70%  net/url.unescape
	 70ms  2.70% 33.59%      590ms 22.78%  runtime.mallocgc
	 60ms  2.32% 35.91%       60ms  2.32%  aeshashbody

----
PS D:\code\go\wx\wx\cmd> go tool pprof cpu_v2.out
File: cmd.test.exe
Build ID: D:\code\go\wx\wx\cmd\cmd.test.exe2025-10-04 14:34:48.0584204 +0700 +07
Type: cpu
Time: 2025-10-04 14:35:00 +07
Duration: 1.41s, Total samples = 1.80s (127.71%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 730ms, 40.56% of 1800ms total
Showing top 10 nodes out of 222

	 flat  flat%   sum%        cum   cum%
	180ms 10.00% 10.00%      180ms 10.00%  runtime.stdcall2
	 90ms  5.00% 15.00%      440ms 24.44%  runtime.mallocgcSmallScanNoHeader
	 80ms  4.44% 19.44%       80ms  4.44%  runtime.stdcall1
	 70ms  3.89% 23.33%       70ms  3.89%  runtime.memclrNoHeapPointers
	 70ms  3.89% 27.22%       90ms  5.00%  runtime.scanblock
	 50ms  2.78% 30.00%       50ms  2.78%  internal/chacha8rand.block
	 50ms  2.78% 32.78%       60ms  3.33%  runtime.(*mspan).writeHeapBitsSmall
	 50ms  2.78% 35.56%       50ms  2.78%  runtime.nextFreeFast (inline)
	 50ms  2.78% 38.33%       60ms  3.33%  runtime.typePointers.next
	 40ms  2.22% 40.56%       50ms  2.78%  net/url.escape

(pprof)
--- duny syn.pool----
PS D:\code\go\wx\wx\cmd> go tool pprof cpu_v2.out
File: cmd.test.exe
Build ID: C:\Users\MSICYB~1\AppData\Local\Temp\go-build3923474888\b001\cmd.test.exe2025-10-04 15:03:23.6377762 +0700 +07
Type: cpu
Time: 2025-10-04 15:03:23 +07
Duration: 1.72s, Total samples = 2s (116.49%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 820ms, 41.00% of 2000ms total
Showing top 10 nodes out of 226

	 flat  flat%   sum%        cum   cum%
	160ms  8.00%  8.00%      160ms  8.00%  runtime.stdcall2
	140ms  7.00% 15.00%      140ms  7.00%  runtime.nextFreeFast (inline)
	100ms  5.00% 20.00%      100ms  5.00%  runtime.stdcall1
	 90ms  4.50% 24.50%      210ms 10.50%  runtime.scanobject
	 70ms  3.50% 28.00%       70ms  3.50%  net/url.unescape
	 60ms  3.00% 31.00%       60ms  3.00%  runtime.findObject
	 50ms  2.50% 33.50%       80ms  4.00%  internal/runtime/maps.(*Iter).Next
	 50ms  2.50% 36.00%      540ms 27.00%  runtime.mallocgc
	 50ms  2.50% 38.50%      410ms 20.50%  runtime.mallocgcSmallScanNoHeader
	 50ms  2.50% 41.00%       50ms  2.50%  runtime.procyield

(pprof)
*/
//go test -benchmem -run=^$ -bench ^BenchmarkCompareV1AndV2InParallel$ -cpuprofile cpu_parallel.out github.com/vn-go/wx/cmd
//go tool pprof cpu_parallel.out
func BenchmarkCompareV1AndV2InParallel(b *testing.B) {
	wx.HandlerGet("hello2", func(controllerInstance *ControllerTest, ctx *wx.HttpContext[any]) (any, error) {
		return "Hello", nil
	})
	b.Run("v1", func(b *testing.B) {
		h, err := wx.MakeHandlerFromMethod[ControllerTest]("Hello")
		assert.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			req, _ := wx.Mock.NewGetRequest(h.GetUriHandler())
			res := wx.Mock.NewRes()
			for p.Next() {
				h.Handler().ServeHTTP(res, req)
			}
		})

	})
	b.Run("V2", func(b *testing.B) {
		h, err := wx.GetHandler[ControllerTest]("hello2")
		assert.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			req, _ := wx.Mock.NewGetRequest("controllertest/hello2") //<-- dat o day de trang tinh vao ket qua test

			res := wx.Mock.NewRes() //<-- dat o day de trang tinh vao ket qua test
			for p.Next() {
				h.ServeHTTP(res, req)
			}

		})

	})

}
func TestInspect(t *testing.T) {
	// NewInstance[Controller, UserTest]("test-001")
	// fx := NewInstance[Controller, UserTest]("test-002")
	//fx := cachFn["test-001"]
	//fx(UserTest{})
	h, _ := wx.MakeHandlerFromMethod[Media]("Files3")
	//{*FilePath}
	fmt.Println(h.GetUriHandler())
	fmt.Println(h)

}

type inspectUriInfo struct {
	// mainUri is the base path used for mapping handlers.
	// 1. For "files/{*FilePath}", mainUri is 'files'.
	// 2. For "file?path={path}&name={name}", mainUri is 'file'.
	// 3. For "file/{path}/{name}", mainUri is 'file'.
	mainUri string
	// swaggerUri is the standardized URI template used for documentation (Swagger/OpenAPI).
	// 1. For "files/*path", swaggerUri is "files/{path}".
	// 2. For "/media/files3/{FilePath}/{filename}?test={code}", swaggerUri is the same.
	swaggerUri string
	// uriRegex is the regular expression used by the routing engine to match incoming requests
	// and extract path parameters.
	// 1. For "media/files/{*FilePath}" -> uriRegex = "media/files/(.*)"
	// 2. For "media/files2/{FilePath}/{filename}" -> uriRegex = "^media/files2/([^/]+)/([^/]+)$"
	// 3. For "media/files3/{FilePath}/{filename}?test={code}" -> uriRegex = "^media/files3/([^/]+)/([^/]+)$" (captures up to the first '?')
	uriRegex string
	// isRegex is set to true if the URI contains path parameters that require a regex match
	// (i.e., contains "{" or "*").
	isRegex bool
}

// placeholderRegex matches the two types of path placeholders: {param} and {*param}

// Hỗ trợ {param}, {*param}, *param
var placeholderRegex = regexp.MustCompile(`\{[*]?[^}]+\}|\*[^/]+`)

func inspectUri(uri string) inspectUriInfo {
	info := inspectUriInfo{
		swaggerUri: uri,
	}

	// --- 1. Xác định phần path chính ---
	pathPart := uri
	if qIndex := strings.Index(uri, "?"); qIndex != -1 {
		pathPart = uri[:qIndex]
	}
	segments := strings.Split(strings.Trim(pathPart, "/"), "/")
	if len(segments) > 0 && segments[0] != "" {
		info.mainUri = segments[0]
	} else {
		info.mainUri = strings.Trim(uri, "/")
	}

	// --- 2. Nếu không có { hoặc * thì không cần regex ---
	if !strings.Contains(uri, "{") && !strings.Contains(uri, "*") {
		info.isRegex = false
		return info
	}

	// --- 3. Xây swaggerUri và uriRegex ---
	swaggerBuilder := strings.Builder{}
	regexBuilder := strings.Builder{}
	regexBuilder.WriteString("^")

	lastIndex := 0
	matches := placeholderRegex.FindAllStringIndex(pathPart, -1)

	for _, match := range matches {
		start, end := match[0], match[1]
		literal := pathPart[lastIndex:start]
		regexBuilder.WriteString(regexp.QuoteMeta(literal))
		swaggerBuilder.WriteString(literal)

		placeholder := pathPart[start:end]

		// --- Xác định tên tham số thật ---
		name := extractParamName(placeholder)

		if strings.Contains(placeholder, "*") {
			// Catch-all: {*FilePath} hoặc *FilePath
			swaggerBuilder.WriteString("{" + name + "}")
			regexBuilder.WriteString("(.*)")
		} else {
			// Normal placeholder: {filename}
			swaggerBuilder.WriteString("{" + name + "}")
			regexBuilder.WriteString("([^/]+)")
		}

		lastIndex = end
	}

	// --- 4. Thêm phần literal còn lại ---
	if lastIndex < len(pathPart) {
		literal := pathPart[lastIndex:]
		regexBuilder.WriteString(regexp.QuoteMeta(literal))
		swaggerBuilder.WriteString(literal)
	}

	regexBuilder.WriteString("$")

	// --- 5. Gán kết quả ---
	info.uriRegex = regexBuilder.String()
	info.swaggerUri = swaggerBuilder.String()

	// Nếu có query string, nối thêm vào swaggerUri (Swagger có thể mô tả query params riêng)
	if qIndex := strings.Index(uri, "?"); qIndex != -1 {
		info.swaggerUri += uri[qIndex:]
	}

	info.isRegex = true
	return info
}

// extractParamName trích tên thật bên trong {param}, {*param}, hoặc *param
func extractParamName(s string) string {
	s = strings.Trim(s, "{}")
	s = strings.TrimPrefix(s, "*")
	s = strings.TrimSpace(s)
	if s == "" {
		return "param"
	}
	return s
}
func TestInspecuUri(t *testing.T) {
	tests := []string{
		"files/{*FilePath}",
		"file?path={path}&name={name}",
		"file/{path}/{name}",
		"/media/files3/{FilePath}/{filename}?test={code}",
		"media/files/*wild",
	}

	for _, t := range tests {
		info := inspectUri(t)
		fmt.Printf("URI: %-50s\n", t)
		fmt.Printf(" mainUri:    %s\n", info.mainUri)
		fmt.Printf(" swaggerUri: %s\n", info.swaggerUri)
		fmt.Printf(" uriRegex:   %s\n", info.uriRegex)
		fmt.Printf(" isRegex:    %v\n", info.isRegex)
		fmt.Println(strings.Repeat("-", 60))
	}
}
