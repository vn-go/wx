package wx

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var useSwagger bool = false

type httpServer struct {
	http.Server
	mux *http.ServeMux
	// Port is the port the server will listen on.
	Port string
	// BaseUrl is the base URL of the server.
	BaseUrl string
	// Host is the host the server will listen on.
	Bind string
	// Handler is the HTTP handler for the server.
	handler http.Handler
	// server is the underlying http.Server.

	mws           []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	MaxBodySize   uint64
	MaxUploadSize uint64
	IsReleaseMode bool
	uriParamList  []string
}

var currentServer *httpServer

/*
This function create a new htttpServer, after calling this function call htttpServer.Start
*/
func NewHttpServer(baseUrl string, port string, bind string) *httpServer {
	if baseUrl[0] != '/' {
		baseUrl = "/" + baseUrl
	}
	if baseUrl[len(baseUrl)-1] == '/' {
		baseUrl = baseUrl[:len(baseUrl)-1]
	}
	baseUrl = strings.ReplaceAll(baseUrl, "//", "/")
	mux := http.NewServeMux()
	currentServer = &httpServer{
		Port:    port,
		Bind:    bind,
		BaseUrl: baseUrl,

		mux:           mux,
		mws:           []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){},
		MaxBodySize:   20 << 20,
		MaxUploadSize: 20 << 20,
	}
	return currentServer

}
func (s *httpServer) SetMaxBodySize(val uint64) {
	s.MaxBodySize = val
}
func (s *httpServer) SetMaxUploadSize(val uint64) {
	s.MaxUploadSize = val
}
func (s *httpServer) loadController() error {
	for _, x := range utils.Routes.UriList {
		url := utils.Routes.Data[x].Info.uriHandler
		//utils.Routes.Data[x].Info.httpMethod = "get"
		urlShow := x
		if currentServer.BaseUrl != "" && !utils.Routes.Data[x].Info.IsAbsUri() {
			url = currentServer.BaseUrl + "/" + utils.Routes.Data[x].Info.uriHandler
			urlShow = currentServer.BaseUrl + "/" + utils.Routes.Data[x].Info.uri
		}

		fmt.Printf("%s\t->\tt%s\n", urlShow, url)
		s.mux.HandleFunc(url, utils.Routes.Data[x].Info.Handler())

	}
	for k, v := range handlerMapping {
		uri := fmt.Sprintf("/%s/%s", s.BaseUrl, k)
		for strings.HasPrefix(uri, "//") {
			uri = uri[1:]
		}
		fmt.Println("Registering route:", uri)
		s.mux.HandleFunc(uri, v)
	}
	return nil

}

type keyHandler string

var keyBeforeRequestCompleted = keyHandler("hack-before-request-completed")

func (s *httpServer) BeforeRequestCompleted(key string, w http.ResponseWriter, r *http.Request, fn http.HandlerFunc) (http.ResponseWriter, *http.Request) {

	if fnList, ok := r.Context().Value(keyBeforeRequestCompleted).([]http.HandlerFunc); ok {
		fnList = append(fnList, fn)
		ctx := context.WithValue(r.Context(), keyBeforeRequestCompleted, fnList)
		*r = *r.WithContext(ctx)
	} else {
		fnList = []http.HandlerFunc{fn}
		ctx := context.WithValue(r.Context(), keyBeforeRequestCompleted, fnList)
		*r = *r.WithContext(ctx)
	}
	// }
	// ctx := context.WithValue(r.Context(), keyBeforeRequestCompleted, afterHandlers)
	// *r = *r.WithContext(ctx)
	return w, r

}

func (s *httpServer) Use(mw func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	s.mws = append(s.mws, mw)
}
func (s *httpServer) Start() error {
	err := s.loadController()
	if err != nil {
		return err
	}
	// Handler cuối cùng (gốc)
	// final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	s.mux.ServeHTTP(w, r)
	// })
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		s.mux.ServeHTTP(w, r)

		// ⚠️ Đến đây có thể header đã gửi
		// => không nên set thêm header, chỉ nên thêm Trailer
		w.Header().Add("Trailer", "X-Process-Time")

		if rc := http.NewResponseController(w); rc != nil {
			// Dùng trailer — header này được phép gửi sau response
			w.Header().Set("X-Process-Time", fmt.Sprintf("%.2fms", float64(time.Since(start).Microseconds())/1000))
		}
	})
	// Xây chuỗi middleware từ cuối về đầu
	for i := len(s.mws) - 1; i >= 0; i-- {
		mw := s.mws[i]
		next := final // ⚠️ tạo biến cục bộ để tránh capture lỗi
		final = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw(w, r, next.ServeHTTP)
		})
	}

	s.handler = final
	s.Addr = fmt.Sprintf("%s:%s", s.Bind, s.Port)
	s.Handler = s.handler

	fmt.Println("Server listening at", s.Addr)
	return s.ListenAndServe()
}
func (s *httpServer) StartOld() error {
	err := s.loadController()
	if err != nil {
		return err
	}
	// handler cuối cùng gọi mux
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s.mux.ServeHTTP(w, r) // api handler
	})

	for i := len(s.mws) - 1; i >= 0; i-- {
		mw := s.mws[i]
		next := final // tạo bản sao tạm thời trong scope mới

		final = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw(w, r, next.ServeHTTP)
		})
	}
	s.handler = final

	addr := fmt.Sprintf("%s:%s", s.Bind, s.Port)
	// fmt.Println("Server listening at", addr)
	// return http.ListenAndServe(addr, s.handler)
	s.Addr = addr
	s.Handler = s.handler

	fmt.Println("Server listening at", addr)
	return s.ListenAndServe()
}
func (s *httpServer) Middleware(fn func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *httpServer {
	s.mws = append(s.mws, fn)
	return s
}
