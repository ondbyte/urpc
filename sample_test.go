package urpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ondbyte/urpc"
	"github.com/ondbyte/urpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCServerUnixSocket(t *testing.T) {
	socketPath := urpc.GetSocketPath()

	go urpc.StartGRPCServer(socketPath)

	// Wait for the server to start
	time.Sleep(1 * time.Second)
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
	}
	// Dial using the Unix socket
	conn, err := grpc.NewClient(
		"unix:",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	client := protos.NewKVStoreClient(conn)
	resp, err := client.Set(context.Background(), &protos.SetRequest{Key: "a", Value: "1"})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("unexpected response: %v", resp.Message)
	}
}
