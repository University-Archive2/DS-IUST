package server

import (
	"Task/internal"
	"Task/internal/repositories"
	"Task/internal/socket"
	"Task/pkg/signal"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var dataset []repositories.Data
var clients []net.Conn

func StartServer() {
	fmt.Println("Server Started.")

	// Load the dataset from the CSV
	fmt.Println(">> Loading data...")
	dataset = repositories.LoadCSV()
	fmt.Println("\t", len(dataset), "rows of data added.")

	server, err := net.Listen(internal.ServerType, internal.ServerHost+":"+internal.ServerPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer server.Close()

	fmt.Println(">> Listening on " + internal.ServerHost + ":" + internal.ServerPort)

	acceptClients(server)

	distributeData()

	for i, connection := range clients {
		go handleServerSocketConnections(connection, i)
	}

	signal.Wait()
}

func acceptClients(server net.Listener) {
	fmt.Println(">> Waiting 10 seconds for clients to connect...")

	// Wait for clients before distributing the data
	timer := time.NewTimer(10 * time.Second)

	isExpired := false

	go func(isExpired *bool) {
		for {
			connection, err := server.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					fmt.Println("Server is shutting down...")
					return
				}

				fmt.Println("Error accepting:", err.Error())
				os.Exit(1)
			}

			if !*isExpired {
				fmt.Printf("\tClient [%v]: %s connected.\n", len(clients), connection.RemoteAddr().String())
				clients = append(clients, connection)
			} else {
				return
			}
		}
	}(&isExpired)

	<-timer.C
	isExpired = true
}

// Distribute data among connected clients
func distributeData() {
	fmt.Println(">> Start distributing data...")

	datasetToDistribute := dataset

	totalClients := len(clients)

	remainingPercentage := 100

	for i, conn := range clients {
		var percentage int

		if i == totalClients-1 { // if it's the last client, assign remaining percentage
			percentage = remainingPercentage
		} else {
			// Calculate random percentage for this client
			percentage = rand.Intn(remainingPercentage-(totalClients-i-1)) + 1 // Ensure that there's enough percentage left for the remaining clients
			remainingPercentage -= percentage
		}

		chunkSize := percentage * len(dataset) / 100
		chunk := datasetToDistribute[:chunkSize]
		datasetToDistribute = datasetToDistribute[chunkSize:]

		fmt.Printf("\tSending %v rows of data to client %v\n", len(chunk), i)

		conn.Write([]byte("START_DATA\n"))

		// Convert chunk data into bytes and send
		// For simplicity, using a delimiter to join and send the data
		for _, data := range chunk {
			conn.Write([]byte(data.ToString() + "\n"))
		}

		conn.Write([]byte("END_DATA"))
	}

	fmt.Println("\tData distributed among clients.")
}

func handleServerSocketConnections(connection net.Conn, i int) {
	for {
		buffer := make([]byte, 1024*1024)
		msgLen, err := socket.ReadSocketMessage(connection, buffer)
		if err != nil {
			return
		}

		receivedData := string(buffer[:msgLen])

		for _, line := range strings.Split(receivedData, "\n") {
			if strings.Contains(line, "REQUEST_ID") {
				broadcastRequest(line, i)
			} else if strings.Contains(line, "RESPONSE") {
				row := strings.Split(line, ":")
				responseFrom := row[1]
				requestFrom := row[2]
				data := repositories.CommaStringToData(row[3])
				fmt.Printf("Data of ID %s requested from %s found in %s\n", data.ID, requestFrom, responseFrom)
				sendMessage(fmt.Sprintf("RESPONSE:%s\n", data.ToString()), requestFrom)
			}
		}
	}
}

func broadcastRequest(message string, from int) {
	for to, connection := range clients {
		if from == to {
			continue
		}

		connection.Write([]byte(fmt.Sprintf("%v:%v:%s\n", strconv.Itoa(to), strconv.Itoa(from), message)))
	}
}

func sendMessage(message string, to string) {
	i, _ := strconv.Atoi(to)
	clients[i].Write([]byte(message))
}
