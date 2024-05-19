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
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)
var mutex = &sync.Mutex{} //Se crea un mutex para evitar problemas de concurrencia



type server struct {  //Crea el servidor rcp con sus variables globales
    pb.UnimplementedMercDirServer
	pb.UnimplementedNameDataServer
	fase int
	suma int
	nmercenarios int
	grpcServer *grpc.Server
	grpcServer1 *grpc.Server
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
	NameNode pb.DirNameClient

	
}

//Implementa la funcion SolicitarM de la interfaz ServicioRecursos de RCP


func (s *server) SolicitarM(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Mercenario %d ha llegado!\n", req.GetId())
	s.mercenarios = append(s.mercenarios, strconv.Itoa(int(req.GetId()))) //Se agrega a lista de mercenarios
	s.suma += 1

	if s.suma == s.nmercenarios { //Si ya llegaron todos los mercenarios
		s.fase += 1
		fmt.Printf("Se empieza mision! Fase actual: %d\n", s.fase)
		close(s.startSignalCh) //Se da se;al de inicio
	}

	return &pb.DirectorMensaje{Inicio: 1, Fase: int32(s.fase)}, nil
}

func (s *server) IniciarMision(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {
	<-s.startSignalCh // Se inicia la mision
	return &pb.DirectorMensaje{Inicio: 1, Fase: int32(s.fase)}, nil
}




func (s *server) Fase1(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {  //Implementa funcion fase1, le envia a cada mercenario si vive o muere
	mutex.Lock()
	s.decisions1 = append(s.decisions1, req.GetDecision()) //Agrego decision a lista con todas las deciisones
	s.NameNode.RegistrosDirector(context.Background(), &pb.EnviarDecision{Nombre: strconv.Itoa(int(req.GetId())), Piso: 1, Decision: req.GetDecision()}) //Envio decision al NameNode
	//fmt.Printf("Mercenario %d ha enviado su decision: %d\n", req.GetId(), req.GetDecision())
	if len(s.decisions1) == s.nmercenarios { //Si ya llegaron todas las decisiones
		for s.X == s.Y{ //Genero valores X e Y globales, para que todos los mercenarios tengan los mismo valores
			s.X = rand.Intn(100)
			s.Y = rand.Intn(100)
		}
		close(s.Señal1) //Libero señal de inicio
	}
	mutex.Unlock()
	
	<-s.Señal1 // Espero a que lleguen todas las decisiones y se;al de inicio
	mutex.Lock()
	defer mutex.Unlock()

	

	decisor := rand.Intn(100) //Numero que indica si vive o muere

	if s.X < s.Y { 
		prob1 := s.X
		prob2 := s.Y - s.X
		prob3 := 100 - s.Y

		if req.GetDecision() == 1{ //Si eligio el arma 1
			if decisor <= prob1{  //Gano la probabildiad
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}

			if decisor > prob1{ //Perdio la probabilidad
				 
				fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
				s.nmercenarios -= 1 //Decremento la cantidad de mercenarios
				nombre := strconv.Itoa(int(req.GetId())) //Quito mercenario de la lista
				for i := range s.mercenarios {
					if s.mercenarios[i] == nombre { 
						s.mercenarios = append(s.mercenarios[:i], s.mercenarios[i+1:]...)
						break
					}
				}
				if s.nmercenarios == 1 { //Si solo queda un mercenario
					go func ()  {
						time.Sleep(1 * time.Second) //Espero un segundo para cerrar el servidor
						fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil //Devuelvo que termino la mision
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil //Devuelvo que murio
				
			}
		}

		if req.GetDecision() == 2{	 //Elgii arma 2
			//Misma logica de antes con probabildiades
			if decisor <= prob2{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil

			}

			if decisor > prob2{
				 
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
					go func ()  {
						time.Sleep(1 * time.Second)
						fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}

		if req.GetDecision() == 3{ //Eligio arma 3
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
					go func ()  {
						time.Sleep(1 * time.Second)
						fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
						s.grpcServer.Stop()
						}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}
	}
	if s.Y <= s.X { //En el caso que Y sea menor o igual a X, se cambian las probabilidades pero misma logica
		prob1 := s.Y
		prob2 := s.X - s.Y
		prob3 := 100 - s.X

		if req.GetDecision() == 1{ //Arma 1
			if decisor <= prob1{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}
			if decisor > prob1{
				 
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
					go func ()  {
						time.Sleep(1 * time.Second)
						fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				}
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}
		
		if req.GetDecision() == 2{ //Arma 2
			if decisor <= prob2{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}

			if decisor > prob2{
				 
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
					go func ()  {
						time.Sleep(1 * time.Second)
						fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
						s.grpcServer.Stop()
					}()
					return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
				} 
				return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
			}
		}

		if req.GetDecision() == 3{ //Arma 3
			if decisor <= prob3{
				fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
				return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
			}

			if decisor > prob3{
				 
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
					go func ()  {
						time.Sleep(1 * time.Second)
						fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
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

	mutex.Lock()
	s.decisions2 = append(s.decisions2, req.GetDecision())
	s.NameNode.RegistrosDirector(context.Background(), &pb.EnviarDecision{Nombre: strconv.Itoa(int(req.GetId())), Piso: 2, Decision: req.GetDecision()}) //Envio decision al NameNode
	if len(s.decisions2) == s.nmercenarios && s.nmercenarios > 1{
		s.camino = rand.Intn(2) // 0 = A, 1 = B Eligo entre 2 caminos cual va a ser el correcto
		fmt.Printf("Comienza Piso 2!!\n")
		fmt.Printf("Mercenarios vivos:\n")
		for i := 0; i < len(s.mercenarios); i++ {
			fmt.Printf("%s\n", s.mercenarios[i])
		}
		s.fase += 1
		close(s.Señal2)
	}
	mutex.Unlock()
	if s.nmercenarios == 1 {
		go func ()  {
			time.Sleep(1 * time.Second)
			fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
			s.grpcServer.Stop()
		}()
		return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
	}
	
	<-s.Señal2 
	mutex.Lock()
	defer mutex.Unlock()

	if s.camino == 0 && req.GetDecision() == 1{ //Si eligio el camino A y era el A
		fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
		return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
	}
	if s.camino == 1 && req.GetDecision() == 1{ //Si eligio el camino A y era el B
		
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
			go func ()  {
				time.Sleep(1 * time.Second)
				fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
				s.grpcServer.Stop()
			}()
			return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
		}
		return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
		
	}

	if s.camino == 0 && req.GetDecision() == 2{ //Si eligio el camino B y era el A
		
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
			go func ()  {
				time.Sleep(1 * time.Second)
				fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
				s.grpcServer.Stop()
			}()
			return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
		}
		return &pb.DirectorMensaje{Estado: 0, Fase: 0}, nil
	}
	

	if s.camino == 1 && req.GetDecision() == 2{ //Si eligio el camino B y era el B
		fmt.Printf("Mercenario %d ha sobrevivido\n", req.GetId())
		return &pb.DirectorMensaje{Estado: 1, Fase: int32(s.fase)}, nil
	}
	
	return nil, nil
}

func (s *server) Fase3(ctx context.Context, req *pb.MercenarioMensaje) (*pb.DirectorMensaje, error) {
	if s.nmercenarios > 1 {
		mutex.Lock()
		defer mutex.Unlock()
		if _, exists := s.decisions3[int(req.GetId())]; !exists { //Si no existe el mercenario en el mapa de decisiones lo agrego
			s.decisions3[int(req.GetId())] = make([]int32, 0, 5)
		}
		s.NameNode.RegistrosDirector(context.Background(), &pb.EnviarDecision{Nombre: strconv.Itoa(int(req.GetId())), Piso: 3, Decisiones: req.GetDecisiones()}) //Envio decision al NameNode
		for i := 0; i < 5; i++ { //Agrego las decisiones del mercenario
			s.decisions3[int(req.GetId())] = append(s.decisions3[int(req.GetId())], req.Decisiones[i])
		}
		
			if len(s.decisions3) == s.nmercenarios && s.nmercenarios > 1{ //Si ya llegaron todas las decisiones
				fmt.Printf("Comienza Piso 3!!\n")
				fmt.Printf("Mercenarios vivos:\n")
				for i := 0; i < len(s.mercenarios); i++ {  //Imprimo los mercenarios vivos
					fmt.Printf("%s\n", s.mercenarios[i])
				}
				for mercID, decisiones := range s.decisions3 { //Imprimo las decisiones de los mercenarios
					fmt.Printf("Mercenario %d: %v\n", mercID, decisiones)
				}

				ganadores := make([]int, 0) //Lista de ganadores para almacenarlos

				for i := 0; i < 5; i++ {
					s.decisor[i] = rand.Int31n(15) + 1
					fmt.Printf("El número elegido por Director es %d\n", s.decisor[i]) //Genero los numeros del director
				}

				for mercID, decisiones := range s.decisions3 { //Comparo las decisiones de los mercenarios con las del director
					
					aciertos := 0 //Cantidad de aciertos
					for i, decision := range decisiones { //Recorro las decisiones
						//fmt.Printf("Mercenario %d eligió %d y la opcion era %d\n", mercID, decision,s.decisor[i])
						if decision == s.decisor[i] { //Si acierta
							aciertos += 1 //Incremento aciertos
						}
					}
					if aciertos >= 2 { //Si acierta 2 o mas veces lo agrego a la lista de ganadores
						ganadores = append(ganadores, mercID)
						
					}
					
					if aciertos < 2 { //Si no acierta 2 o mas veces muere
						fmt.Printf("Mercenario %d ha muerto\n", req.GetId())
						s.nmercenarios -= 1
						fmt.Printf("Mercenarios restantes: %d\n",s.nmercenarios )
						if s.nmercenarios == 1 {
							go func ()  {
								time.Sleep(1 * time.Second)
								fmt.Printf("Misión terminada por falta de mercenarios, se le entregara el monto al sobreviviente\n")
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

func StartServerMerc(s *server, grpcServer *grpc.Server){
	ip := "10.35.169.91:8080"
	pb.RegisterMercDirServer(grpcServer, s) //Se registra el servidor
	

	 //Se asigna la direccion del servidor
	lis, err := net.Listen("tcp", ip) //Se crea el listener
	
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func StartServerData(s *server, grpcServer *grpc.Server){
	fmt.Printf("DataNode\n")
	ip := "10.35.169.91:8084"
	pb.RegisterNameDataServer(grpcServer, s) //Se registra el servidor
	 //Se asigna la direccion del servidor
	lis, err := net.Listen("tcp", ip) //Se crea el listener
	
    if err != nil {
		log.Fatalf("Fallo al escuchar %v", err)
    }
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
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


func main() {

	//Conexion a NameNode
	conn, err := grpc.NewClient("10.35.169.93:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))  //10.35.169.93:8080
    if err != nil {
		log.Fatalf("Fallo al conectarse a NameNode: %v", err)
    }
    defer conn.Close()
    NameNode := pb.NewDirNameClient(conn)
	


	grpcServer := grpc.NewServer() //Se crea el servidor
	grpcServer1 := grpc.NewServer() //Se crea el servidor
	
	mutex.Lock()
	s := &server{ //Se le asignan los recursos al servidor
		fase: 0,
		grpcServer: grpcServer,
		grpcServer1: grpcServer1,
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
		NameNode: NameNode,
    }
	
	mutex.Unlock()
	
	//Grpc Mercenarios - Director
	go StartServerMerc(s, grpcServer)

	//Grpc DataNode
	go StartServerData(s, grpcServer1)

	time.Sleep(20 * time.Second)

	
	
	
}
