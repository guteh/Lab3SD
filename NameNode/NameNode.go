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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	dNode1 pb.NameDataClient
	dNode2 pb.NameDataClient
	dNode3 pb.NameDataClient
}

func (s *server) RegistrosDirector(ctx context.Context, req *pb.EnviarDecision) (*emptypb.Empty, error) {
	if !s.txt {  //Si direcciones.txt no existe, crearlo, y si existe vaciarlo
		file, err := os.Create("direcciones.txt")
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()
		s.txt = true
	}
    log.Printf("Recibida decision de %s en el piso %d: %d", req.GetNombre(), req.GetPiso(), req.GetDecision())
	piso := strconv.Itoa(int(req.GetPiso()))
	if _, exists := s.mercenarios[req.GetNombre()+"_"+piso]; !exists {  //Crear registo de ID de mercenario y piso
		s.mercenarios[req.GetNombre()] = rand.Intn(3)
	}
	
	if s.mercenarios[req.GetNombre()] == 0 {
		file, err := os.OpenFile("direcciones.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		line := fmt.Sprintf("%s Piso_%d 10.35.169.91\n", req.GetNombre(), req.GetPiso())
		if _, err := file.WriteString(line); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		
		if req.GetPiso() < 3 {
			s.dNode1.RegistroMercenario(context.Background(), &pb.EnviarDecision{Nombre: req.GetNombre(), Piso: req.GetPiso(), Decision: req.GetDecision()})
		}
		if req.GetPiso() == 3 {
			s.dNode1.RegistroMercenario(context.Background(), &pb.EnviarDecision{Nombre: req.GetNombre(), Piso: req.GetPiso(), Decisiones: req.GetDecisiones()})
		}
		
	
	}
	if s.mercenarios[req.GetNombre()] == 1 {
		file, err := os.OpenFile("direcciones.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		line := fmt.Sprintf("%s Piso_%d 10.35.169.92\n", req.GetNombre(), req.GetPiso())
		if _, err := file.WriteString(line); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		/* FALTA IMPLEMENTACION DOSH BANK
		if req.GetPiso() < 3 {
			s.dNode2.RegistroMercenario(context.Background(), &pb.EnviarDecision{Nombre: req.GetNombre(), Piso: req.GetPiso(), Decision: req.GetDecision()})
		}
		if req.GetPiso() == 3 {
			s.dNode2.RegistroMercenario(context.Background(), &pb.EnviarDecision{Nombre: req.GetNombre(), Piso: req.GetPiso(), Decisiones: req.GetDecisiones()})
		}
		*/
	}

	if s.mercenarios[req.GetNombre()] == 2 {
		file, err := os.OpenFile("direcciones.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		line := fmt.Sprintf("%s Piso_%d 10.35.169.94\n", req.GetNombre(), req.GetPiso())
		if _, err := file.WriteString(line); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		if req.GetPiso() < 3 {
			s.dNode3.RegistroMercenario(context.Background(), &pb.EnviarDecision{Nombre: req.GetNombre(), Piso: req.GetPiso(), Decision: req.GetDecision()})
		}
		if req.GetPiso() == 3 {
			s.dNode3.RegistroMercenario(context.Background(), &pb.EnviarDecision{Nombre: req.GetNombre(), Piso: req.GetPiso(), Decisiones: req.GetDecisiones()})
		}
	
	}
    return &emptypb.Empty{}, nil
}

func StartServer(s *server, grpcServer *grpc.Server){
	pb.RegisterDirNameServer(grpcServer, s) //Se registra el servidor
	addr := "172.17.0.1:8080"  //Se asigna la direccion del servidor
	lis, err := net.Listen("tcp", addr) //Se crea el listener
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	log.Println("NameNode escuchando solicitudes", addr)
	if err := grpcServer.Serve(lis); err != nil {  //Se inicia el servidor
        log.Fatalf("Fallo al crear servidor: %s", err)
    }
}


func main() {
	conn1, err := grpc.NewClient("172.17.0.1:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))  //DATA NODE1 10.35.169.91:8084
    if err != nil {
        log.Fatalf("Fallo al conectarse a NameNode: %v", err)
    }
    defer conn1.Close()
    DataNode1 := pb.NewNameDataClient(conn1)

	
	conn2, err := grpc.NewClient("0.0.0.0:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))  //DATANODE 2
    if err != nil {
        log.Fatalf("Fallo al conectarse a NameNode: %v", err)
    }
    defer conn2.Close()
    DataNode2 := pb.NewNameDataClient(conn2)

	conn3, err := grpc.NewClient("10.35.169.94:8086", grpc.WithTransportCredentials(insecure.NewCredentials()))  //DATANODE 3 10.35.169.94:8086
    if err != nil {
        log.Fatalf("Fallo al conectarse a NameNode: %v", err)
    }
    defer conn3.Close()
    DataNode3 := pb.NewNameDataClient(conn3)
	
	grpcServer := grpc.NewServer() //Se crea el servidor

	s := &server{ //Se le asignan los recursos al servidor
		grpcServer:  grpcServer,
		decisions1:  make(map[int][]int32),
		decisions2:  make(map[int][]int32),
		decisions3:  make(map[int][]int32),
		mercenarios: make(map[string]int),
		txt : false,
		dNode1: DataNode1,
		dNode2: DataNode2,
		dNode3: DataNode3,
	}

	go StartServer(s, grpcServer) //Se inicia el servidor

	time.Sleep(20 * time.Second)

}
