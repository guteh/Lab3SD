package main

import (
	pb "Lab3SD/Proto"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {  //Crea el servidor rcp con sus variables globales
    pb.UnimplementedNameDataServer
	pb.UnimplementedDirNameServer
	grpcServer *grpc.Server
	grpcServer1 *grpc.Server
	mercenarios map[string]int
	dinero int
	txt bool

}

func (s *server) RegistrosDirector(ctx context.Context, req *pb.EnviarDecision) (*emptypb.Empty, error) {
	s.dinero += 100000000
	if !s.txt {  //Si direcciones.txt no existe, crearlo, y si existe vaciarlo
		file, err := os.Create("montos.txt")
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()
		s.txt = true
	}
    log.Printf("Mercenario %s muerto en el piso %d: Nuevo monto: %d", req.GetNombre(), req.GetPiso(), s.dinero)
	
	file, err := os.OpenFile("montos.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Mercenario %s Piso_%d %d\n", req.GetNombre(), req.GetPiso(), s.dinero)
	if _, err := file.WriteString(line); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	
    return &emptypb.Empty{}, nil
}

func (s *server) RegistroMercenario(ctx context.Context, req *pb.EnviarDecision) (*emptypb.Empty, error) {

	piso := strconv.Itoa(int(req.GetPiso()))
	nombretxt := "DataNode/Mercenario"+req.GetNombre()+"_"+piso+".txt"

	file, err := os.Create(nombretxt)
	if err != nil {
		log.Fatalf("Fallo al crear archivo: %v", err)
	}
	defer file.Close()

	file, err = os.OpenFile(nombretxt, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Fallo al abrir archivo: %v", err)
		}
		defer file.Close()
		decision := strconv.Itoa(int(req.GetDecision()))
		if req.GetPiso() < 3 {
			line := fmt.Sprintf("* "+decision+"\n")
			if _, err := file.WriteString(line); err != nil {
				log.Fatalf("Fallo al escribir en el archivo: %v", err)
			}
		}
		if req.GetPiso() == 3 {
			for i := 0; i < 5; i++ {
				decision := strconv.Itoa(int(req.GetDecisiones()[i]))
				line := fmt.Sprintf("* "+decision+"\n")
				if _, err := file.WriteString(line); err != nil {
					log.Fatalf("Fallo al escribir en el archivo: %v", err)
				}
			}
		}
	return &emptypb.Empty{}, nil
}

func StartServerDosh(s *server, grpcServer *grpc.Server){
	pb.RegisterDirNameServer(grpcServer, s) //Se registra el servidor
	addr := "10.35.169.92:8089"  //Se asigna la direccion del servidor 10.35.169.92:8089
	lis, err := net.Listen("tcp", addr) //Se crea el listener
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	log.Println("DoshBank escuchando solicitudes", addr)
	if err := grpcServer.Serve(lis); err != nil {  //Se inicia el servidor
        log.Fatalf("Fallo al crear servidor: %s", err)
    }
}

func StartServerData(s *server, grpcServer *grpc.Server){
	pb.RegisterNameDataServer(grpcServer, s) //Se registra el servidor
	addr := "10.35.169.94:8085"  //DataNode3 //10.35.169.94:8085
	lis, err := net.Listen("tcp", addr) //Se crea el listener
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	if err := grpcServer.Serve(lis); err != nil {  //Se inicia el servidor
        log.Fatalf("Fallo al crear servidor: %s", err)
    }
}


func main() {
	
	grpcServer := grpc.NewServer() //Se crea el servidor
	grpcServer1 := grpc.NewServer() //Se crea el servidor

	s := &server{ //Se le asignan los recursos al servidor
		grpcServer:  grpcServer,
		grpcServer1: grpcServer1,
		mercenarios: make(map[string]int),
		txt : false,
	}

	go StartServerDosh(s, grpcServer) //Se inicia el servidor
	go StartServerData(s, grpcServer1) //Se inicia el servidor

	time.Sleep(20 * time.Second)

}
