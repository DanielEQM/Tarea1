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
	addressL = "lester:50051"
	addressF = "franklin:50052"
	addressT = "trevor:50053" // La dirección del servidor
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
			log.Printf("Michael rechazó 3 veces. Lester lo hace esperar 10s ...")
			time.Sleep(10 * time.Second)
			rechazo = 0
		}
		if r.GetDisp() {
			log.Printf("Oferta recibida -> Botín: %s, ProbF: %s, ProbT: %s, Riesgo: %s", r.GetBotin(), r.GetProbF(), r.GetProbT(), r.GetRiesgo())
			if r.GetBotin() == "" || r.GetProbF() == "" || r.GetProbT() == "" || r.GetRiesgo() == "" {
				rechazo++
				continue
			}
			probF, _ := strconv.Atoi(r.GetProbF())
			probT, _ := strconv.Atoi(r.GetProbT())
			riesgo, _ := strconv.Atoi(r.GetRiesgo())
			if (probF <= 50 && probT <= 50) || riesgo >= 80 {
				log.Printf("Michael rechaza la oferta...")
				rechazo++
				continue
			}

			log.Println("¡Michael acepta la oferta!")
			_, err = L.ConfirmMission(ctx, &pb.ConfirmRequest{Conf: true})
			if err != nil {
				log.Fatalf("Error al confirmar oferta: %v", err)
			}
			break
		} else {
			log.Println("Lester no tiene ofertas en este momento, reintentando...")
			time.Sleep(2 * time.Second)
		}
	}

	log.Println("=== Fase 1 completada con éxito ===")

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
		log.Println("¡Franklin inicia la distracción!")
		turnos := 200 - probF
		log.Printf("%d", turnos)
		f, err := F.Distraccion(ctx, &pb.DistraccionRequest{Turnos: int32(turnos)})
		if err != nil {
			log.Fatalf("No se pudo saludar: %v", err)
		}
		if f.GetConfirmacion() {
			log.Println("=== Fase 2 completada con éxito ===")
			log.Println("Franklin logró distraer, Trevor ejecuta el golpe!")
			turnos = 200 - probT
			t, err := T.Golpe(ctx, &pb.GolpeRequest{Turnos: int32(turnos)})
			if err != nil {
				log.Fatalf("Error en el golpe de Trevor: %v", err)
			}
			if t.GetConfirmacion() {
				log.Printf("=== Atraco completado con éxito ===")
			} else {
				log.Println("Trevor falló en el golpe")
				victoria = false
			}
		} else {
			log.Println("Franklin falló en la distracción")
			victoria = false
		}
	} else {
		log.Println("=== Trevor inicia la distracción ===")
		turnos := 200 - probT
		log.Printf("%d", turnos)
		t, err := T.Distraccion(ctx, &pb.DistraccionRequest{Turnos: int32(turnos)})
		if err != nil {
			log.Fatalf("No se pudo saludar: %v", err)
		}
		if t.GetConfirmacion() {
			log.Println("=== Fase 2 completada con éxito ===")
			log.Println("Trevor logró distraer, Franklin ejecuta el golpe!")
			turnos = 200 - probF
			f, err := F.Golpe(ctx, &pb.GolpeRequest{Turnos: int32(turnos)})
			if err != nil {
				log.Fatalf("No se pudo saludar: %v", err)
			}
			if f.GetConfirmacion() {
				botinExtra += int(f.GetBotinExtra())
				log.Printf("=== Atraco completado con éxito ===")
			} else {
				log.Println("Franklin falló en el golpe")
				victoria = false
			}
		} else {
			log.Println("Trevor falló en la distracción")
			victoria = false
		}
	}

	if victoria {
		log.Println("=== Fase 3 completada con éxito ===")
	} else {
		log.Println("=== Fase 3 fracasó ===")
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
		line += "Misión: Asalto al banco\n"
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
