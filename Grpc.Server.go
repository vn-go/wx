package wx

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/libp2p/go-reuseport"
	"github.com/vn-go/wx/grpc/pb"
	"google.golang.org/grpc"
)

type endpoint struct {
	receiverType reflect.Type
	method       reflect.Method
	receiver     reflect.Value
	paramType    reflect.Type
}

var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func (e endpoint) gobInvoke(args []byte) ([]byte, error) {

	paramPtr := reflect.New(e.paramType)
	paramValue := paramPtr.Interface()

	// === 2. Decode args []byte → struct bằng gob ===
	if err := gob.NewDecoder(bytes.NewReader(args)).Decode(paramValue); err != nil {
		return nil, err
	}

	// receiver := reflect.New(e.receiverType)
	results := e.method.Func.Call([]reflect.Value{e.receiver, paramPtr.Elem()})

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	if err := gob.NewEncoder(buf).Encode(results[0].Interface()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *endpoint) jsonInvoke(args []byte) ([]byte, error) {
	in := reflect.New(e.receiverType.In(0)).Interface()
	err := json.Unmarshal(args, in)
	if err != nil {
		return nil, err
	}
	out := e.method.Func.Call([]reflect.Value{reflect.ValueOf(in)})[0].Interface()
	return json.Marshal(out)
}

type server struct {
	pb.UnimplementedDynamicServiceServer
	enpointMap map[string]endpoint
}

func (s *server) Invoke(ctx context.Context, in *pb.GrpcCall) (*pb.GrpcResult, error) {

	key := strings.ToLower(in.GrpcEndpoint)
	if endpoint, ok := s.enpointMap[key]; ok {
		if in.Endcoder == "json" {
			ret := &pb.GrpcResult{}
			var err error
			ret.Result, err = endpoint.jsonInvoke(in.Args)
			if err != nil {
				ret.Error = err.Error()
			}
			return ret, nil
		}
		if in.Endcoder == "gob" {
			ret := &pb.GrpcResult{}
			var err error
			ret.Result, err = endpoint.gobInvoke(in.Args)
			if err != nil {
				ret.Error = err.Error()
			}
			return ret, nil
		}
	}
	log.Printf("Received: endpoint=%s, args=%s", in.GrpcEndpoint, string(in.Args))
	return &pb.GrpcResult{
		Result: append([]byte("ECHO: "), in.Args...),
	}, nil
}

type grpcServerType struct {
	server *server
	//lis    net.Listener
}

func (grpcServer *grpcServerType) AddServices(services ...any) {
	for _, service := range services {
		typ := reflect.TypeOf(service)
		if typ.Kind() != reflect.Ptr {
			typ = reflect.PointerTo(typ)
		}
		eleType := typ.Elem()
		for i := 0; i < typ.NumMethod(); i++ {
			method := typ.Method(i)
			key := strings.ToLower(fmt.Sprintf("%s.%s.%s", eleType.PkgPath(), eleType.Name(), method.Name))
			if _, ok := grpcServer.server.enpointMap[strings.ToLower(key)]; !ok {
				grpcServer.server.enpointMap[key] = endpoint{
					receiverType: eleType,
					method:       method,
					receiver:     reflect.New(eleType),
					paramType:    method.Type.In(1),
				}
			}
		}
	}
}

func NewGrpcTPCServer() *grpcServerType {

	return &grpcServerType{
		server: &server{
			enpointMap: make(map[string]endpoint),
		},
		// lis:    lis,
	}
}
func (grpcServer *grpcServerType) Start(port string) error {
	n := runtime.NumCPU()
	for i := 0; i < n; i++ {
		go func() {
			lis, err := reuseport.Listen("tcp", port)
			if err != nil {
				log.Fatalf("reuseport listen: %v", err)
			}
			s := grpc.NewServer(
				grpc.MaxConcurrentStreams(2048),
				grpc.NumStreamWorkers(uint32(runtime.NumCPU())),
			)
			pb.RegisterDynamicServiceServer(s, grpcServer.server)
			log.Printf("[CPU %d] Listening on %s", i, port)
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Serve: %v", err)
			}
		}()
	}
	select {} // block main
}
