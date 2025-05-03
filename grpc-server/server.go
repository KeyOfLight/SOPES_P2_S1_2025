package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-server/proto"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *server) SendWeatherData(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	log.Printf("📬 grpc-server recibió:\n➡️  description: %s\n🌍 country: %s\n⛅ weather: %s",
		req.Description, req.Country, req.Weather)

	// responseMsg := fmt.Sprintf("✅ Datos recibidos correctamente en grpc-server: %s, %s, %s",
	// 	req.Description, req.Country, req.Weather)

	message := fmt.Sprintf("description: %s, country: %s, weather: %s", req.Description, req.Country, req.Weather)

	broker := "my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092" // Dirección del broker Kafka
	topic := "my.topic"

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(message), // Usamos el mensaje recibido como contenido para Kafka
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error al enviar mensaje a Kafka: %v", err)
	}

	respMessage := fmt.Sprintf("Mensaje enviado correctamente: %s", message)
	return &pb.WeatherResponse{Message: respMessage}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ No se pudo iniciar en puerto 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpcServer, &server{})

	log.Println("🚀 grpc-server escuchando en el puerto 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Fallo al ejecutar grpc-server: %v", err)
	}
}
