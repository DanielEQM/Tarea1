package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	// Importamos el código generado por protoc
	pb "Lester/proto/lester-sys/proto" // Reemplaza con el path de tu módulo

	"google.golang.org/grpc"
)

// Definimos una struct para nuestro servidor. Debe embeber el UnimplementedGreeterServer.
// Esto asegura la compatibilidad hacia adelante si se añaden más RPCs al servicio.
type server struct {
	pb.UnimplementedGreeterServer
}

// var rechazo = 0

// SayHello es la implementación de la función definida en el archivo .proto.
// Esta es la lógica real que se ejecuta cuando un cliente llama a este RPC.
func (s *server) GiveMission(ctx context.Context, in *pb.MissionRequest) (*pb.MissionResponse, error) {
	log.Printf("Recibida petición de: %v", in.GetPregunta())
	// Creamos y devolvemos la respuesta.
	prob := rand.Intn(100)
	if prob >= 90 {
		return &pb.MissionResponse{Hay: "NO", Botin: "", ProbF: "", ProbT: "", Riesgo: ""}, nil
	}
	return &pb.MissionResponse{Hay: "YES", Botin: "1", ProbF: "2", ProbT: "3", Riesgo: "4"}, nil

}

func main() {
	// 1. Abrimos un puerto para escuchar (en este caso, el 50051)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	// 2. Creamos una nueva instancia del servidor gRPC
	s := grpc.NewServer()

	// 3. Registramos nuestro servicio 'Greeter' en el servidor gRPC.
	//    Esto conecta nuestra implementación lógica (la struct 'server') con el
	//    servicio definido en el .proto.
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())

	// 4. Iniciamos el servidor para que empiece a aceptar peticiones en el puerto.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
