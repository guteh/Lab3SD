package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	pb "Lab3SD/Proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
var mutex = &sync.Mutex{} //Se crea un mutex para evitar problemas de concurrencia

func main() {
	var wg sync.WaitGroup



	// Se crean los cuatro grupos
	for i := 0; i < 10; i++ {  // Se crean 4 grupos, si se quiere modificar se debe cambiar el 4 por otro numero y el archivo Central.go linea 53
		wg.Add(1)   
		go InicioMercenario(i + 1, &wg)  //Empieza ejecucion de equipo
		fmt.Printf("Mercenario %d redi!\n", i+1)
	}

	// Espera que todas las ejecuciones terminen para finalizar la ejecucion del codigo.
	wg.Wait()

	fmt.Println("Todos los equipos han terminado")

}

func InicioMercenario(id int,wg *sync.WaitGroup) { //Toma como parametros el id del equipo y el grupo de espera

	defer wg.Done()

	serverAddr := "0.0.0.0:8080"
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))  //Se conecta al servidor central
	if err != nil {
		fmt.Println("Error al conectar al servidor central:", err)
		return
	}
	defer conn.Close()

	c := pb.NewMercDirClient(conn)
	stream, err := c.MensajeDirector(context.Background(), &pb.MercenarioMensaje{Peticion: 1})
	if err != nil {
		log.Fatalf("Error on stream: %v", err)
	}
	
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.StreamMessages(_) = _, %v", c, err)
		}
		if msg.GetInicio() == 1 {
			mutex.Lock()
			fmt.Println("Empieza mision mercenario:", id)
			mutex.Unlock()
		}
	}
}
