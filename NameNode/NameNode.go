package main

import (
	pb "Lab3SD/Proto"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {  //Crea el servidor rcp con sus variables globales
    pb.UnimplementedDirNameServer
	grpcServer *grpc.Server
	decisions1 map[int][]int32
	decisions2 map[int][]int32
	decisions3 map[int][]int32
	mercenarios map[string]int
	txt bool
}

func (s *server) RegistrosDirector(ctx context.Context, req *pb.EnviarDecision) (*emptypb.Empty, error) {
	if !s.txt {
		file, err := os.Create("direcciones.txt")
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()
		s.txt = true
	}
    log.Printf("Recibida decision de %s en el piso %d: %d", req.GetNombre(), req.GetPiso(), req.GetDecision())
	piso := strconv.Itoa(int(req.GetPiso()))
	if _, exists := s.mercenarios[req.GetNombre()+"_"+piso]; !exists {
		s.mercenarios[req.GetNombre()] = rand.Intn(3)
	}
	
	if s.mercenarios[req.GetNombre()] == 0 {
		file, err := os.OpenFile("direcciones.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		line := fmt.Sprintf("%s %d 10.35.169.91\n", req.GetNombre(), req.GetPiso())
		if _, err := file.WriteString(line); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	
	}
	if s.mercenarios[req.GetNombre()] == 1 {
		file, err := os.OpenFile("direcciones.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		line := fmt.Sprintf("%s %d 10.35.169.93\n", req.GetNombre(), req.GetPiso())
		if _, err := file.WriteString(line); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	}

	if s.mercenarios[req.GetNombre()] == 2 {
		file, err := os.OpenFile("direcciones.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		line := fmt.Sprintf("%s %d 10.35.169.94\n", req.GetNombre(), req.GetPiso())
		if _, err := file.WriteString(line); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	
	}
    return &emptypb.Empty{}, nil
}


func main() {
	grpcServer := grpc.NewServer() //Se crea el servidor


	s := &server{ //Se le asignan los recursos al servidor
		grpcServer:  grpcServer,
		decisions1:  make(map[int][]int32),
		decisions2:  make(map[int][]int32),
		decisions3:  make(map[int][]int32),
		mercenarios: make(map[string]int),
		txt : false,
	}

	pb.RegisterDirNameServer(grpcServer, s) //Se registra el servidor

	addr := "0.0.0.0:8081"  //Se asigna la direccion del servidor
	lis, err := net.Listen("tcp", addr) //Se crea el listener
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	log.Println("NameNode escuchando solicitudes", addr)
	if err := grpcServer.Serve(lis); err != nil {  //Se inicia el servidor
        log.Fatalf("Fallo al crear servidor: %s", err)
    }
}
