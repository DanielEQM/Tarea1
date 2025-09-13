package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	pb "trevor/proto/trevor-sys/proto"

	amqp "github.com/rabbitmq/amqp091-go"

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
		log.Printf("Error en conexión (intento %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(delay)
	}
	return nil, err
}

/*********************
** Nombre: connectWithRetry
**********************
** Parametros: uri (string)
**********************
** Retorno: *amqp.Connection, error
**********************
** Descripción: Intenta realizar una conexion a amqp hasta alcanzar un maximo de intentos, con un delay determinado entre cada uno. Retorna la tupla conexion, nil en caso de conexion exitosa, o la tupla nil, err cuando no la logra luego del maximo de intentos.
 */

func fallo(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

/*********************
** Nombre: fallo
**********************
** Parametros: err (error), msg (string)
**********************
** Retorno:
**********************
** Descripción: En caso de recibir un error como parametro, permite imprimirlo junto al mensaje recibido como parametro, para luego finalizar el programa.
 */

func (s *server) Distraccion(ctx context.Context, in *pb.DistraccionRequest) (*pb.DistraccionResponse, error) {
	log.Printf("Trevor comienza la distracción por %d turnos", in.GetTurnos())
	prob := rand.Intn(101)
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
	//
	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = "amqp://guest:guest@10.53.168.35:5672/"
	}
	conn, err := connectWithRetry(amqpURI)
	fallo(err, "Se excedio el tiempo")
	defer conn.Close()

	ch, err := conn.Channel()
	fallo(err, "No se pudo abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare("Trevor", false, false, false, false, nil)
	fallo(err, "No se pudo declarar la cola")

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
	fallo(err, "No se pudo registrar consumidor")

	var estrellas int = 0
	limite := 5
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
	lis, err := net.Listen("tcp", ":50053")
	fallo(err, "Fallo al escuchar")

	s := grpc.NewServer()
	pb.RegisterMissionServer(s, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}

}
