package main

import (
	pb "Lab3SD/Proto"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)
var mutex = &sync.Mutex{} //Se crea un mutex para evitar problemas de concurrencia



type server struct {  //Crea el servidor rcp
    pb.UnimplementedMercDirServer
	fase int
	suma int
	nmercenarios int
	grpcServer *grpc.Server
}

//Implementa la funcion SolicitarM de la interfaz ServicioRecursos de RCP

func (s *server) MensajeDirector(req *pb.MercenarioMensaje, stream pb.MercDir_MensajeDirectorServer) error {
    if req.GetPeticion() == 1 { 
		mutex.Lock() //Bloquea recursos
		s.suma += 1 //Se suma 1 a la variable suma
		mutex.Unlock() // Desbloquea recursos
		fmt.Printf("Se ha sumado 1 a la suma\n")
	}
	if s.suma == s.nmercenarios {
        if err := stream.Send(&pb.DirectorMensaje{Inicio: 1}); err != nil {
            return err
        }
        
	}

	return nil
}




func main() {

	grpcServer := grpc.NewServer() //Se crea el servidor

	mutex.Lock()
	s := &server{ //Se le asignan los recursos al servidor
		fase: 0,
		grpcServer: grpcServer,
		suma: 0,
		nmercenarios: 10,
    }
	mutex.Unlock()

	pb.RegisterMercDirServer(grpcServer, s) //Se registra el servidor
	
	addr := "0.0.0.0:8080"  //Se asigna la direccion del servidor
	lis, err := net.Listen("tcp", addr) //Se crea el listener
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	
	if err := grpcServer.Serve(lis); err != nil {  //Se inicia el servidor
        log.Fatalf("Fallo al crear servidor: %s", err)
    }
	
	fmt.Printf("Director jeje\n")

}
