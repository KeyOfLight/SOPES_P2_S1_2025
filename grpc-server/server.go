package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-server/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *server) SendWeatherData(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	// Imprimir en consola cada solicitud recibida
	log.Printf("Recibido: Descripción=%s, País=%s, Clima=%s\n", req.Description, req.Country, req.Weather)

	// Crear mensaje de respuesta
	message := fmt.Sprintf("Datos recibidos correctamente: %s, %s, %s", req.Description, req.Country, req.Weather)
	return &pb.WeatherResponse{Message: message}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpcServer, &server{})

	log.Println("Servidor gRPC corriendo en el puerto 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error al ejecutar el servidor: %v", err)
	}
}
