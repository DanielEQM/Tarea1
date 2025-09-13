package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	pb "franklin/proto/franklin-sys/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMissionServer
}

func connectWithRetry(uri string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	const maxRetries = 10
	const delay = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(uri)
		if err == nil {
			log.Println("[*] Conexión exitosa")
			return conn, nil
		}
		log.Printf("Error en conexión (intendo %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(delay)
	}
	return nil, err
}

func fallo(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (s *server) Distraccion(ctx context.Context, in *pb.DistraccionRequest) (*pb.DistraccionResponse, error) {
	log.Printf("Franklin comienza la distracción por %d turnos...", in.GetTurnos())
	prob := rand.Intn(101)
	log.Printf("la probabilidad es: %d", prob)
	if prob < 10 {
		log.Printf("Chop, para de ladrar!")
		log.Printf("Oh no, fallamos")
		return &pb.DistraccionResponse{Confirmacion: false, Razon: "Chop me distrajo"}, nil
	}
	log.Print("Distracción completada con éxito por Franklin")
	log.Println("=== Fase 2 completada con éxito ===")
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
	//
	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := connectWithRetry(amqpURI)
	fallo(err, "Se excedio el tiempo")
	defer conn.Close()

	ch, err := conn.Channel()
	fallo(err, "Fallo al abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare("Franklin", false, false, false, false, nil)
	fallo(err, "Fallo al declarar la cola")
	//
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	fallo(err, "Fallo al registrar consumidor")

	var estrellas int = 0
	limite := 5
	extra := 0
	victoria := true
	razon := ""
	for i := 1; i < int(in.GetTurnos()); i++ {
		go func() {
			for d := range msgs {
				new := string(d.Body)
				estrellas, _ = strconv.Atoi(new)
				log.Printf("Ahora tengo %d estrellas", estrellas)
			}
		}()
		time.Sleep(50 * time.Millisecond)
		if estrellas == limite {
			victoria = false
			razon = "Se llego a 5 estrellas..."
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
	fallo(err, "Fallo al escuchar")

	s := grpc.NewServer()
	pb.RegisterMissionServer(s, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())
	fallo(err, "Fallo al servir")
}
