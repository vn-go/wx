package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

// ======================
//  Cấu trúc và registry
// ======================

type GrpcCall struct {
	GrpcEndPoint string
	Args         []byte
}

var grpcEndpoints = map[string]reflect.Value{}

func GrpcEndPointAdd(name string, fn interface{}) {
	grpcEndpoints[name] = reflect.ValueOf(fn)
}

// ======================
//  Encode / Decode args
// ======================

func EncodeArgs(args ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(args); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeArgs(data []byte) ([]reflect.Value, error) {
	var args []interface{}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&args); err != nil {
		return nil, err
	}
	vals := make([]reflect.Value, len(args))
	for i, a := range args {
		vals[i] = reflect.ValueOf(a)
	}
	return vals, nil
}

// ======================
//  Hàm gọi động
// ======================

func (g *GrpcCall) Call() (interface{}, error) {
	fn, ok := grpcEndpoints[g.GrpcEndPoint]
	if !ok {
		return nil, fmt.Errorf("endpoint not found: %s", g.GrpcEndPoint)
	}

	args, err := DecodeArgs(g.Args)
	if err != nil {
		return nil, fmt.Errorf("decode args: %w", err)
	}

	results := fn.Call(args)
	switch len(results) {
	case 0:
		return nil, nil
	case 1:
		return results[0].Interface(), nil
	default:
		out := make([]interface{}, len(results))
		for i, r := range results {
			out[i] = r.Interface()
		}
		return out, nil
	}
}

// ======================
//  Ví dụ hàm logic server
// ======================

func Test(a, b int) int {
	return a + b
}

// ======================
//  Mô phỏng client gọi
// ======================

func main() {
	// Đăng ký endpoint
	GrpcEndPointAdd("Test", Test)

	// Client chuẩn bị args
	argsData, err := EncodeArgs(5, 7)
	if err != nil {
		panic(err)
	}

	call := &GrpcCall{
		GrpcEndPoint: "Test",
		Args:         argsData,
	}

	// Server xử lý
	result, err := call.Call()
	if err != nil {
		panic(err)
	}

	fmt.Println("Kết quả trả về:", result) // Kết quả trả về: 12
}
