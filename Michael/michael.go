package main

import (
	"context"
	"log"
	"os"
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
	c := pb.NewGreeterClient(conn)

	// 3. Preparamos el contexto y los datos para la llamada remota.
	//    Un contexto puede llevar deadlines, cancelaciones, y otros valores a través
	//    de las llamadas.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Tomamos un nombre de los argumentos de la línea de comandos, o usamos "Mundo"
	pregunta := "Hola"
	if len(os.Args) > 1 {
		pregunta = os.Args[1]
	}

	// 4. ¡Llamamos a la función remota!
	//    Esto parece una llamada a una función local, pero gRPC se encarga de
	//    serializar los datos, enviarlos al servidor, esperar la respuesta y
	//    deserializarla.
	r, err := c.GiveMission(ctx, &pb.MissionRequest{Pregunta: pregunta})
	if err != nil {
		log.Fatalf("No se pudo saludar: %v", err)
	}
	// 5. Imprimimos la respuesta del servidor.

	log.Printf("Respuesta del servidor: %s", r.GetHay())
	log.Printf("          del servidor: %s", r.GetBotin())
	log.Printf("          del servidor: %s", r.GetProbF())
	log.Printf("          del servidor: %s", r.GetProbT())
	log.Printf("          del servidor: %s", r.GetRiesgo())

	r1, err := c.ConfirmMission(ctx, &pb.MissionRequest{Pregunta: "HEy hey hey"})
	if err != nil {
		log.Fatalf("No se pudo saludar: %v", err)
	}

	log.Printf("Respuesta del servidor: %s", r1.GetConf())
}
