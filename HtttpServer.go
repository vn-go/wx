package wx

import (
	"fmt"
	"net/http"
	"strings"
)

var useSwagger bool = false

type htttpServer struct {
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

	mws []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

/*
This function create a new htttpServer, after calling this function call htttpServer.Start
*/
func NewHtttpServer(baseUrl string, port string, bind string) *htttpServer {
	if baseUrl[0] != '/' {
		baseUrl = "/" + baseUrl
	}
	if baseUrl[len(baseUrl)-1] == '/' {
		baseUrl = baseUrl[:len(baseUrl)-1]
	}
	baseUrl = strings.ReplaceAll(baseUrl, "//", "/")
	mux := http.NewServeMux()
	return &htttpServer{
		Port:    port,
		Bind:    bind,
		BaseUrl: baseUrl,

		mux: mux,
		mws: []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){},
	}

}
func (s *htttpServer) loadController() error {
	for _, x := range utils.Routes.UriList {
		fmt.Println("Registering route:", x)
		// s.mux.HandleFunc(x, func(w http.ResponseWriter, r *http.Request) {
		// 	route := utils.Routes.Data[x]
		// 	route.Info.Invoke(w, r)
		// })
		s.mux.HandleFunc(x, utils.Routes.Data[x].Info.Handler())

	}
	return nil

}
func (s *htttpServer) Start() error {
	err := s.loadController()
	if err != nil {
		return err
	}
	// handler cuối cùng gọi mux
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s.mux.ServeHTTP(w, r)
	})

	// Gắn middleware vào handler chain
	for i := len(s.mws) - 1; i >= 0; i-- {
		mw := s.mws[i]
		next := final
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
	// s. = &http.Server{
	// 	Addr:         addr,
	// 	Handler:      s.handler,
	// 	ReadTimeout:  10 * time.Second, // Giới hạn đọc request
	// 	WriteTimeout: 10 * time.Second, // Giới hạn ghi response
	// 	IdleTimeout:  60 * time.Second, // Cho keep-alive
	// }

	fmt.Println("Server listening at", addr)
	return s.ListenAndServe()
}
func (s *htttpServer) Middleware(fn func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *htttpServer {
	s.mws = append(s.mws, fn)
	return s
}
