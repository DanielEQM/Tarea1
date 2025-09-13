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
	addressL = "10.35.168.35:50051"
	addressF = "10.35.168.37:50052"
	addressT = "10.35.168.38:50053" // La dirección del servidor
)

func fallo(err error, msgs string) {
	if err != nil {
		log.Printf("%s: %v", msgs, err)
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

func main() {
	connL, err := grpc.Dial(addressL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	fallo(err, "No se pudo conectar")
	defer connL.Close()

	connF, err := grpc.Dial(addressF, grpc.WithTransportCredentials(insecure.NewCredentials()))
	fallo(err, "No se pudo conectar")
	defer connF.Close()

	connT, err := grpc.Dial(addressT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	fallo(err, "No se pudo conectar")
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
		fallo(err, "No se pudo conectar")
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
			fallo(err, "No se pudo conectar")
			break
		} else {
			log.Println("Lester no tiene ofertas en este momento, reintentando...")
			time.Sleep(2 * time.Second)
		}
	}

	log.Println("=== Fase 1 completada con éxito ===")

	_, err = L.ConfirmMission(ctx, &pb.ConfirmRequest{Conf: true})
	fallo(err, "No se pudo conectar")

	botin, _ := strconv.Atoi(r.GetBotin())
	probF, _ := strconv.Atoi(r.GetProbF())
	probT, _ := strconv.Atoi(r.GetProbT())
	botinExtra := 0
	riesgo, _ := strconv.Atoi(r.GetRiesgo())
	log.Printf("probF: %d y probT: %d", probF, probT)

	victoria := true
	if probF > probT {
		log.Println("¡Franklin inicia la distracción!")
		turnos := 200 - probF
		log.Printf("%d", turnos)
		f, err := F.Distraccion(ctx, &pb.DistraccionRequest{Turnos: int32(turnos)})
		fallo(err, "No se pudo conectar")

		if f.GetConfirmacion() {
			log.Println("=== Fase 2 completada con éxito ===")
			log.Println("Franklin logró distraer, Trevor ejecuta el golpe!")
			turnos = 200 - probT
			_, err = L.NotificarEstrellas(ctx, &pb.AvisoRequest{Pj: "Trevor", Riesgo: int32(riesgo), Turnos: int32(turnos)})
			fallo(err, "No se pudo conectar")

			t, err := T.Golpe(ctx, &pb.GolpeRequest{Turnos: int32(turnos)})
			fallo(err, "No se pudo conectar")

			if t.GetConfirmacion() {
				_, err = L.NotificarGolpe(ctx, &pb.ConfirmRequest{Conf: false})
				fallo(err, "No se pudo conectar")

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
		t, err := T.Distraccion(ctx, &pb.DistraccionRequest{Turnos: int32(turnos)})
		fallo(err, "No se pudo conectar")

		if t.GetConfirmacion() {
			log.Println("=== Fase 2 completada con éxito ===")
			log.Println("Trevor logró distraer, Franklin ejecuta el golpe!")
			turnos = 200 - probF
			_, err = L.NotificarEstrellas(ctx, &pb.AvisoRequest{Pj: "Franklin", Riesgo: int32(riesgo), Turnos: int32(turnos)})
			fallo(err, "No se pudo conectar")

			f, err := F.Golpe(ctx, &pb.GolpeRequest{Turnos: int32(turnos)})
			fallo(err, "No se pudo conectar")

			if f.GetConfirmacion() {
				botinExtra = int(f.GetBotinExtra())
				log.Printf("=== Atraco completado con éxito ===")
			} else {
				_, err = L.NotificarGolpe(ctx, &pb.ConfirmRequest{Conf: false})
				fallo(err, "No se pudo conectar")

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
	file, err := os.Create("informe.txt")
	fallo(err, "No se pudo conectar")

	defer file.Close()

	if victoria {
		line := "==============================================\n"
		line += "==        REPORTE FINAL DE LA MISIÓN        ==\n"
		line += "==============================================\n"
		line += "Misión: Asalto al banco\n"
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
		err = os.WriteFile("informe.txt", []byte(line), 0606)
		fallo(err, "No se pudo conectar")

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
		err = os.WriteFile("informe.txt", []byte(line), 0606)
		fallo(err, "No se pudo conectar")

	}
	log.Printf("=== FASE 4 COMPLETADA ===")
}
