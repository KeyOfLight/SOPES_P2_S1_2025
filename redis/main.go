package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "redis/proto" // Asegúrate de que esto apunta a tu paquete generado

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

type server struct {
	pb.UnimplementedWeatherServiceServer
}

// Esta función maneja las solicitudes gRPC y guarda en Redis
func (s *server) SendWeatherData(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	log.Printf("📥 Datos recibidos: %s, %s, %s\n", req.Description, req.Country, req.Weather)

	timestamp := time.Now().Format("20060102_150405")
	key := fmt.Sprintf("weather:%s", timestamp)

	value := fmt.Sprintf("description: %s | country: %s | weather: %s", req.Description, req.Country, req.Weather)

	if err := rdb.Set(ctx, key, value, 0).Err(); err != nil {
		log.Printf("❌ Error al guardar en Redis: %v", err)
	} else {
		log.Printf("✅ Guardado en Redis [%s]: %s", key, value)
	}

	msg := fmt.Sprintf("Datos guardados en Redis: %s", key)
	return &pb.WeatherResponse{Message: msg}, nil
}

func main() {
	// Configuración de Redis desde variables de entorno
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	// Inicializar el servidor gRPC
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ No se pudo escuchar en el puerto 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpcServer, &server{})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ No se pudo conectar a Redis: %v", err)
	}
	log.Println("✅ Conexión exitosa a Redis")

	log.Println("🚀 Servidor gRPC escuchando en el puerto 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor gRPC: %v", err)
	}
}
