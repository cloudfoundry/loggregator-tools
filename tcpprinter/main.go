package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"code.cloudfoundry.org/rfc5424"
)

func main() {
	l, err := net.Listen("tcp4", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Print("Listening on " + os.Getenv("PORT"))

	for {
		conn, err := l.Accept()
		log.Printf("Accepted connection")
		if err != nil {
			log.Printf("Error accepting: %s", err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	var msg rfc5424.Message
	for {
		_, err := msg.ReadFrom(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s", string(msg.Message))
	}
}
