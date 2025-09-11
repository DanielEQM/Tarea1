package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "michael/proto/michael-sys/proto" // Reemplaza con el path de tu módulo

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addressL = "localhost:50051"
	addressF = "localhost:50052"
	addressT = "localhost:50053" // La dirección del servidor
)

func main() {
	// 1. Establecemos una conexión con el servidor.
	//    Usamos WithTransportCredentials(insecure.NewCredentials()) porque no estamos
	//    usando SSL/TLS en este ejemplo simple.
	connL, err := grpc.Dial(addressL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer connL.Close() // Importante cerrar la conexión al final.

	connF, err := grpc.Dial(addressF, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer connF.Close()

	connT, err := grpc.Dial(addressT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer connT.Close()

	// 2. Creamos un "stub" de cliente a partir de la conexión.
	//    Este objeto 'c' es el que tiene los métodos remotos que podemos llamar.
	L := pb.NewMissionClient(connL)
	F := pb.NewMissionClient(connF)
	T := pb.NewMissionClient(connT)

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
	var r *pb.MissionResponse
	for {
		r, err = L.Oferta(ctx, &pb.MissionRequest{Rechazo: rechazo})
		if err != nil {
			log.Fatalf("No se pudo saludar: %v", err)
		}
		if rechazo == 3 {
			rechazo = 0
		}
		if r.GetDisp() {
			log.Printf("y el rechazo %d", rechazo)
			log.Printf("Respuesta del servidor: %t", r.GetDisp())
			log.Printf("          del servidor: %s", r.GetBotin())
			log.Printf("          del servidor: %s", r.GetProbF())
			log.Printf("          del servidor: %s", r.GetProbT())
			log.Printf("          del servidor: %s", r.GetRiesgo())
			log.Printf("")
			if r.GetBotin() == "" || r.GetProbF() == "" || r.GetProbT() == "" || r.GetRiesgo() == "" {
				rechazo++
				continue
			}
			probF, _ := strconv.Atoi(r.GetProbF())
			probT, _ := strconv.Atoi(r.GetProbT())
			riesgo, _ := strconv.Atoi(r.GetRiesgo())
			if (probF <= 50 && probT <= 50) || riesgo >= 80 {
				rechazo++
				continue
			}
			break
		}
	}

	_, err = L.ConfirmMission(ctx, &pb.ConfirmRequest{Conf: true})
	if err != nil {
		log.Fatalf("No se pudo saludar: %v", err)
	}
	//----- Se guardan los datos
	botin, _ := strconv.Atoi(r.GetBotin())
	probF, _ := strconv.Atoi(r.GetProbF())
	probT, _ := strconv.Atoi(r.GetProbT())
	//riesgo, _ := strconv.Atoi(r.GetRiesgo())
	log.Printf("probF: %d y probT: %d", probF, probT)
	//-----

	victoria := true
	if probF > probT {
		turnos := 200 - probF
		log.Printf("%d", turnos)
		f, err := F.Distraccion(ctx, &pb.DistraccionRequest{Turnos: int32(turnos)})
		if err != nil {
			log.Fatalf("No se pudo saludar: %v", err)
		}
		if f.GetConfirmacion() {
			turnos = 200 - probT
			t, err := T.Golpe(ctx, &pb.GolpeRequest{Turnos: int32(turnos)})
			if err != nil {
				log.Fatalf("No se pudo saludar: %v", err)
			}
			if t.GetConfirmacion() {
				log.Printf("Nice")
			} else {
				victoria = false
			}
		} else {
			victoria = false
		}
	} else {
		turnos := 200 - probT
		log.Printf("%d", turnos)
		t, err := T.Distraccion(ctx, &pb.DistraccionRequest{Turnos: int32(turnos)})
		if err != nil {
			log.Fatalf("No se pudo saludar: %v", err)
		}
		if t.GetConfirmacion() {
			turnos = 200 - probF
			f, err := F.Golpe(ctx, &pb.GolpeRequest{Turnos: int32(turnos)})
			if err != nil {
				log.Fatalf("No se pudo saludar: %v", err)
			}
			if f.GetConfirmacion() {
				botin += int(f.GetBotinExtra())
				log.Printf("Nice")
			} else {
				victoria = false
			}
		} else {
			victoria = false
		}
	}
	if victoria {
		file, err := os.Create("algo.txt")
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer file.Close()

		extra := botin % 4
		botin -= extra
		botin /= 4
		log.Printf("botin lester: %d + extra: %d. Total: %d", botin, extra, botin+extra)
		log.Printf("botin franklin: %d", botin)
		log.Printf("botin trevor: %d", botin)
		log.Printf("nice")
	} else {
		log.Printf("bad")
	}
}
