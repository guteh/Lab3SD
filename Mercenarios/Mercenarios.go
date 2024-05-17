package main

import (
	"context"
	"fmt"
	"sync"

	pb "Lab3SD/Proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var wg sync.WaitGroup



	// Se crean los cuatro grupos
	for i := 0; i < 1; i++ {  // Se crean 4 grupos, si se quiere modificar se debe cambiar el 4 por otro numero y el archivo Central.go linea 53
		wg.Add(1)
		go InicioMercenario(i + 1, &wg)  //Empieza ejecucion de equipo
		fmt.Printf("Mercenario %d redi!\n", i+1)
	}

	// Espera que todas las ejecuciones terminen para finalizar la ejecucion del codigo.
	wg.Wait()

}

func InicioMercenario(id int,wg *sync.WaitGroup) { //Toma como parametros el id del equipo y el grupo de espera

	defer wg.Done()

	// Genera cantidades random de recursos

	serverAddr := "0.0.0.0:8080"
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))  //Se conecta al servidor central
	if err != nil {
		fmt.Println("Error al conectar al servidor central:", err)
		return
	}
	defer conn.Close()

	c := pb.NewMercDirClient(conn)

	for {
		response, err := c.MensajeDirector(context.Background(), &pb.MercenarioMensaje{Peticion: 1, Decision: 1})  //Envia peticion de recursos a servidor central
		if err != nil {
			fmt.Println("Error al enviar el mensaje al servidor central:", err)
		}
		if response.Inicio == 1 {
			fmt.Printf("Se inicia mision")
			break  //Si la respuesta es 1, se cierra la comunicacion
		} else {
			fmt.Printf("No hay mision")
		}
	}
}
