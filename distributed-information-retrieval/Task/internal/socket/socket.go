package socket

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func ReadSocketMessage(conn net.Conn, buffer []byte) (int, error) {
	msgLen, err := conn.Read(buffer)
	if err != nil {
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			fmt.Println("Connection has been closed.")
		}

		fmt.Println("Error reading from connection:", err.Error())
		return 0, err
	}

	return msgLen, nil
}
