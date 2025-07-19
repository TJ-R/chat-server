package main

import (
	"net"
	"log"
)

func main() {
	l, err := net.Listen("tcp", "localhost:1234")	
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on ws://%v", l.Addr())
}
