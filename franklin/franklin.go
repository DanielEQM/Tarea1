package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	// Importamos el código generado por protoc
	pb "franklin/proto/franklin-sys/proto" // Reemplaza con el path de tu módulo

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMissionServer
}

func (s *server) Distraccion(ctx context.Context, in *pb.DistraccionRequest) (*pb.DistraccionResponse, error) {
	prob := rand.Intn(100)
	log.Printf("la probabilidad es: %d", prob)
	if prob < 10 {
		log.Printf("Chop, para de ladrar!")
		log.Printf("Oh no, fallamos")
		return &pb.DistraccionResponse{Confirmacion: false}, nil
	}
	return &pb.DistraccionResponse{Confirmacion: true}, nil
}

func main() {
	// 1. Abrimos un puerto para escuchar (en este caso, el 50051)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	// 2. Creamos una nueva instancia del servidor gRPC
	s := grpc.NewServer()

	// 3. Registramos nuestro servicio 'Greeter' en el servidor gRPC.
	//    Esto conecta nuestra implementación lógica (la struct 'server') con el
	//    servicio definido en el .proto.
	pb.RegisterMissionServer(s, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())

	// 4. Iniciamos el servidor para que empiece a aceptar peticiones en el puerto.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
