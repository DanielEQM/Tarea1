package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	pb "lester/proto/lester-sys/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMissionServer
}

func (s *server) Oferta(ctx context.Context, in *pb.MissionRequest) (*pb.MissionResponse, error) {
	log.Printf("Recibida petición con %d rechazos", in.GetRechazo())
	if in.GetRechazo() == 3 {
		log.Printf("Ahora te esperas")
		time.Sleep(10 * time.Second)
	}

	prob := rand.Intn(100)
	log.Printf("La probabilidad es de %d", prob)
	if prob > 90 {
		return &pb.MissionResponse{Disp: false, Botin: "0", ProbF: "0", ProbT: "0", Riesgo: "0"}, nil
	}

	return &pb.MissionResponse{Disp: true, Botin: "100", ProbF: "70", ProbT: "60", Riesgo: "70"}, nil
}

/*********************
** Nombre: Oferta
**********************
** Parametros: ctx (context.Context), in (*pb.MissionRequest)
**********************
** Retorno: *pb.MissionResponse, error
**********************
** Descripción: Parte de la fase 1, donde Lester le da a Michael las ofertas.
 Si se rechaza 3 veces se hace esperar por 10 segundos, y si se acepta
 retorna el botin, probabilidad de Franklin, probabilidad de Lester y el Riesgo.
*/

func (s *server) ConfirmMission(ctx context.Context, in *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	// debe hacer algo
	return &pb.ConfirmResponse{}, nil
}

/*********************
** Nombre: confirmMission
**********************
** Parametros: ctx (context.Context), in (*pb.ConfirmRequest)
**********************
** Retorno: *pb.ConfirmResponse, error
**********************
** Descripción: Parte de la fase 1. Si Michael confirma una misión
Lester se entera a través de *pb.ConfirmResponse.
*/

func main() {
	lis, err := net.Listen("tcp", ":50051")
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
