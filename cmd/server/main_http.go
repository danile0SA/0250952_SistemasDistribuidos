/*package main_http

import (
	"log"

	server "github.com/danile0SA/0250952_SistemasDistribuidos/Internal/Server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}*/

package main

import (
	"context"
	"log"
	"net"

	pb "github.com/danile0SA/0250952_SistemasDistribuidos/api/v1" // Importa el código generado
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogService implementa el servicio Log definido en protobuf
type LogService struct {
	// Embebe la estructura que implementa los métodos no implementados
	pb.UnimplementedLogServer
	records []*pb.Record
}

// Produce recibe un record y lo añade al log
func (s *LogService) Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceResponse, error) {
	offset := uint64(len(s.records))
	s.records = append(s.records, &pb.Record{
		Value:  req.Record.Value,
		Offset: offset,
	})
	return &pb.ProduceResponse{Offset: offset}, nil
}

// Consume devuelve un record dado un offset
func (s *LogService) Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeResponse, error) {
	// Retorna el record en el offset solicitado
	if req.Offset >= uint64(len(s.records)) {
		return nil, status.Error(codes.Unavailable, "offset fuera de rango")
	}
	record := s.records[req.Offset]
	return &pb.ConsumeResponse{Record: record}, nil
}

// ConsumeStream devuelve un flujo de records comenzando desde el offset solicitado
func (s *LogService) ConsumeStream(req *pb.ConsumeRequest, stream pb.Log_ConsumeStreamServer) error {
	for i := req.Offset; i < uint64(len(s.records)); i++ {
		err := stream.Send(&pb.ConsumeResponse{
			Record: s.records[i],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ProduceStream recibe un flujo de records y los añade al log
func (s *LogService) ProduceStream(stream pb.Log_ProduceStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		offset := uint64(len(s.records))
		s.records = append(s.records, &pb.Record{
			Value:  req.Record.Value,
			Offset: offset,
		})
		if err := stream.Send(&pb.ProduceResponse{Offset: offset}); err != nil {
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLogServer(grpcServer, &LogService{}) // Se registra el servicio correctamente

	log.Println("Servidor escuchando en el puerto 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
