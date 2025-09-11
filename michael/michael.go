package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "michael/proto/michael-sys/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
<<<<<<< HEAD
	addressL = "localhost:50051"
	addressF = "localhost:50052"
	addressT = "localhost:50053" // La dirección del servidor
=======
	addressL = "lester-container:50051"
	addressF = "franklin-container:50052"
	addressT = "trevor-container:50053"
>>>>>>> ce0c21c9088dde9f9826d76919f4243391d600fa
)

func main() {
	connL, err := grpc.Dial(addressL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer connL.Close()

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

	L := pb.NewMissionClient(connL)
	F := pb.NewMissionClient(connF)
	T := pb.NewMissionClient(connT)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

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
			/*log.Printf("y el rechazo %d", rechazo)
			log.Printf("Respuesta del servidor: %t", r.GetDisp())
			log.Printf("          del servidor: %s", r.GetBotin())
			log.Printf("          del servidor: %s", r.GetProbF())
			log.Printf("          del servidor: %s", r.GetProbT())
			log.Printf("          del servidor: %s", r.GetRiesgo())
			log.Printf("")*/
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
	botin, _ := strconv.Atoi(r.GetBotin())
	probF, _ := strconv.Atoi(r.GetProbF())
	probT, _ := strconv.Atoi(r.GetProbT())
	//riesgo, _ := strconv.Atoi(r.GetRiesgo())
	log.Printf("probF: %d y probT: %d", probF, probT)

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
