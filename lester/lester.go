package main

import (
	"bufio"
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	pb "lester/proto/lester-sys/proto"

	amqp "github.com/rabbitmq/amqp091-go"

	"google.golang.org/grpc"
)

var conn *amqp.Connection
var ch *amqp.Channel
var flag bool = true
var sss [][]string

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

/*********************
** Nombre: connectWithRetry
**********************
** Parametros: uri (string)
**********************
** Retorno: *amqp.Connection, error
**********************
** Descripción: Intenta realizar una conexion a amqp hasta alcanzar un maximo de intentos, con un delay determinado entre cada uno. Retorna la tupla conexion, nil en caso de conexion exitosa, o la tupla nil, err cuando no la logra luego del maximo de intentos.
 */

func (s *server) NotificarGolpe(ctx context.Context, in *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	flag = in.GetConf()
	return &pb.ConfirmResponse{}, nil
}

/*********************
** Nombre: NotificarGolpe
**********************
** Parametros: ctx (context.Context), in (*pb.ConfirmRequest)
**********************
** Retorno: *pb.ConfirmResponse, error
**********************
** Descripción:
 */

func (s *server) NotificarEstrellas(ctx context.Context, in *pb.AvisoRequest) (*pb.AvisoResponse, error) {
	go func() {
		amqpURI := os.Getenv("AMQP_URI")
		if amqpURI == "" {
			amqpURI = "amqp://guest:guest@localhost:5672/"
		}
		conn, err := connectWithRetry(amqpURI)
		fallo(err, "Se excedió el tiempo maximo")
		defer conn.Close()

		ch, err = conn.Channel()
		fallo(err, "No se puede abrir el canal")
		defer ch.Close()

		q, err := ch.QueueDeclare(in.GetPj(), false, false, false, false, nil)
		fallo(err, "No se puede declarar la cola")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		frecuencia := 100 - in.GetRiesgo()
		estrellas := 1

		for i := 0; i < int(in.GetTurnos()); i++ {
			if !flag {
				break
			}
			if i == int(frecuencia) {
				log.Printf("Aumento a %d estrellas!", estrellas)
				body := strconv.Itoa(estrellas)
				err = ch.PublishWithContext(ctx,
					"",
					q.Name,
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(body),
					})
				fallo(err, "No se pudo publicar el mensaje.")
				frecuencia += 100 - in.GetRiesgo()
				estrellas++
			}
			time.Sleep(50 * time.Millisecond)
		}

	}()
	return &pb.AvisoResponse{}, nil
}

/*********************
** Nombre: NotificarEstrellas
**********************
** Parametros: ctx (context.Context), in (*pb.AvisoRequest)
**********************
** Retorno: *pb.AvisoResponse, error
**********************
** Descripción: Genera la mensajeria asincronica que permite notificar a Trevor y Franklin de sus estrellas. Se conecta a RabbitMQ, abre un canal, declara una cola y envia un mensaje cuando frecuencia coincide con la iteracion de los turnos de la mision.
 */

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
	log.Printf("Entendido, estaré atento.")
	return &pb.ConfirmResponse{}, nil
}

/*********************
** Nombre: ConfirmMission
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
	fallo(err, "No se pudo cargar el archivo.")
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

	// Servidor
	lis, err := net.Listen("tcp", ":50051")
	fallo(err, "Fallo al escuchar.")
	s := grpc.NewServer()
	pb.RegisterMissionServer(s, &server{})
	log.Printf("Servidor escuchando en %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
