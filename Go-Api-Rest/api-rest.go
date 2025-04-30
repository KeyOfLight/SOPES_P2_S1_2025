package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "go-deployment-1/proto" // Ajusta al path correcto

	"google.golang.org/grpc"
)

type WeatherData struct {
	Description string `json:"Description"`
	Country     string `json:"Country"`
	Weather     string `json:"Weather"`
}

var grpcClient pb.WeatherServiceClient

func handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo el cuerpo", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	req := &pb.WeatherRequest{
		Description: data.Description,
		Country:     data.Country,
		Weather:     data.Weather,
	}

	resp, err := grpcClient.SendWeatherData(context.Background(), req)
	if err != nil {
		http.Error(w, "Error al contactar gRPC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Respuesta del gRPC: %s", resp.Message)
}

func main() {
	// Conectar al servidor gRPC
	conn, err := grpc.Dial("grpc-server-service:50051", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient = pb.NewWeatherServiceClient(conn)

	http.HandleFunc("/input", handlePost)

	log.Println("API REST escuchando en :8081 ...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
