package wx

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var LogAccessTokenClaims = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		fmt.Println("[Token] Không có Bearer token trong request")
		next.ServeHTTP(w, r)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Thử decode JWT
	parts := strings.Split(token, ".")
	if len(parts) == 3 {
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			fmt.Println("[Token] Lỗi decode payload:", err)
		} else {
			var claims map[string]interface{}
			if err := json.Unmarshal(payload, &claims); err == nil {
				fmt.Println("[Token] JWT Claims:", claims)
			}
		}
	} else {
		// Token không phải JWT → gọi introspection endpoint
		introspectOpaqueToken(token)
	}

	next.ServeHTTP(w, r)
}

// Hàm gọi introspection endpoint cho token opaque
func introspectOpaqueToken(token string) {
	// ⚠️ URL và client_id, client_secret để trong config hoặc env
	req, _ := http.NewRequest("POST", "https://auth.example.com/oauth/introspect",
		strings.NewReader("token="+token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("client_id_here", "client_secret_here")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[Token] Lỗi introspection:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("[Token] Opaque token introspection result:", string(body))
}
