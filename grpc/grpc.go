package grpc

import (
	"context"
	"log"
	"net"

	db "btools/Db"
	pb "btools/protoc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement BtoolsManager Server.
type server struct{}

// DdsCheckClientDoc implements helloworld.GreeterServer
func (s *server) DdsCheckClientDoc(ctx context.Context,
	in *pb.DdsCheckRequest) (*pb.DdsCheckReply, error) {
	ret := db.DdsDoesClientDocExist(in.Projno, in.Docno)
	return &pb.DdsCheckReply{Message: "Hello " + in.Name, Res: ret}, nil
}

func (s *server) DdsCheckContractorDoc(ctx context.Context,
	in *pb.DdsCheckRequest) (*pb.DdsCheckReply, error) {
	return &pb.DdsCheckReply{Message: "Goodbyes " + in.Name}, nil
}

func StartServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterBtoolsManagerServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
