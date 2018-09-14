package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := bufio.NewReader(conn)
	for {
		data, err := buf.ReadBytes(0x0A)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("----------")
		fmt.Println(string(data))
		fmt.Println("----------")
	}
}
