package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	// Importamos el código generado por protoc
	pb "trevor/proto/trevor-sys/proto" // Reemplaza con el path de tu módulo

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMissionServer
}

func (s *server) Distraccion(ctx context.Context, in *pb.DistraccionRequest) (*pb.DistraccionResponse, error) {
	prob := rand.Intn(100)
	log.Printf("la probabilidad es: %d", prob)
	if prob < 10 {
		log.Printf("Parece que bebi más de la cuenta... zzz")
		return &pb.DistraccionResponse{Confirmacion: false, Razon: "Bebi más de la cuenta"}, nil
	}
	return &pb.DistraccionResponse{Confirmacion: true}, nil
}

func (s *server) Golpe(stx context.Context, in *pb.GolpeRequest) (*pb.GolpeResponse, error) {
	var estrellas int = 0
	limite := 5
	victoria := true
	razon := ""
	for i := 1; i < int(in.GetTurnos()); i++ {
		log.Printf("Turno: %d", i)
		// Acá lee?
		if estrellas == 5 {
			limite = 7
		}
		if estrellas == limite {
			victoria = false
			razon = "Se llego a 7 estrellas..."
			break
		}
	}
	estrellas++
	return &pb.GolpeResponse{Confirmacion: victoria, Razon: razon}, nil
}

func main() {
	// 1. Abrimos un puerto para escuchar (en este caso, el 50051)
	lis, err := net.Listen("tcp", ":50053")
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
