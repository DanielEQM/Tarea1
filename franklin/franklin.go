package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "franklin/proto/franklin-sys/proto"

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
		return &pb.DistraccionResponse{Confirmacion: false, Razon: "Chop me distrajo"}, nil
	}
	return &pb.DistraccionResponse{Confirmacion: true}, nil
}

/*********************
** Nombre: Distraccion
**********************
** Parametros: ctx (context.Context), in (*pb.DistraccionRequest)
**********************
** Retorno: *pb.DistraccionResponse, error
**********************
** Descripción: Parte de la fase 2, donde Franklin lleva a cabo la distracción.
 Tiene un 10% de que falle porque su perro Chop ladre, asignando el valor false a
 Confirmacion. Caso contrario, se asigna el valor true a Confirmacion.
*/

func (s *server) Golpe(ctx context.Context, in *pb.GolpeRequest) (*pb.GolpeResponse, error) {
	var estrellas int = 0
	limite := 5
	extra := 0
	victoria := true
	razon := ""
	for i := 1; i < int(in.GetTurnos()); i++ {
		if estrellas == limite {
			victoria = false
			razon = "Se llego a 5 estrellas..."
			extra = 0
			break
		}
		if estrellas >= 3 {
			extra += 1000
		}
	}
	return &pb.GolpeResponse{Confirmacion: victoria, BotinExtra: int32(extra), Razon: razon}, nil
}

/*********************
** Nombre: Golpe
**********************
** Parametros: ctx (context.Context), in (*pb.GolpeRequest)
**********************
** Retorno: *pb.GolpeResponse, error
**********************
** Descripción: Parte de la fase 3, donde Franklin lleva a cabo el golpe.
En el caso de obtener 3 estrellas activa su habilidad especial de sumar $1000
a la variable extra por cada turno que pase. Retorna la confirmacion de victoria
o derrota, el botin extra y la razon asociada.
*/

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMissionServer(s, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
