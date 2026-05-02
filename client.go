package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	Conn       net.Conn
	buf        [1024]byte
	server_buf [1024]byte
}

func (client *Client) connect() error {
	if conn, err := net.Dial("tcp", "localhost:8080"); err != nil {
		fmt.Printf("error occured! %v \n", err)
		return err
	} else {
		client.Conn = conn
		return nil
	}
}

func (client *Client) sendData(data []byte) (int, error) {
	if n, err := client.Conn.Write(data); err != nil {
		fmt.Printf("error occured! %v \n", err)
		return 0, err
	} else {
		return n, nil
	}
}

func (client *Client) readData() (int, error) {
	n, err := os.Stdin.Read(client.buf[:])
	return n, err
}

func (client *Client) RecieveFromServer() (int, error) {
	if n, err := client.Conn.Read(client.server_buf[:]); err == nil {
		return n, nil
	} else {
		return 0, err
	}
}

func (client *Client) HandleServerMessage() {
	for {
		n, err := client.RecieveFromServer()
		if err != nil {
			return
		}
		input := string(client.server_buf[:n])
		fmt.Printf("%s\n", input)
	}
}

func main() {

	client := Client{}

	if err := client.connect(); err != nil {
		fmt.Printf("error occured! %v \n", err)
		return
	} else {
		go client.HandleServerMessage()
		defer client.Conn.Close()
		for {
			n, err := client.readData()
			if err != nil {
				fmt.Printf("error occured! %v \n", err)
				return
			}
			input := strings.Fields(string(client.buf[:n]))
			if len(input) == 0 {
				continue
			}
			command := input[0]
			message := strings.Join(input[1:], " ")
			switch command {
			case "send":
				if _, err := client.sendData(client.buf[:n]); err != nil {
					fmt.Printf("error occured! %v \n", err)
					return
				} else {
					fmt.Printf("sent %s\n", message)
				}
			case "exit":
				if _, err := client.sendData(client.buf[:n]); err != nil {
					fmt.Printf("error occured! %v \n", err)
				} else {
					fmt.Printf("sent %s\n", string(client.buf[:n]))
					fmt.Println("exited!!")
				}
				return
			default:
				fmt.Println("incorrect input")
			}

		}
	}
}
