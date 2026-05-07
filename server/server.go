package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Server struct {
	listener      net.Listener
	clients       map[string]net.Conn
	clients_mutex sync.Mutex
}

func (s *Server) Launch() error {
	if listener, err := net.Listen("tcp", "localhost:8080"); err != nil {
		fmt.Printf("server launch error: %v\n", err)
		return err
	} else {
		s.listener = listener
	}
	return nil
}

func (s *Server) Accept() (net.Conn, error) {

	if conn, err := s.listener.Accept(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return nil, err
	} else {
		addr := conn.RemoteAddr().String()
		s.clients_mutex.Lock()
		s.clients[addr] = conn
		s.clients_mutex.Unlock()
		fmt.Println("user connected!")
		return conn, nil
	}
}

func give_event(comand, message, addr string) string {
	var str string
	switch comand {
	case "exit":
		str = "user exitted!: " + addr
	case "send":
		str = "user " + addr + " sent: " + message
	}
	return str
}

func (s *Server) HandleClient(conn net.Conn) error {
	addr := conn.RemoteAddr().String()
	defer func() {
		s.clients_mutex.Lock()
		delete(s.clients, addr)
		s.clients_mutex.Unlock()
		conn.Close()
	}()

	pendingBuffer := make([]byte, 0, 1024)

	for {
		comand, message, err := s.ReadFromClient(conn, &pendingBuffer)
		if err != nil {
			fmt.Printf("error occurred: %v\n", err)
			return err
		}
		if comand == "" {
			continue
		}
		event := give_event(comand, message, addr)
		if event != "" {
			s.broadcast(event)
			fmt.Println(event)
		}
		if comand == "exit" {
			return nil
		}
	}
}

func (s *Server) ReadFromClient(conn net.Conn, pending *[]byte) (string, string, error) {

	var (
		readBuf     [1024]byte
		parse_input = make([]string, 0)
		line        = make([]byte, 0)
	)

	for i, b := range *pending {
		if b == '\n' {
			line = (*pending)[:i]
			*pending = (*pending)[i+1:]
			parse_input = strings.Fields((string(line[:])))
			if len(parse_input) == 0 {
				break
			}
			comand := parse_input[0]
			message := strings.Join(parse_input[1:], " ")
			return comand, message, nil
		}
	}

	n, err := conn.Read(readBuf[:])
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return "", "", err
	}

	for _, b := range readBuf[:n] {
		*pending = append(*pending, b)
	}

	for i, b := range *pending {
		if b == '\n' {
			line = (*pending)[:i]
			*pending = (*pending)[i+1:]
			parse_input = strings.Fields((string(line[:])))
			break
		}
	}

	if len(parse_input) == 0 {
		return "", "", nil
	}
	comand := parse_input[0]
	message := strings.Join(parse_input[1:], " ")
	return comand, message, err
}

func (s *Server) response(message string, conn net.Conn) (int, error) {
	n, err := conn.Write([]byte(message))
	return n, err
}

func (s *Server) broadcast(event string) error {
	s.clients_mutex.Lock()
	var temp_clients = make([]net.Conn, 0, len(s.clients))
	for _, conn := range s.clients {
		temp_clients = append(temp_clients, conn)
	}
	s.clients_mutex.Unlock()

	for _, conn := range temp_clients {
		_, err := s.response(event, conn)
		if err != nil {
			fmt.Printf("error occurred: %v\n", err)
			return err
		}
	}
	return nil
}

func main() {

	s := Server{clients: make(map[string]net.Conn)}
	err := s.Launch()
	if err != nil {
		return
	}
	for {
		conn, err := s.Accept()
		if err == nil {
			go s.HandleClient(conn)
		}

	}

}
