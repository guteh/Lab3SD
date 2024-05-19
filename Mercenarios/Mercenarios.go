package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	pb "Lab3SD/Proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Synchronize access to shared resources
var mutex = &sync.Mutex{}

func main() { 
	var wg sync.WaitGroup

	numMercenarios := 8 

	// Inicio mercenarios
	for i := 0; i < numMercenarios; i++ {  
		wg.Add(1)
		go InicioMercenario(i+1, &wg, 1)
	}


	wg.Wait()
	time.Sleep(3 * time.Second)
	fmt.Println("Termina la mision de los mercenarios")
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

	_, err = c.SolicitarM(context.Background(), &pb.MercenarioMensaje{Peticion: 1, Id: int32(id)})  //Mando solicitud para inciar mision
	if err != nil {
		log.Fatalf("Error al solicitar mision: %v", err)
	}

	_, err = c.IniciarMision(context.Background(), &pb.MercenarioMensaje{Id: int32(id)})  //Espero para inicar mision
	if err != nil {
		log.Fatalf("Error al iniciar mision: %v", err)
	}

	fmt.Printf("Mercenario %d ha iniciado la misión\n", id)
	Nivel = 1


	if Nivel == 1 {
		fmt.Println("Mercenario ", id, " en piso 1")  //Mercenario en nivel 1
		randomNumber := rand.Intn(3) + 1  //Genero numero aleatorio entre 1 y 3
		for {
			resp, err := c.Fase1(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)})  //Envio decision
			if err != nil {
				log.Fatalf("Error en Fase1: %v", err)
			}
			if resp.GetEstado() == 1 {  //Sobrevive mercenario, pasa al siguiente piso
				fmt.Println("Mercenario", id, "pasa al piso 2")  
				Nivel = 2
				break

			}
			if resp.GetEstado() == 0 {  //Recibe que murio el mercenario
				fmt.Printf("Muere mercenario %d\n", id)  
  				Nivel = 0
				return
			}
		}
	}

	if Nivel == 2 {  //Mercenario en nivel 2
		randomNumber := rand.Intn(2) + 1  //Genero numero aleatorio entre 1 y 2
		for {  
			resp, err := c.Fase2(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)}) //Envio decision
			if err != nil {
				log.Fatalf("Error en Fase2: %v", err)
			}
			if resp.GetEstado() == 1 { //Sobrevive mercenario, pasa al siguiente piso
				fmt.Println("Mercenario", id, "pasa a piso 3")
				Nivel = 3
				break
			}
			if resp.GetEstado() == 0 { //Recibe que murio el mercenario
				fmt.Printf("Muere mercenario %d\n", id)
				Nivel = 0
				return
			}
		}
	}

	if Nivel == 3{ //Mercenario en nivel 3
		numeros := make([]int32, 5) //Genero 5 numeros aleatorios entre 1 y 15
		for i := 0; i < 5; i++ {
			numeros[i] = rand.Int31n(15) + 1
			fmt.Println("Numero elegido: ", numeros[i], " en nivel 3")
		}
		
		resp, err := c.Fase3(context.Background(), &pb.MercenarioMensaje{Decisiones: numeros, Id: int32(id)}) //Envio lista de numeros
		if err != nil {
			log.Fatalf("Error en Fase1: %v", err)
		}
		if resp.GetEstado() == 1 {
			fmt.Println("Mercenario", id, "ha derrotado al Patriarca!") //Mercenario derrota al patriarca
			Nivel = 0
			return
		}
		if resp.GetEstado() == 0 {
			fmt.Printf("Muere mercenario %d\n", id) //Recibe que murio el mercenario
			Nivel = 0
			return
		}
	}
	
}
