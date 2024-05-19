package main

import (
	pb "Lab3SD/Proto"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
)
var mutex = &sync.Mutex{} //Se crea un mutex para evitar problemas de concurrencia



type server struct {  //Crea el servidor rcp
    pb.UnimplementedMercDirServer
	fase int
	suma int
	nmercenarios int
	grpcServer *grpc.Server
	startSignalCh  chan struct{}
	Señal1 chan struct{}
	Señal2 chan struct{}
	Señal3 chan struct{}
	decisions1 []int32
	decisions2 []int32
	decisions3 map[int][]int32
	decisor []int32
	mercenarios []string
	X int
	Y int
	camino int
}

//Implementa la funcion SolicitarM de la interfaz ServicioRecursos de RCP


func (s *server) SolicitarM(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Mercenario %d ha llegado!\n", req.GetId())
	s.mercenarios = append(s.mercenarios, strconv.Itoa(int(req.GetId())))
	s.suma += 1

	if s.suma == s.nmercenarios {
		s.fase += 1
		fmt.Printf("Se empieza mision! Fase actual: %d\n", s.fase)
		close(s.startSignalCh) // Close the channel to signal all mercenaries to start
	}

	return &pb.DirectorMensaje{Inicio: 1, Fase: int32(s.fase)}, nil
}

