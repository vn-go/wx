package grpctest

import (
	"log"
	"testing"

	"github.com/vn-go/wx"
)

type UserController struct {
}

func (u *UserController) Login(h *wx.Handler, data struct {
	Username string
	Password string
}) (bool, error) {
	return true, nil
}
func TestLoging(t *testing.T) {
	httpClient := wx.NewHttpClient("http://localhost:8080")
	res, err := httpClient.PostJson("/api/user-controller/login", struct {
		Username string
		Password string
	}{"admin", "123456"})
	if err != nil {
		panic(err)
	}
	t.Log(res)
	// h, err := wx.MakeHandlerFromMethod[UserController]("Login")
	// if err != nil {
	// 	panic(err)
	// }
	// req, err := wx.Mock.JsonRequest(h.GetHttpMethod(), h.GetUriHandler(), struct {
	// 	Username string
	// 	Password string
	// }{"admin", "123456"})
	// if err != nil {
	// 	panic(err)
	// }
	// res := wx.Mock.NewRes()
	// h.Handler().ServeHTTP(res, req)
}
func TestXxx(t *testing.T) {
	client, err := wx.NewGrpcClient("localhost:50051")
	if err != nil {
		panic(err)
	}
	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("không kết nối được: %v", err)
	// }
	defer client.Close()
	res, err := client.Call(t.Context(), "main.user.Login", struct {
		Username string
		Password string
	}{"admin", "123456"})
	if err != nil {
		panic(err)
	}
	log.Printf("Kết quả: %s", res.Result)
	if res.Error != "" {
		log.Printf("Lỗi: %s", res.Error)
	}

}
func BenchmarkXxx(b *testing.B) {

	client, err := wx.NewGrpcClient("localhost:50051")
	if err != nil {
		panic(err)
	}

	defer client.Close()
	b.Run("grpc", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res, err := client.Call(b.Context(), "main.user.Login", struct {
				Username string
				Password string
			}{"admin", "123456"})
			if err != nil {
				panic(err)
			}

			if res.Error != "" {
				panic(res.Error)
			}
		}
	})
	b.Run("http", func(b *testing.B) {
		httpClient := wx.NewHttpClient("http://localhost:8080")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			_, err := httpClient.PostJson("/api/user-controller/login", struct {
				Username string
				Password string
			}{"admin", "123456"})
			if err != nil {
				panic(err)
			}
		}
	})
	b.Run("grpc-parallel", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {

			for pb.Next() {
				res, err := client.Call(b.Context(), "main.user.Login", struct {
					Username string
					Password string
				}{"admin", "123456"})
				if err != nil {
					panic(err)
				}

				if res.Error != "" {
					panic(res.Error)
				}
			}
		})
	})
	b.Run("http-parallel", func(b *testing.B) {
		httpClient := wx.NewHttpClient("http://localhost:8080")
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {

				_, err := httpClient.PostJson("/api/user-controller/login", struct {
					Username string
					Password string
				}{"admin", "123456"})
				if err != nil {
					panic(err)
				}
			}
		})
	})
}
