package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"

	"github.com/rabbitmq/amqp091-go"

	pb "trevor/proto/trevor-sys/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMissionServer
}

func fallo(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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

/*********************
** Nombre: Distraccion
**********************
** Parametros: ctx (context.Context), in (*pb.DistraccionRequest)
**********************
** Retorno: *pb.DistraccionResponse, error
**********************
** Descripción: Parte de la fase 2, donde Trevor lleva a cabo la distracción.
 Tiene un 10% de que falle porque se emborrachó, asignando el valor false a
 Confirmacion. Caso contrario, se asigna el valor true a Confirmacion.
*/

func (s *server) Golpe(stx context.Context, in *pb.GolpeRequest) (*pb.GolpeResponse, error) {
	var estrellas int = 0
	limite := 5
	victoria := true
	razon := ""
	for i := 1; i < int(in.GetTurnos()); i++ {
		log.Printf("Turno: %d", i)
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

/*********************
** Nombre: Golpe
**********************
** Parametros: ctx (context.Context), in (*pb.GolpeRequest)
**********************
** Retorno: *pb.GolpeResponse, error
**********************
** Descripción: Parte de la fase 3, donde Trevor lleva a cabo el golpe.
En el caso de obtener más de 5 estrellas, activa su habilidad especial de aumentar
su límite de fracaso a 7 estrellas. Retorna la confirmacion de victoria
o derrota, el botin extra y la razon asociada.
*/

func main() {
	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = "amqp://guest:guest@rabbitmq:5672/"
	}

	conn, err := amqp091.Dial(amqpURI)
	fallo(err, "No se pudo conectar a RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	fallo(err, "No se pudo abrir un canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cola.Trevor",
		false, false, false, false, nil,
	)
	fallo(err, "No se pudo declarar la cola")

	err = ch.QueueBind(
		q.Name,
		"Trevor",
		"notificaciones",
		false,
		nil,
	)
	fallo(err, "No se pudo hacer el bind")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	fallo(err, "No se pudo registrar el consumidor")

	log.Printf("Trevor esperando notificaciones...\n")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Trevor recibió: %s", d.Body)
		}
	}()
	go func() {
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("Fallo al escuchar: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterMissionServer(s, &server{})
		log.Printf("Servidor escuchando en %v", lis.Addr())

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Fallo al servir: %v", err)
		}
	}()
	<-forever
}
