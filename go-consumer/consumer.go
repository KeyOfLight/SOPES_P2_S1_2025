package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var ctx = context.Background()

func main() {

	topic := "my.topic"
	broker := "my-cluster-kafka-bootstrap.kafka:9092"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "go-consumer-group",
	})

	// redisAddr := "redis-kafka-service.kafka.svc.cluster.local:6379"
	// redisPassword := "password_sopes_2025"
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     redisAddr,
	// 	Password: redisPassword,
	// 	DB:       0,
	// })

	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis-service.kafka.svc.cluster.local:6379",
		// No Password field
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("No se pudo conectar a Redis: %v", err)
	}
	fmt.Println("Conectado a Redis correctamente.")

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Error leyendo mensaje:", err)
			break
		}
		msg := string(m.Value)
		fmt.Printf("Mensaje recibido: %s\n", msg)

		key := fmt.Sprintf("kafka:%s:offset:%d", topic, m.Offset)
		err = redisClient.Set(ctx, key, msg, 0).Err()
		if err != nil {
			fmt.Printf("Error al guardar en Redis: %v\n", err)
		} else {
			fmt.Printf("Mensaje guardado en Redis con clave %s\n", key)
		}
	}

	r.Close()
}
