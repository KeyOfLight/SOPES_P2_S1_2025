package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	pb "api-rest/proto" // Asegúrate que el path sea correcto

	"google.golang.org/grpc"
)

type WeatherData struct {
	Description string `json:"Description"`
	Country     string `json:"Country"`
	Weather     string `json:"Weather"`
}

var grpcClient pb.WeatherServiceClient

func handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo el cuerpo de la petición", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Printf("📥 Datos recibidos en API REST: %+v\n", data)

	req := &pb.WeatherRequest{
		Description: data.Description,
		Country:     data.Country,
		Weather:     data.Weather,
	}

	resp, err := grpcClient.SendWeatherData(ctx, req)
	if err != nil {
		http.Error(w, "Error al contactar gRPC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Respuesta del gRPC: %s", resp.Message)
}

func main() {
	// Conectar al servidor gRPC en el namespace grpc-namespace
	conn, err := grpc.Dial("grpc-client-service.grpc-namespace.svc.cluster.local:50052", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("No se pudo conectar al cliente gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient = pb.NewWeatherServiceClient(conn)

	http.HandleFunc("/input", handlePost)

	log.Println("🌐 API REST escuchando en :8081 y enviando en 50052 ...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
