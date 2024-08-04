package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	var clientName string
	fmt.Fprintln(conn, "Welcome to TCP-Chat!")
	fmt.Fprintln(conn, "         _nnnn_")
	fmt.Fprintln(conn, "        dGGGGMMb")
	fmt.Fprintln(conn, "       @p~qp~~qMb")
	fmt.Fprintln(conn, "       M|@||@) M|")
	fmt.Fprintln(conn, "       @,----.JM|")
	fmt.Fprintln(conn, "      JS^\\__/  qKL")
	fmt.Fprintln(conn, "     dZP        qKRb")
	fmt.Fprintln(conn, "    dZP          qKKb")
	fmt.Fprintln(conn, "   fZP            SMMb")
	fmt.Fprintln(conn, "   HZM            MMMM")
	fmt.Fprintln(conn, "   FqM            MMMM")
	fmt.Fprintln(conn, " __| \"\".        |\\dS\"qML")
	fmt.Fprintln(conn, " |    `.`       | `' \\Zq")
	fmt.Fprintln(conn, " _)      \\.___.,|     .'")
	fmt.Fprintln(conn, " \\____   )MMMMMP|   .'")
	fmt.Fprintln(conn, "     `-'       `--'")
	fmt.Fprintln(conn, "[ENTER YOUR NAME]:")

	scanner := bufio.NewScanner(conn)
	for {
		scanner.Scan()
		clientName = scanner.Text()

		if clientName == "" {
			fmt.Fprintln(conn, "Invalid client name. Please enter a valid name:")
			continue
		}

		Mutex.Lock()
		if _, exists := Connections[clientName]; exists {
			fmt.Fprintln(conn, "Client name already in use. Please choose another name:")
			Mutex.Unlock()
			continue
		}
		Connections[clientName] = conn
		Mutex.Unlock()
		break
	}

	fmt.Printf("Client %s connected\n", clientName)
	defer func() {
		fmt.Printf("Client %s disconnected\n", clientName)
	}()

	fmt.Fprintln(conn, "Connected to the chat!")
	

	// Send message history to new client
	Mutex.Lock()
	for _, msg := range MessageHistory {
		fmt.Fprintln(conn, msg)
	}
	Mutex.Unlock()

	go func() {
		for message := range Messages {
			Mutex.Lock()
			for _, conn := range Connections {
				fmt.Fprintln(conn, message)
			}
			Mutex.Unlock()
		}
	}()

	scanner = bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		message = fmt.Sprintf("[%s][%s]: %s", timestamp, clientName, message)
		Messages <- message

		// Store the message in history
		Mutex.Lock()
		MessageHistory = append(MessageHistory, message)
		Mutex.Unlock()
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from client %s: %v", clientName, err)
	}

	Mutex.Lock()
	delete(Connections, clientName)
	Mutex.Unlock()

	Messages <- fmt.Sprintf("%s has left the chat", clientName)
}
