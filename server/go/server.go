package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"progetto/server/go/fejs"
	"progetto/server/go/menu"
	pb "progetto/server/go/proto"
	qr "progetto/server/go/query"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedSensorServiceServer
	db *sql.DB
}

func (s *server) SendData(ctx context.Context, in *pb.SensorData) (*pb.Response, error) {
	//fmt.Println("Ricevuti dati dal sensore: ", in)
	_, err := s.db.Exec(qr.GetMeasure(), in.DeviceID, in.ParkID, in.Temperature, in.Humidity, in.Brightness, in.AirQuality, in.Timestamp)
	if err != nil {
		fmt.Println("ERROR: error in inserting new measure: ", err)
	}

	response := &pb.Response{
		Message: "dati ricevuti correttamente",
		Success: true,
	}

	return response, nil
}

func (s *server) Configuration(ctx context.Context, in *pb.SensorIdentification) (*pb.CommunicationConfiguration, error) {
	//fmt.Println("Ricevuto primo feedback dal sensore ", in.SerialNumber)

	//recupero del record associato al sensore e preparazione dellla response
	rows, err := s.db.Query(qr.GetSensor(), in.SerialNumber)
	if err != nil {
		fmt.Println("ERROR: query error: ", err)
	}

	var sensor menu.Sensor
	var count int

	for rows.Next() {
		if err := rows.Scan(&sensor.Id, &sensor.Is_active, &sensor.Park_id, &sensor.Serial_number); err != nil {
			fmt.Println("ERROR: scan error: ", err)
		}
		count += 1
	}

	var response *pb.CommunicationConfiguration

	if err := rows.Err(); err != nil {
		fmt.Println("ERROR: rows error: ", err)
	} else if count == 0 {
		fmt.Println("INFO: no record found")
		return response, nil
	}

	response = &pb.CommunicationConfiguration{
		DeviceID: strconv.Itoa(sensor.Id),
		ParkID:   strconv.FormatInt(sensor.Park_id.Int64, 10),
		Interval: 10,
	}

	//modifica del valore is_active del sensore per indicare che è attivo nel parco e può iniziare ad inviare dati
	res, err := s.db.Exec(qr.UpdateSensorStatus(), true, sensor.Id)
	if err != nil {
		fmt.Println("ERROR: update error: ", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println("ERROR: rows affected error: ", err)
	}

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Println("ERROR: failed to listen: ", err)
	}

	conn := "root:root@tcp(localhost:3306)/edgedb"
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("ERROR: problem in opening db connection: ", err)
	}
	defer db.Close()

	s := grpc.NewServer()
	pb.RegisterSensorServiceServer(s, &server{db: db})

	//verifica della connessione
	if err := db.Ping(); err != nil {
		fmt.Println("ERROR: ping error: ", err)
	} else {
		fmt.Println("INFO: successful connection to db")
	}

	//goroutine per il menu
	//go menu.ShowMenu(db)

	//goroutine per mostrare servizio per backend
	go fejs.StartFrontendSetup(db)

	fmt.Println("Server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Println("ERROR: failed to serve: ", err)
	}
}
