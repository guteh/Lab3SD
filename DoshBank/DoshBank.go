package main

import (

	"fmt"
	"log"
	"os"
	"sync"

    
	"github.com/streadway/amqp"

)

var (
	mutex               = &sync.Mutex{}
	mercenariosMuertos  = make(map[int]bool)
	montoAcumulado      int
	archivoSalida       = "datos.txt"
	archivoSalidaHeader = "Mercenario | Fase | Monto Acumulado\n"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func consumeMessages() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mercenario_status", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			body := string(d.Body)
			fmt.Printf("Received a message: %s\n", body)

			var mercenarioID, fase int
			_, err := fmt.Sscanf(body, "Merceario %d falleció en el piso %d", &mercenarioID, &fase)
			if err == nil {
				mutex.Lock()
				if _, exists := mercenariosMuertos[mercenarioID]; !exists {
					mercenariosMuertos[mercenarioID] = true
					montoAcumulado += calcularMontoPorMercenario()
					guardarDatosEnArchivo(mercenarioID, fase, montoAcumulado)
				}
				mutex.Unlock()
			}
		}
	}()
}

func calcularMontoPorMercenario() int {
	// Implementar la lógica de cálculo del monto aquí
	return 100 // Ejemplo: cada mercenario muerto agrega 100 al monto acumulado
}

func guardarDatosEnArchivo(mercenarioID, fase, monto int) {
	f, err := os.OpenFile(archivoSalida, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer f.Close()

	// Si el archivo está vacío, escribir el encabezado primero
	fileInfo, _ := f.Stat()
	if fileInfo.Size() == 0 {
		if _, err := f.WriteString(archivoSalidaHeader); err != nil {
			log.Fatalf("Error al escribir el encabezado: %v", err)
		}
	}

	// Escribir los datos del mercenario en el archivo
	line := fmt.Sprintf("%d | %d | %d\n", mercenarioID, fase, monto)
	if _, err := f.WriteString(line); err != nil {
		log.Fatalf("Error al escribir en el archivo: %v", err)
	}
}

func main() {
	consumeMessages()

	// Bloquear el main para mantener el consumidor corriendo
	select {}
}
