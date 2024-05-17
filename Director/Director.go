package main

import (
	pb "Lab3SD/Proto"
	"context"
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
	grpcServer *grpc.Server
}

//Implementa la funcion SolicitarM de la interfaz ServicioRecursos de RCP

func (s *server) MensajeDirector(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {
	mutex.Lock() //Bloquea recursos


    if req.GetPeticion() == 1 { 
        print("ostia chaval")
        return &pb.DirectorMensaje{Inicio: 1}, nil //Se retorna en funcion SolicitarM un mensaje de aprobacion

    } else {
        // No hay suficientes recursos

        return &pb.DirectorMensaje{Inicio: 0}, nil  //Se retorna en funcion SolicitarM un mensaje de denegacion
    }
}




func main() {

	grpcServer := grpc.NewServer() //Se crea el servidor

	mutex.Lock()
	s := &server{ //Se le asignan los recursos al servidor
		fase: 0,
		grpcServer: grpcServer,
    }
	mutex.Unlock()

	pb.RegisterMercDirServer(grpcServer, s) //Se registra el servidor
	
	addr := "10.35.169.91:8080"  //Se asigna la direccion del servidor
	lis, err := net.Listen("tcp", addr) //Se crea el listener
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	
	if err := grpcServer.Serve(lis); err != nil {  //Se inicia el servidor
        log.Fatalf("Fallo al crear servidor: %s", err)
    }
	
	fmt.Printf("Director jeje\n")

}
