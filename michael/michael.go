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
	addressL = "localhost:50051"
	addressF = "localhost:50052"
	addressT = "localhost:50053" // La dirección del servidor
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
	botinExtra := 0
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
				botinExtra += int(f.GetBotinExtra())
				log.Printf("Nice")
			} else {
				victoria = false
			}
		} else {
			victoria = false
		}
	}
	//---
	botinTotal := botin + botinExtra
	extraLester := (botinTotal) % 4
	botinTotal -= extraLester
	botinTotal /= 4
	//---
	file, err := os.Create("Informe.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	if victoria {
		line := "==============================================\n"
		line += "==        REPORTE FINAL DE LA MISIÓN        ==\n"
		line += "==============================================\n"
		line += "Misión: Asalato al banco\n"
		line += "Resultado Global: MISSION COMPLETADA CON EXITO\n\n"
		line += "       <--    REPARTO DEL BOTIN    -->        \n"
		line += "Botin Base: $" + strconv.Itoa(botin) + "\n"
		line += "Botin Extra (Habilidad de Chop): $" + strconv.Itoa(botinExtra) + "\n"
		line += "Botin Total: $" + strconv.Itoa(botinTotal) + "\n"
		line += "----------------------------------------------\n"
		line += "botin franklin: $" + strconv.Itoa(botin) + "\n"
		line += "\n"
		line += "botin trevor: $" + strconv.Itoa(botin) + "\n"
		line += "\n"
		line += "botin lester: $" + strconv.Itoa(botin) + " + extra: $" + strconv.Itoa(extraLester) + ". Total: $" + strconv.Itoa(botin+extraLester) + "\n"
		line += "\n"
		line += "----------------------------------------------\n"
		line += "Saldo Final: $" + strconv.Itoa(botinTotal) + "\n"
		line += "=============================================="
		err = os.WriteFile("Informe.txt", []byte(line), 0606)
		if err != nil {
			log.Printf("%v", err)
		}

		log.Printf("nice")
	} else {
		line := "============================================\n"
		line += "==       REPORTE FINAL DE LA MISIÓN       ==\n"
		line += "============================================\n"
		line += "Misión: Asalato al banco\n"
		line += "Resultado Global: MISSION INCOMPLETA...\n\n"
		line += "      <--    REPARTO DEL BOTIN    -->       \n"
		line += "Botin Base: $" + strconv.Itoa(botin) + "\n"
		line += "Botin Extra (Habilidad de Chop): $" + strconv.Itoa(botinExtra) + "\n"
		line += "Botin Total Perdido: $" + strconv.Itoa(botinTotal) + "\n"
		line += "--------------------------------------------\n"
		line += "botin franklin: $0\n"
		line += "\n"
		line += "botin franklin: $0\n"
		line += "\n"
		line += "botin lester: $0 + extra: $0. Total: $0\n"
		line += "\n"
		line += "--------------------------------------------\n"
		line += "Saldo Final Perdido: $0\n"
		line += "============================================"
		err = os.WriteFile("Informe.txt", []byte(line), 0606)
		if err != nil {
			log.Printf("%v", err)
		}
		log.Printf("bad")
	}
}
