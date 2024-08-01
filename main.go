package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	server "server/functions"
)

const (
	defaultPort = "8989"
)

func main() {
	var port string
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) == 1 {
		port = defaultPort
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on the port :%s\n", port)

	// Add a signal handler for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		close(server.Messages)
		listener.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		server.Mutex.Lock()
		if len(server.Connections) >= server.MaxConnections {
			fmt.Fprintln(conn, "Maximum connections reached")
			conn.Close()
			server.Mutex.Unlock()
			continue
		}
		server.Mutex.Unlock()

		go server.HandleConnection(conn)
	}
}
