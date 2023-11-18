package client

import (
	"Task/internal"
	"Task/internal/repositories"
	"Task/internal/socket"
	"Task/pkg/signal"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var clientDataset map[string]repositories.Data

func StartClient() {
	fmt.Println(">> Client Started.")

	// establish connection
	connection, err := net.Dial(internal.ServerType, internal.ServerHost+":"+internal.ServerPort)
	if err != nil {
		panic(err)
	}

	fmt.Println(">> Connected to server.")

	defer connection.Close()

	receiveDataset(connection)

	go handleUserInput(connection)
	go handleClientSocketConnection(connection)

	signal.Wait()
}

func receiveDataset(connection net.Conn) {
	fmt.Println(">> Waiting for data...")

	clientDataset = make(map[string]repositories.Data)

	started := false
	for {
		buffer := make([]byte, 1024*1024)
		msgLen, err := socket.ReadSocketMessage(connection, buffer)
		if err != nil {
			return
		}

		receivedData := string(buffer[:msgLen])
		lines := strings.Split(receivedData, "\n")

		for _, line := range lines {
			if line == "" {
				continue
			}

			if !started && line != "START_DATA" {
				fmt.Println("\tData should start with 'START_DATA'")
				os.Exit(1)
			}

			if !started {
				fmt.Println("\tStart receiving data")
				started = true
				continue
			}

			if line == "END_DATA" {
				fmt.Println("\tEnd of receiving data")
				fmt.Println("\t"+strconv.Itoa(len(clientDataset)), "rows of data added.")
				return
			}

			data := repositories.CommaStringToData(line)
			clientDataset[data.ID] = data
		}
	}

}

func handleUserInput(conn net.Conn) {
	for {
		var IDs string
		fmt.Scanln(&IDs)

		for _, ID := range strings.Split(IDs, ",") {
			fmt.Println("\tSearching for ID " + ID + "...")
			ID = "'" + ID + "'"
			data, exist := clientDataset[ID]
			if exist {
				printResult(data)
			} else {
				// Request server for data
				conn.Write([]byte("REQUEST_ID:" + ID + "\n"))
			}
		}
	}
}

func handleClientSocketConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1024*1024)
		msgLen, err := socket.ReadSocketMessage(conn, buffer)
		if err != nil {
			return
		}

		receivedData := string(buffer[:msgLen])
		for _, line := range strings.Split(receivedData, "\n") {
			if strings.Contains(line, "REQUEST_ID") {
				row := strings.Split(line, ":")
				responseFrom := row[0]
				requestFrom := row[1]
				ID := row[3]
				data, exist := clientDataset[ID]
				if exist {
					conn.Write([]byte("RESPONSE:" + responseFrom + ":" + requestFrom + ":" + data.ToString() + "\n"))
				}
			} else if strings.Contains(line, "RESPONSE") {
				row := strings.Split(line, ":")
				data := repositories.CommaStringToData(row[1])
				printResult(data)
			}
		}
	}
}

func printResult(data repositories.Data) {
	fmt.Println("* Data of ID "+data.ID+" is :", data)
}
