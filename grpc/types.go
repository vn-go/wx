package grpc

import (
	"bytes"
	"encoding/gob"
	"reflect"

	pb "github.com/vn-go/wx/grpc/pb"
)

// ======================
//  Registry
// ======================

var grpcEndpoints = map[string]reflect.Value{}

func GrpcEndPointAdd(name string, fn interface{}) {
	grpcEndpoints[name] = reflect.ValueOf(fn)
}

// ======================
//  Encode / Decode
// ======================

func decodeArgs(data []byte) ([]reflect.Value, error) {
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

func encodeResult(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ======================
//  Service implement
// ======================

type dynamicServer struct {
	pb.UnimplementedDynamicServiceServer
}
