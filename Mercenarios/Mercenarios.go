package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	pb "Lab3SD/Proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Synchronize access to shared resources
var mutex = &sync.Mutex{}

func main() {
	var wg sync.WaitGroup

	numMercenarios := 8 // Number of mercenaries to run concurrently

	// Create and start mercenaries
	for i := 0; i < numMercenarios; i++ {
		wg.Add(1)
		go InicioMercenario(i+1, &wg, 1)
	}

	// Wait for all mercenaries to complete
	wg.Wait()

	fmt.Println("Todos los equipos han terminado")
}

func InicioMercenario(id int, wg *sync.WaitGroup, Estado int32) {
	defer wg.Done()

	var Nivel int32
	serverAddr := "0.0.0.0:8080"

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error al conectar al servidor central:", err)
		return
	}
	defer conn.Close()

	c := pb.NewMercDirClient(conn)

	_, err = c.SolicitarM(context.Background(), &pb.MercenarioMensaje{Peticion: 1, Id: int32(id)})
	if err != nil {
		log.Fatalf("Error al solicitar mision: %v", err)
	}

	_, err = c.IniciarMision(context.Background(), &pb.MercenarioMensaje{Id: int32(id)})
	if err != nil {
		log.Fatalf("Error al iniciar mision: %v", err)
	}

	fmt.Printf("Mercenario %d ha iniciado la misión\n", id)
	Nivel = 1
	for {
		if Estado == 0 {
			fmt.Printf("Muere mercenario %d\n", id)
			break
		}


		switch Nivel {
		case 1:
			fmt.Println("Mercenario ", id, " en nivel 1")
			randomNumber := rand.Intn(3) + 1
			for {
				resp, err := c.Fase1(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)})
				if err != nil {
					log.Fatalf("Error en Fase1: %v", err)
				}
				if resp.GetEstado() == 1 {
					fmt.Println("Mercenario", id, "pasa a nivel 2")
					Estado = resp.GetEstado()
					Nivel = 2
					break
				}
				if resp.GetEstado() == 0 {
					Estado = resp.GetEstado()
					Nivel = 0
					break
				}
			}

		case 2:
			fmt.Println("Mercenario ", id, " en nivel 2")
			randomNumber := rand.Intn(2) + 1
			for {
				resp, err := c.Fase2(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)})
				if err != nil {
					log.Fatalf("Error en Fase2: %v", err)
				}
				if resp.GetEstado() == 1 {
					fmt.Println("Mercenario", id, "pasa a nivel 3")
					Estado = resp.GetEstado()
					Nivel = 3
					break
				}
				if resp.GetEstado() == 0 {
					Estado = resp.GetEstado()
					Nivel = 0
					break
				}
			}

		case 8:
			var aciertos int = 0
			for i := 1; i <= 5; i++ {
				randomNumber := rand.Intn(15) + 1
				fmt.Println("Numero elegido: ", randomNumber, " en nivel 3")
				resp, err := c.Fase3(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)})
				if err != nil {
					log.Fatalf("Error en Fase1: %v", err)
				}
				if resp.GetEstado() == 1 {
					fmt.Printf("Mercenario %d ha acertado\n", id)
					aciertos += 1
				}
				if aciertos >= 2 {
					fmt.Println("Mercenario", id, "ha derrotado al Patriarca!")
					Nivel = 0
					break
				}
			}
		}
	}
}
