package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	// Importamos el código generado por protoc
	pb "Lester/proto/lester-sys/proto" // Reemplaza con el path de tu módulo

	"google.golang.org/grpc"
)

// Definimos una struct para nuestro servidor. Debe embeber el UnimplementedGreeterServer.
// Esto asegura la compatibilidad hacia adelante si se añaden más RPCs al servicio.
type server struct {
	pb.UnimplementedMissionServer
}

// var rechazo = 0

// SayHello es la implementación de la función definida en el archivo .proto.
// Esta es la lógica real que se ejecuta cuando un cliente llama a este RPC.
func (s *server) Oferta(ctx context.Context, in *pb.MissionRequest) (*pb.MissionResponse, error) {
	log.Printf("Recibida petición con %d rechazos", in.GetRechazo())
	if in.GetRechazo() == 3 {
		log.Printf("Ahora te esperas")
		time.Sleep(10 * time.Second)
	}
	// Creamos y devolvemos la respuesta.
	prob := rand.Intn(100)
	log.Printf("La probabilidad es de %d", prob)
	if prob > 90 {
		return &pb.MissionResponse{Disp: false, Botin: 0, ProbF: 0, ProbT: 0, Riesgo: 0}, nil
	}

	return &pb.MissionResponse{Disp: true, Botin: 100, ProbF: 70, ProbT: 70, Riesgo: 70}, nil
}

func (s1 *server) ConfirmMission(ctx context.Context, in *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	// debe hacer algo
	return &pb.ConfirmResponse{}, nil
}

func main() {
	// 1. Abrimos un puerto para escuchar (en este caso, el 50051)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	// 2. Creamos una nueva instancia del servidor gRPC
	s := grpc.NewServer()
	s1 := grpc.NewServer()

	// 3. Registramos nuestro servicio 'Greeter' en el servidor gRPC.
	//    Esto conecta nuestra implementación lógica (la struct 'server') con el
	//    servicio definido en el .proto.
	pb.RegisterMissionServer(s, &server{})
	pb.RegisterMissionServer(s1, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())

	// 4. Iniciamos el servidor para que empiece a aceptar peticiones en el puerto.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
