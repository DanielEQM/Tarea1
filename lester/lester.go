package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/streadway/amqp"

	pb "lester/proto/lester-sys/proto"

	"google.golang.org/grpc"
)

var rabbitConn *amqp.Connection
var rabbitCh *amqp.Channel
var sss [][]string

func notificarEstrellas(personaje string, riesgo int, stopChan chan bool) {
	frecuencia := 100 - riesgo
	estrellas := 0

	ticker := time.NewTicker(time.Duration(frecuencia) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			log.Printf("Se detuvieron las notificaciones a %s", personaje)
			return
		case <-ticker.C:
			estrellas++
			msg := fmt.Sprintf("%s tiene ahora %d estrellas", personaje, estrellas)
			publishMessage(personaje, msg)

			if personaje == "Trevor" && estrellas >= 7 {
				log.Printf("Trevor alcanzó su limite de estrellas (7). Fracaso")
				stopChan <- true
				return
			}
			if personaje == "Franklin" && estrellas >= 5 {
				log.Printf("Franklin alcanzó su limite de estrellas (5). Fracasó")
				stopChan <- true
				return
			}
		}
	}
}

func fallo(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func initRabbit() {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	fallo(err, "No se pudo conectar a RabbitMQ")

	rabbitCh, err = rabbitConn.Channel()
	fallo(err, "No se pudo abrir el canal")

	_, err = rabbitCh.QueueDeclare(
		"notificaciones",
		false,
		false,
		false,
		false,
		nil,
	)
	fallo(err, "No se pudo declarar la cola")
}

func publishMessage(personaje, msg string) {
	err := rabbitCh.Publish(
		"notificaciones",
		personaje,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	fallo(err, "Error publicando un mensaje")
	log.Printf(" [x] Lester notificó a %s: %s", personaje, msg)
}

type server struct {
	pb.UnimplementedMissionServer
}

func (s *server) Oferta(ctx context.Context, in *pb.MissionRequest) (*pb.MissionResponse, error) {
	log.Printf("Recibida petición con %d rechazos", in.GetRechazo())
	if in.GetRechazo() == 3 {
		log.Printf("Ahora te esperas")
		time.Sleep(10 * time.Second)
	}

	prob := rand.Intn(101)
	log.Printf("La probabilidad es de %d", prob)
	o := sss[rand.Intn(len(sss))]
	if prob > 90 {
		return &pb.MissionResponse{Disp: false, Botin: "0", ProbF: "0", ProbT: "0", Riesgo: "0"}, nil
	}

	return &pb.MissionResponse{Disp: true, Botin: o[0], ProbF: o[1], ProbT: o[2], Riesgo: o[3]}, nil
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
	file, err := os.Open("ofertas/ofertas_grande.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	cont := 0
	line, err := reader.ReadString('\n')
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("%v", err)
		}
		sep := 0
		aux := ""
		sss = append(sss, []string{})
		for i := 0; i < len(line); i++ {
			if line[i] == ',' {
				sss[cont] = append(sss[cont], aux)
				sep++
				aux = ""
				continue
			}
			aux += string(line[i])
		}
		sss[cont] = append(sss[cont], aux[:len(aux)-2])
		cont++
	}

	stopChan := make(chan bool)
	go notificarEstrellas("Franklin", 20, stopChan)
	go notificarEstrellas("Trevor", 20, stopChan)

	initRabbit()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMissionServer(s, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())

	defer rabbitConn.Close()
	defer rabbitCh.Close()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
