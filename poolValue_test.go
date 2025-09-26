package wx

import (
	"reflect"
	"testing"
)

type UserData struct {
	Username string
	Password string
}

func Login(d *Form[UserData]) bool {
	return d.Data.Password == "admin" && d.Data.Username == "admon"
}
func BenchmarkTestPoolValue(b *testing.B) {
	b.Run("use sync.pool", func(b *testing.B) {
		typ := reflect.TypeOf(Form[UserData]{})
		for i := 0; i < b.N; i++ {
			func() {
				dataVale := (PoolValues.Get(typ))
				defer PoolValues.Put(dataVale)
				formData := dataVale.Elem().Interface().(Form[UserData])
				Login(&formData)
			}()

		}

	})
	b.Run("no sync.pool", func(b *testing.B) {
		typ := reflect.TypeOf(Form[UserData]{})
		for i := 0; i < b.N; i++ {
			func() {
				dataVale := reflect.New(typ)
				formData := dataVale.Elem().Interface().(Form[UserData])
				Login(&formData)
			}()

		}

	})
}

/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkTestPoolValue$ github.com/vn-go/wx

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkTestPoolValue/use_sync.pool-16         	 5955298	       206.3 ns/op	     152 B/op	       4 allocs/op
PASS
ok  	github.com/vn-go/wx	2.804s
----
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkTestPoolValue$ github.com/vn-go/wx

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkTestPoolValue/use_sync.pool-16         	 7422556	       150.0 ns/op	     120 B/op	       3 allocs/op
PASS
ok  	github.com/vn-go/wx	3.570s


*/
/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkTestPoolValue$ github.com/vn-go/wx

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkTestPoolValue/use_sync.pool-16         	 7760961	       147.6 ns/op	      48 B/op	       2 allocs/op
BenchmarkTestPoolValue/no_sync.pool-16          	19967086	        57.85 ns/op	      64 B/op	       2 allocs/op
PASS
ok  	github.com/vn-go/wx	4.407s
----
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkTestPoolValue$ github.com/vn-go/wx

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkTestPoolValue/use_sync.pool-16         	 6902581	       175.2 ns/op	     104 B/op	       3 allocs/op
BenchmarkTestPoolValue/no_sync.pool-16          	21106173	        58.65 ns/op	      64 B/op	       2 allocs/op
PASS
ok  	github.com/vn-go/wx	4.151s


*/
/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkTestPoolValue$ github.com/vn-go/wx

goos: windows
goarch: amd64
pkg: github.com/vn-go/wx
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkTestPoolValue/use_sync.pool-16         	13306413	        94.04 ns/op	      64 B/op	       2 allocs/op
BenchmarkTestPoolValue/no_sync.pool-16          	20193316	        58.56 ns/op	      64 B/op	       2 allocs/op
PASS
*/
/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkTestPoolValue$ github.com/vn-go/wx

go: downloading github.com/stretchr/testify v1.11.1
go: downloading gopkg.in/yaml.v3 v3.0.1
go: downloading github.com/davecgh/go-spew v1.1.1
go: downloading github.com/pmezard/go-difflib v1.0.0
goos: windows
goarch: amd64
pkg: github.com/vn-go/wx
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkTestPoolValue/use_sync.pool-16         	17633218	        65.25 ns/op	      64 B/op	       2 allocs/op
BenchmarkTestPoolValue/no_sync.pool-16          	20829933	        61.52 ns/op	      64 B/op	       2 allocs/op
PASS
ok  	github.com/vn-go/wx	5.366s
*/
