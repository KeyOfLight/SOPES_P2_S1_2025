package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "go-deployment-1/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *server) SendWeatherData(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	// Aquí iría la lógica para publicar a RabbitMQ o Kafka
	log.Printf("Recibido en gRPC Client Server: Description=%s, Country=%s, Weather=%s",
		req.Description, req.Country, req.Weather)

	// Aquí puedes llamar a los publicadores (según sea RabbitMQ o Kafka)
	// Por ahora solo responde
	msg := fmt.Sprintf("Mensaje recibido: %s en %s (%s)", req.Description, req.Country, req.Weather)
	return &pb.WeatherResponse{Message: msg}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo escuchar en el puerto 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWeatherServiceServer(s, &server{})

	log.Println("gRPC Client Server escuchando en el puerto 50051...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Fallo al iniciar el servidor gRPC: %v", err)
	}
}
