package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "grpc-client/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *server) SendWeatherData(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	// Aquí se reciben los datos enviados desde la API REST
	log.Printf("✅ Datos recibidos en grpc-client:\n📌 Description: %s\n📌 Country: %s\n📌 Weather: %s\n",
		req.Description, req.Country, req.Weather)
	log.Printf("Mensaje recibido en grpc-client: %s en %s (%s)", req.Description, req.Country, req.Weather)

	conn, err := grpc.Dial("grpc-server:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("❌ Error conectando con grpc-server: %v", err)
		return nil, err
	}
	defer conn.Close()

	grpcServerClient := pb.NewWeatherServiceClient(conn)

	ctx2, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	forwardReq := &pb.WeatherRequest{
		Description: req.Description,
		Country:     req.Country,
		Weather:     req.Weather,
	}

	resp, err := grpcServerClient.SendWeatherData(ctx2, forwardReq)
	if err != nil {
		log.Printf("❌ Error reenviando al grpc-server: %v", err)
		return nil, err
	}

	return &pb.WeatherResponse{Message: "🔁 Reenviado al grpc-server: " + resp.Message}, nil
}

func main() {

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ No se pudo escuchar en el puerto 50052: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWeatherServiceServer(s, &server{})

	log.Println("🚀 grpc-client escuchando en el puerto 50052...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("❌ Fallo al iniciar el grpc-client: %v", err)
	}
}
