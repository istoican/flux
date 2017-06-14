package flux

import (
	"net"

	"github.com/istoican/flux/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server :
func Server() error {
	lis, err := net.Listen("tcp", ":2323")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterRPCServer(server, new(RPCServer))
	go server.Serve(lis)
}

func (flux *Flux) Put(ctx context.Context, in *pb.Req) (*pb.Void, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (flux *Flux) Get(key *pb.Key, stream *pb.Value) error {
	node := flux.peers.Get(key.Id)

	if v, err := flux.config.datastore.Get(key.Id); err != nil {
		return err
	}
	stream.Send(v)
	
	if node.Address == 
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}
