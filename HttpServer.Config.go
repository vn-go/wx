package wx

// import (
// 	"context"
// 	"crypto/tls"
// 	"log"
// 	"net"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// type ServerConfig struct {
// 	// Addr optionally specifies the TCP address for the server to listen on,
// 	// in the form "host:port". If empty, ":http" (port 80) is used.
// 	// The service names are defined in RFC 6335 and assigned by IANA.
// 	// See net.Dial for details of the address format.
// 	Addr string

// 	Handler Handler // handler to invoke, http.DefaultServeMux if nil

// 	// DisableGeneralOptionsHandler, if true, passes "OPTIONS *" requests to the Handler,
// 	// otherwise responds with 200 OK and Content-Length: 0.
// 	DisableGeneralOptionsHandler bool

// 	// TLSConfig optionally provides a TLS configuration for use
// 	// by ServeTLS and ListenAndServeTLS. Note that this value is
// 	// cloned by ServeTLS and ListenAndServeTLS, so it's not
// 	// possible to modify the configuration with methods like
// 	// tls.Config.SetSessionTicketKeys. To use
// 	// SetSessionTicketKeys, use Server.Serve with a TLS Listener
// 	// instead.
// 	TLSConfig *tls.Config

// 	// ReadTimeout is the maximum duration for reading the entire
// 	// request, including the body. A zero or negative value means
// 	// there will be no timeout.
// 	//
// 	// Because ReadTimeout does not let Handlers make per-request
// 	// decisions on each request body's acceptable deadline or
// 	// upload rate, most users will prefer to use
// 	// ReadHeaderTimeout. It is valid to use them both.
// 	ReadTimeout time.Duration

// 	// ReadHeaderTimeout is the amount of time allowed to read
// 	// request headers. The connection's read deadline is reset
// 	// after reading the headers and the Handler can decide what
// 	// is considered too slow for the body. If zero, the value of
// 	// ReadTimeout is used. If negative, or if zero and ReadTimeout
// 	// is zero or negative, there is no timeout.
// 	ReadHeaderTimeout time.Duration

// 	// WriteTimeout is the maximum duration before timing out
// 	// writes of the response. It is reset whenever a new
// 	// request's header is read. Like ReadTimeout, it does not
// 	// let Handlers make decisions on a per-request basis.
// 	// A zero or negative value means there will be no timeout.
// 	WriteTimeout time.Duration

// 	// IdleTimeout is the maximum amount of time to wait for the
// 	// next request when keep-alives are enabled. If zero, the value
// 	// of ReadTimeout is used. If negative, or if zero and ReadTimeout
// 	// is zero or negative, there is no timeout.
// 	IdleTimeout time.Duration

// 	// MaxHeaderBytes controls the maximum number of bytes the
// 	// server will read parsing the request header's keys and
// 	// values, including the request line. It does not limit the
// 	// size of the request body.
// 	// If zero, DefaultMaxHeaderBytes is used.
// 	MaxHeaderBytes int

// 	// TLSNextProto optionally specifies a function to take over
// 	// ownership of the provided TLS connection when an ALPN
// 	// protocol upgrade has occurred. The map key is the protocol
// 	// name negotiated. The Handler argument should be used to
// 	// handle HTTP requests and will initialize the Request's TLS
// 	// and RemoteAddr if not already set. The connection is
// 	// automatically closed when the function returns.
// 	// If TLSNextProto is not nil, HTTP/2 support is not enabled
// 	// automatically.
// 	TLSNextProto map[string]func(*Server, *tls.Conn, Handler)

// 	// ConnState specifies an optional callback function that is
// 	// called when a client connection changes state. See the
// 	// ConnState type and associated constants for details.
// 	ConnState func(net.Conn, ConnState)

// 	// ErrorLog specifies an optional logger for errors accepting
// 	// connections, unexpected behavior from handlers, and
// 	// underlying FileSystem errors.
// 	// If nil, logging is done via the log package's standard logger.
// 	ErrorLog *log.Logger

// 	// BaseContext optionally specifies a function that returns
// 	// the base context for incoming requests on this server.
// 	// The provided Listener is the specific Listener that's
// 	// about to start accepting requests.
// 	// If BaseContext is nil, the default is context.Background().
// 	// If non-nil, it must return a non-nil context.
// 	BaseContext func(net.Listener) context.Context

// 	// ConnContext optionally specifies a function that modifies
// 	// the context used for a new connection c. The provided ctx
// 	// is derived from the base context and has a ServerContextKey
// 	// value.
// 	ConnContext func(ctx context.Context, c net.Conn) context.Context

// 	// HTTP2 configures HTTP/2 connections.
// 	//
// 	// This field does not yet have any effect.
// 	// See https://go.dev/issue/67813.
// 	HTTP2 *HTTP2Config

// 	// Protocols is the set of protocols accepted by the server.
// 	//
// 	// If Protocols includes UnencryptedHTTP2, the server will accept
// 	// unencrypted HTTP/2 connections. The server can serve both
// 	// HTTP/1 and unencrypted HTTP/2 on the same address and port.
// 	//
// 	// If Protocols is nil, the default is usually HTTP/1 and HTTP/2.
// 	// If TLSNextProto is non-nil and does not contain an "h2" entry,
// 	// the default is HTTP/1 only.
// 	Protocols *Protocols

// 	inShutdown atomic.Bool // true when server is in shutdown

// 	disableKeepAlives atomic.Bool
// 	nextProtoOnce     sync.Once // guards setupHTTP2_* init
// 	nextProtoErr      error     // result of http2.ConfigureServer if used

// 	mu         sync.Mutex
// 	listeners  map[*net.Listener]struct{}
// 	activeConn map[*conn]struct{}
// 	onShutdown []func()

// 	listenerGroup sync.WaitGroup
// }
