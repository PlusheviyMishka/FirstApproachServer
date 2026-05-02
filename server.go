package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type id int

type Server struct {
	listener      net.Listener
	clients       map[id]net.Conn
	nextid        id
	clients_mutex sync.Mutex
}

func (s *Server) Launch() error {
	if listener, err := net.Listen("tcp", "0.0.0.0:8080"); err != nil {
		fmt.Printf("server launch error: %v\n", err)
		return err
	} else {
		s.listener = listener
	}
	return nil
}

func (s *Server) Accept() (id, net.Conn, error) {

	if conn, err := s.listener.Accept(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return -1, nil, err
	} else {
		s.clients_mutex.Lock()
		s.clients[s.nextid] = conn
		id := s.nextid
		s.nextid++
		s.clients_mutex.Unlock()
		fmt.Println("user connected!")
		return id, conn, nil
	}
}

func (s *Server) HandleClient(id id, conn net.Conn) error {

	defer func() {
		s.clients_mutex.Lock()
		delete(s.clients, id)
		s.clients_mutex.Unlock()
		conn.Close()
	}()
	for {
		comand, message, err := s.ReadFromClient(id, conn)
		if err != nil {
			fmt.Printf("error occurred: %v\n", err)
			return err
		}
		switch comand {
		case "exit":
			fmt.Printf("user exitted!: %d\n", id)
			return nil
		case "send":
			fmt.Printf("user %d sent: %s\n", id, message)
		case "":
			fmt.Printf("user %d sent: empty string\n", id)
		}
	}
}

func (s *Server) ReadFromClient(id id, conn net.Conn) (string, string, error) {

	var input [1024]byte

	n, err := conn.Read(input[:])
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return "", "", err
	}
	if n == 0 {
		return "", "", nil
	}
	parse_input := strings.Fields((string(input[:n])))
	if len(parse_input) == 0 {
		return "", "", nil
	}
	comand := parse_input[0]
	message := strings.Join(parse_input[1:], " ")
	return comand, message, err
}

func main() {

	s := Server{clients: make(map[id]net.Conn)}
	s.Launch()

	for {
		id, conn, err := s.Accept()
		if err == nil {
			go s.HandleClient(id, conn)
		}
	}

}
