package wx

import "net/http"

type middleWareTypes struct {
	Cors func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Zip  func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

var cors = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Cho phép tất cả origin (cẩn thận với sản phẩm thật!)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Nếu là preflight request (OPTIONS), chỉ phản hồi 200
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Gọi tiếp handler chính
	next.ServeHTTP(w, r)
}

var MiddlWares = &middleWareTypes{
	Cors: cors,
	Zip:  zip,
}