func (s *server) IniciarMision(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {
	<-s.startSignalCh // Wait until the start signal is received
	return &pb.DirectorMensaje{Inicio: 1, Fase: int32(s.fase)}, nil
}




func (s *server) Fase1(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {  //Implementa funcion fase1, le envia a cada mercenario si vive o muere
	mutex.Lock()
	s.decisions1 = append(s.decisions1, req.GetDecision())
	fmt.Printf("Mercenario %d ha enviado su decision: %d\n", req.GetId(), req.GetDecision())
	if len(s.decisions1) == s.nmercenarios {
		for s.X == s.Y{
			s.X = rand.Intn(100)
			s.Y = rand.Intn(100)
		}
		close(s.Señal1)
	}
	mutex.Unlock()
	
	<-s.Señal1 // Wait until all decisions are received
	mutex.Lock()
	defer mutex.Unlock()

	//fmt.Printf("X: %d Y: %d\n", s.X, s.Y)
	

	decisor := rand.Intn(100)

	if s.X < s.Y {
		prob1 := s.X
		prob2 := s.Y - s.X
		prob3 := 100 - s.Y

		if req.GetDecision() == 1{
			if decisor <= prob1{
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}

			if decisor > prob1{
				 
				fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
				s.nmercenarios -= 1
				nombre := strconv.Itoa(int(req.GetId())) // replace with the name you want to remove
				for i := range s.mercenarios {
					if s.mercenarios[i] == nombre { // replace .Name with the actual field name
						s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
						break
					}
				}
				if s.nmercenarios == 1 {
					fmt.Printf("Misión terminada\n")
					go func ()  {
						time.Sleep(1 * time.Second)
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				
			}
		}

		if req.GetDecision() == 2{
			if decisor <= prob2{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil

			}

			if decisor > prob2{
				 
				fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
				s.nmercenarios -= 1
				nombre := strconv.Itoa(int(req.GetId())) // replace with the name you want to remove
				for i := range s.mercenarios {
					if s.mercenarios[i] == nombre { // replace .Name with the actual field name
						s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
						break
					}
				}
				if s.nmercenarios == 1 {
					fmt.Printf("Misión terminada\n")
					go func ()  {
						time.Sleep(1 * time.Second)
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}

		if req.GetDecision() == 3{
			if decisor <= prob3{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}

			if decisor > prob3{
				 
				fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
				s.nmercenarios -= 1
				nombre := strconv.Itoa(int(req.GetId())) // replace with the name you want to remove
				for i := range s.mercenarios {
					if s.mercenarios[i] == nombre { // replace .Name with the actual field name
						s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
						break
					}
				}
				if s.nmercenarios == 1 {
					fmt.Printf("Misión terminada\n")
					go func ()  {
						 time.Sleep(1 * time.Second)
						s.grpcServer.Stop()
						}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}
	}
	if s.Y <= s.X {
		prob1 := s.Y
		prob2 := s.X - s.Y
		prob3 := 100 - s.X
		if req.GetDecision() == 1{
			if decisor <= prob1{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}
			if decisor > prob1{
				 
				fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
				s.nmercenarios -= 1
				nombre := strconv.Itoa(int(req.GetId())) // replace with the name you want to remove
				for i := range s.mercenarios {
					if s.mercenarios[i] == nombre { // replace .Name with the actual field name
						s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
						break
					}
				}
				if s.nmercenarios == 1 {
					fmt.Printf("Misión terminada\n")
					go func ()  {
						 time.Sleep(1 * time.Second)
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}
		
		if req.GetDecision() == 2{
			if decisor <= prob2{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}

			if decisor > prob2{
				 
				fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
				s.nmercenarios -= 1
				nombre := strconv.Itoa(int(req.GetId())) // replace with the name you want to remove
				for i := range s.mercenarios {
					if s.mercenarios[i] == nombre { // replace .Name with the actual field name
						s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
						break
					}
				}
				if s.nmercenarios == 1 {
					fmt.Printf("Misión terminada\n")
					go func ()  {
						 time.Sleep(1 * time.Second)
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				} 
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}

		if req.GetDecision() == 3{
			if decisor <= prob3{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}

			if decisor > prob3{
				 
				fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
				s.nmercenarios -= 1
				nombre := strconv.Itoa(int(req.GetId())) // replace with the name you want to remove
				for i := range s.mercenarios {
					if s.mercenarios[i] == nombre { // replace .Name with the actual field name
						s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
						break
					}
				}
				if s.nmercenarios == 1 {
					fmt.Printf("Misión terminada\n")
					go func ()  {
						 time.Sleep(1 * time.Second)
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}
	}
	return nil, nil
}

func (s *server) Fase2(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) { //Implementa funcion fase2, le envia a cada mercenario si vive o muere
	if s.nmercenarios > 1 {
		mutex.Lock()
		s.decisions2 = append(s.decisions2, req.GetDecision())
		if len(s.decisions2) == s.nmercenarios {
			s.camino = rand.Intn(2) // 0 = A, 1 = B
			fmt.Printf("Comienza Piso 2!!\n")
			fmt.Printf("Mercenarios vivos:\n")
			for i := 0; i < len(s.mercenarios); i++ {
				fmt.Printf("%s\n", s.mercenarios[i])
			}
			s.fase += 1
			close(s.Señal2)
		}
		mutex.Unlock()
		
		<-s.Señal2 // Wait until all decisions are received
		mutex.Lock()
		defer mutex.Unlock()

		if s.camino == 0 && req.GetDecision() == 1{
			fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
			return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
		}
		if s.camino == 1 && req.GetDecision() == 1{
			
			fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
			s.nmercenarios -= 1
			nombre := strconv.Itoa(int(req.GetId()))
			for i := range s.mercenarios {
				if s.mercenarios[i] == nombre {
					s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
					break
				}
			}
			if s.nmercenarios == 1 {
				fmt.Printf("Misión terminada\n")
				go func ()  {
					time.Sleep(1 * time.Second)
					s.grpcServer.Stop()
				}()
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
			return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			
		}

		if s.camino == 0 && req.GetDecision() == 2{
			
			fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
			s.nmercenarios -= 1
			nombre := strconv.Itoa(int(req.GetId())) // replace with the name you want to remove
			for i := range s.mercenarios {
				if s.mercenarios[i] == nombre { // replace .Name with the actual field name
					s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
					break
				}
			}
			if s.nmercenarios == 1 {
				fmt.Printf("Misión terminada\n")
				go func ()  {
					time.Sleep(1 * time.Second)
					s.grpcServer.Stop()
				}()
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
			return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
		}
		

		if s.camino == 1 && req.GetDecision() == 2{
			fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
			return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
		}
	}
	return nil, nil
}

func (s *server) Fase3(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {
	if s.nmercenarios > 1 {
	mutex.Lock()
	defer mutex.Unlock()
	if _, exists := s.decisions3[int(req.GetId())]; !exists {
		s.decisions3[int(req.GetId())] = make([]int32, 0, 5)
	}
	for i := 0; i < 5; i++ {
		s.decisions3[int(req.GetId())] = append(s.decisions3[int(req.GetId())], req.Decisiones[i])
	}
	
		if len(s.decisions3) == s.nmercenarios && s.nmercenarios > 1{
			fmt.Printf("Comienza Piso 3!!\n")
			fmt.Printf("Mercenarios vivos:\n")
			for i := 0; i < len(s.mercenarios); i++ {
				fmt.Printf("%s\n", s.mercenarios[i])
			}
			fmt.Printf("Decisiones de los mercenarios:\n")
			for mercID, decisiones := range s.decisions3 {
				fmt.Printf("Mercenario %d: %v\n", mercID, decisiones)
			}

			ganadores := make([]int, 0)

			for i := 0; i < 5; i++ {
				s.decisor[i] = rand.Int31n(15) + 1
				fmt.Printf("El número elegido por Director es %d\n", s.decisor[i])
			}

			for mercID, decisiones := range s.decisions3 {
				
				aciertos := 0
				for i, decision := range decisiones {
					//fmt.Printf("Mercenario %d eligió %d y la opcion era %d\n", mercID, decision,s.decisor[i])
					if decision == s.decisor[i] {
						//fmt.Printf("Ha acertado un numero!")
						aciertos += 1
					}
				}
				if aciertos >= 2 {
					ganadores = append(ganadores, mercID)
					
				}
				
				if aciertos < 2 {
					fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
					s.nmercenarios -= 1
					fmt.Printf("Mercenarios restantes: %d\n",s.nmercenarios )
					if s.nmercenarios == 1 {
						fmt.Printf("Misión terminada\n")
						go func ()  {
							time.Sleep(1 * time.Second)
							s.grpcServer.Stop()
						}()
						return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
					}
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				
			}

			for _, ganadorID := range ganadores {
				if int32(ganadorID) == req.GetId() {
					fmt.Printf("Mercenario %d ha ganado\n", ganadorID)
					s.nmercenarios -= 1
					return &pb.DirectorMensaje{Estado: 1, Fase: 0}, nil
				}
			}
		}
	}
	return nil, nil
}




func main() {

	grpcServer := grpc.NewServer() //Se crea el servidor

	mutex.Lock()
	s := &server{ //Se le asignan los recursos al servidor
		fase: 0,
		grpcServer: grpcServer,
		decisions1: make([]int32, 0, 8),
		decisions2: make([]int32, 0, 8),
		decisions3: make(map[int][]int32),
		mercenarios: make([]string, 0, 8), 
		suma: 0,
		nmercenarios: 8,
		startSignalCh:  make(chan struct{}),
		Señal1: make(chan struct{}),
		Señal2: make(chan struct{}),
		Señal3: make(chan struct{}),
		decisor : make([]int32, 5),
		camino: 0,
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

}
