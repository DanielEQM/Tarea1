package main

import (
	"context"
	"log"
	"time"

	pb "Michael/proto/michael-sys/proto" // Reemplaza con el path de tu módulo

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "lester-container:50051" // La dirección del servidor
)

func main() {
	// 1. Establecemos una conexión con el servidor.
	//    Usamos WithTransportCredentials(insecure.NewCredentials()) porque no estamos
	//    usando SSL/TLS en este ejemplo simple.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close() // Importante cerrar la conexión al final.

	// 2. Creamos un "stub" de cliente a partir de la conexión.
	//    Este objeto 'c' es el que tiene los métodos remotos que podemos llamar.
	c := pb.NewMissionClient(conn)

	// 3. Preparamos el contexto y los datos para la llamada remota.
	//    Un contexto puede llevar deadlines, cancelaciones, y otros valores a través
	//    de las llamadas.
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	// Tomamos un nombre de los argumentos de la línea de comandos, o usamos "Mundo"

	// 4. ¡Llamamos a la función remota!
	//    Esto parece una llamada a una función local, pero gRPC se encarga de
	//    serializar los datos, enviarlos al servidor, esperar la respuesta y
	//    deserializarla.
	var rechazo int32 = 0
	keep := true
	for keep {
		r, err := c.Oferta(ctx, &pb.MissionRequest{Rechazo: rechazo})
		if rechazo == 3 {
			rechazo = 0
		}
		if err != nil {
			log.Fatalf("No se pudo saludar: %v", err)
		}
		if r.GetDisp() {
			log.Printf("y el rechazo %d", rechazo)
			log.Printf("Respuesta del servidor: %t", r.GetDisp())
			log.Printf("          del servidor: %d", r.GetBotin())
			log.Printf("          del servidor: %d", r.GetProbF())
			log.Printf("          del servidor: %d", r.GetProbT())
			log.Printf("          del servidor: %d", r.GetRiesgo())
			log.Printf("")
			if r.GetRiesgo() == 0 || !(r.GetRiesgo() < 80 && (r.GetProbF() > 50 || r.GetProbT() > 50)) {
				rechazo += 1
			} else {
				break
			}
		}
	}
	// 5. Imprimimos la respuesta del servidor.

	_, err = c.ConfirmMission(ctx, &pb.ConfirmRequest{Conf: true})
	if err != nil {
		log.Fatalf("No se pudo saludar: %v", err)
	}
}
