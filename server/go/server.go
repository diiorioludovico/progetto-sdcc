package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"progetto/server/go/menu"
	pb "progetto/server/go/proto"

	_ "github.com/go-sql-driver/mysql"

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

	conn := "root:root@tcp(localhost:3306)/edgedb"
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("ERROR: problem in opening db connection: ", err)
	}
	defer db.Close()

	//verifica della connessione
	if err := db.Ping(); err != nil {
		fmt.Println("ERROR: ping error: ", err)
	} else {
		fmt.Println("INFO: successful connection to db")
	}

	go menu.ShowMenu(db)
	fmt.Println("Server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Println("ERROR: failed to serve: ", err)
	}
}
