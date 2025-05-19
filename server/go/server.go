package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	pb "progetto/server/go/proto"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer
type server struct {
	pb.UnimplementedSensorServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SendData(ctx context.Context, in *pb.SensorData) (*pb.Response, error) {
	fmt.Println("Ricevuti dati dal sensore: ", in)

	response := &pb.Response{
		Message: "dati ricevuti correttamente",
		Success: true,
	}

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Println("ERROR: failed to listen: ", err)
	}

	s := grpc.NewServer()
	pb.RegisterSensorServiceServer(s, &server{})
	fmt.Println("Server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Println("ERROR: failed to serve: ", err)
	}

}
