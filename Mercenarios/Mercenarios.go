package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"

	pb "Lab3SD/Proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
var mutex = &sync.Mutex{} //Se crea un mutex para evitar problemas de concurrencia

func main() {
	var wg sync.WaitGroup

	// Se crean los cuatro grupos
	for i := 0; i < 1; i++ {  // Se crean 4 grupos, si se quiere modificar se debe cambiar el 4 por otro numero y el archivo Central.go linea 53
		wg.Add(1)  
		go InicioMercenario(i + 1, &wg, 1)  //Empieza ejecucion de equipo
		//fmt.Printf("Mercenario %d redi!\n", i+1)
	}

	// Espera que todas las ejecuciones terminen para finalizar la ejecucion del codigo.
	wg.Wait()

	fmt.Println("Todos los equipos han terminado")

}

func InicioMercenario(id int,wg *sync.WaitGroup, Estado int32) { //Toma como parametros el id del equipo y el grupo de espera
  //1 = vivo, 0 = muerto
	var Nivel int32
	var Sent bool = false
	defer wg.Done()
	serverAddr := "0.0.0.0:8080"
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))  //Se conecta al servidor central
	if err != nil {
		fmt.Println("Error al conectar al servidor central:", err)
		return
	}
	defer conn.Close()
	c := pb.NewMercDirClient(conn)

	for {
		if Estado == 0{ //Muere mercenario
			fmt.Printf("Muere mercenario %d\n", id)
			break
		}

		if !Sent {
			stream, err := c.MensajeDirector(context.Background(), &pb.MercenarioMensaje{Peticion: 1, Id: int32(id)}) //Manda solicitud de inicio a director
			//fmt.Printf("Mensaje enviado")
			if err != nil {
				log.Fatalf("Error on stream: %v", err)
			}
			Sent = true
			for {
				msg, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("%v.StreamMessages(_) = _, %v", c, err)
				}
				if msg.GetInicio() == 1 { //Recibe solicitud de inicio del director
					fmt.Println("Empieza mision mercenario:", id)
					Nivel = msg.GetFase()
					Estado = 1
					break
				}
			}
		}
		
		
		if Nivel == 1{ //Si la fase enviada por director es 1.
			fmt.Println("Mercenario ", id, " en nivel 1")
			randomNumber := rand.Intn(3) + 1
			for {
				stream, err := c.Fase1(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)})
				if err != nil {
					log.Fatalf("Error on stream: %v", err)
				}
				msg, _ := stream.Recv()
				if msg.GetEstado() == 1 {
					fmt.Println("Mercenario", id,"pasa a nivel 2")
					Estado = msg.GetEstado()
					Nivel = msg.Fase
					Sent = false
					break
				}
				if msg.GetEstado() == 0 {
					Estado = msg.GetEstado()
					Nivel = 0
					break
				}
			}
		}
			

		if Nivel == 2{
			fmt.Println("Mercenario ", id, " en nivel 2")
			randomNumber := rand.Intn(2) + 1
			for {
				stream, err := c.Fase2(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)})
				if err != nil {
					log.Fatalf("Error al recibir mensaje de Director: %v", err)
				}
				msg, _ := stream.Recv()
				if msg.GetEstado() == 1 {
					fmt.Println("Mercenario", id,"pasa a nivel 3")
					Estado = msg.GetEstado()
					Nivel = msg.Fase
					Sent = false
					break
				}
				if msg.GetEstado() == 0 {
					Estado = msg.GetEstado()
					Nivel = 0
					break
				}
			}
		}

		if Nivel == 3 {
			var aciertos int = 0
			for i := 1; i <= 5; i++ {
				randomNumber := rand.Intn(15) + 1
				fmt.Println("Numero elegido: ", randomNumber, " en nivel 3")
				stream, err := c.Fase3(context.Background(), &pb.MercenarioMensaje{Decision: int32(randomNumber), Id: int32(id)})
				if err != nil {
					log.Fatalf("Error al recibir mensaje de Director: %v", err)
				}
				msg, err := stream.Recv()
				if err != nil {
					log.Fatalf("Error al recibir mensaje del stream: %v", err)
				}
				if msg.GetEstado() == 1 {
					fmt.Printf("Mercenario %d ha acertado\n", id)
					aciertos += 1
				}
				if aciertos >= 2 {
					fmt.Println("Mercenario", id, "ha derrotado al Patriarca!")
					Nivel = 0 // Resetting Nivel to 0
					break
				}
			}
		}
		
	}
	
}
