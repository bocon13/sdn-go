package p4runtime

import (
	"github.com/bocon13/sdn-go/proto/p4"
	"google.golang.org/grpc"
)

// Cache of address to gRPC client
var clients map[string]*grpc.ClientConn = make(map[string]*grpc.ClientConn)

func GetClient(host string) p4.P4RuntimeClient {
	var err error
	conn, ok := clients[host]
	if !ok {
		conn, err = grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		clients[host] = conn
	}
	return p4.NewP4RuntimeClient(conn)
}
