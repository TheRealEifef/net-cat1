package server

import (
	"net"
	"sync"
)

const (
	MaxConnections = 10
)

var (
	Connections    = make(map[string]net.Conn)
	Messages       = make(chan string)
	Mutex          = &sync.Mutex{}
	MessageHistory []string
)
