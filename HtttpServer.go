package wx

import (
	"fmt"
	"net/http"
	"time"
)

var useSwagger bool = false

type HtttpServer struct {

	// Port is the port the server will listen on.
	Port string
	// BaseUrl is the base URL of the server.
	BaseUrl string
	// Host is the host the server will listen on.
	Bind string
	// Handler is the HTTP handler for the server.
	handler http.Handler
	// server is the underlying http.Server.
	server *http.Server
	mux    *http.ServeMux

	mws []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func NewHtttpServer(baseUrl string, port string, bind string) *HtttpServer {
	mux := http.NewServeMux()
	return &HtttpServer{
		Port:    port,
		Bind:    bind,
		BaseUrl: baseUrl,

		mux: mux,
		mws: []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){},
	}

}
func (s *HtttpServer) Start() error {

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
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  10 * time.Second, // Giới hạn đọc request
		WriteTimeout: 10 * time.Second, // Giới hạn ghi response
		IdleTimeout:  60 * time.Second, // Cho keep-alive
	}

	fmt.Println("Server listening at", addr)
	return s.server.ListenAndServe()
}
