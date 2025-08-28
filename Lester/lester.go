package main

import (
	"context"
	"log"
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

// SayHello es la implementación de la función definida en el archivo .proto.
// Esta es la lógica real que se ejecuta cuando un cliente llama a este RPC.
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Recibida petición de: %v", in.GetName())
	// Creamos y devolvemos la respuesta.
	return &pb.HelloResponse{Message: "Hola, " + in.GetName() + "!"}, nil
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
