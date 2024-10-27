package main

import (
	"context"
	"log"
	"time"

	pb "github.com/danile0SA/0250952_SistemasDistribuidos/api/v1" // Importa el código generado
	"google.golang.org/grpc"
)

func main2() {
	// Conectar al servidor gRPC
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewLogClient(conn)

	// Producir un record
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	produceRes, err := client.Produce(ctx, &pb.ProduceRequest{
		Record: &pb.Record{Value: []byte("hello gRPC")},
	})
	if err != nil {
		log.Fatalf("Error en Produce: %v", err)
	}
	log.Printf("Record producido con offset: %d", produceRes.Offset)

	// Consumir el record
	consumeRes, err := client.Consume(ctx, &pb.ConsumeRequest{Offset: produceRes.Offset})
	if err != nil {
		log.Fatalf("Error en Consume: %v", err)
	}
	log.Printf("Record consumido: %s", string(consumeRes.Record.Value))
}
