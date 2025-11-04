package wx

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"sync"

	"github.com/vn-go/wx/grpc/pb"
	"google.golang.org/grpc"
)

type grpcClientType struct {
	conn   *grpc.ClientConn
	client pb.DynamicServiceClient
}

func (g *grpcClientType) Close() {
	g.conn.Close()
}
func (g *grpcClientType) Call(ctx context.Context, endpoint string, args interface{}) (*pb.GrpcResult, error) {
	// Encode args thành []byte bằng gob (SIÊU NHANH cho Go-to-Go)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(args); err != nil {

		return nil, err
	}

	// Gửi request
	return g.client.Invoke(ctx, &pb.GrpcCall{
		GrpcEndpoint: endpoint,
		Endcoder:     "gob",
		Args:         buf.Bytes(),
	})
}
func (g *grpcClientType) CallJSON(ctx context.Context, endpoint string, args interface{}) (*pb.GrpcResult, error) {
	// Encode args thành []byte bằng gob (SIÊU NHANH cho Go-to-Go)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(args); err != nil {
		return nil, err
	}

	// Gửi request
	return g.client.Invoke(ctx, &pb.GrpcCall{
		GrpcEndpoint: endpoint,
		Endcoder:     "json",
		Args:         buf.Bytes(),
	})
}

var GrpcClient = &grpcClientType{}

type initNewGrpcClient struct {
	val  *grpcClientType
	err  error
	once sync.Once
}

var initNewGrpcClientMap sync.Map

func NewGrpcClient(target string) (*grpcClientType, error) {
	a, _ := initNewGrpcClientMap.LoadOrStore(target, &initNewGrpcClient{})
	i := a.(*initNewGrpcClient)
	i.once.Do(func() {
		conn, err := grpc.NewClient(target, grpc.WithInsecure()) //grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			i.err = err
			return
		}
		i.val = &grpcClientType{
			conn:   conn,
			client: pb.NewDynamicServiceClient(conn),
		}

	})
	return i.val, i.err
}
