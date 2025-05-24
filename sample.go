package urpc

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/ondbyte/urpc/protos"
	"google.golang.org/grpc"
)

type KVServer struct {
	store map[string]string
	protos.UnimplementedKVStoreServer
}

// Get implements protos.KVStoreServer.
func (k *KVServer) Get(ctx context.Context, gr *protos.GetRequest) (*protos.GetResponse, error) {
	v, ok := k.store[gr.Key]
	return &protos.GetResponse{Value: v, Found: ok}, nil
}

// Set implements protos.KVStoreServer.
func (k *KVServer) Set(ctx context.Context, sr *protos.SetRequest) (*protos.SetResponse, error) {
	k.store[sr.Key] = sr.Value
	return &protos.SetResponse{Success: true, Message: "set"}, nil
}

var _ protos.KVStoreServer = (*KVServer)(nil)

func StartGRPCServer(socketPath string) {
	os.Remove(socketPath)

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	protos.RegisterKVStoreServer(grpcServer, &KVServer{store: map[string]string{}})

	log.Println("gRPC server listening on", socketPath)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
