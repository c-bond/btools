package main

import (
	"log"
	"os"
	"time"

	pb "btools/protoc"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	//address = "192.168.0.44:40"
	address     = "0.0.0.0:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBtoolsManagerClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.DdsCheckClientDoc(ctx,
		&pb.DdsCheckRequest{Name: name, Projno: 12345, Docno: 33999999})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if !r.Res {
		log.Println("Doc doesn't exist")
	} else {
		log.Println("Doc exists")
	}

	t, err := c.DdsCheckContractorDoc(ctx,
		&pb.DdsCheckRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", t.Message)
}
