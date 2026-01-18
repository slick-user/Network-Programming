package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Reads lines from r then sets them to uppercase and writes them to w
func echoUpper(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("This is what the client wrote: %s", scanner.Text())

		// Fprintf writes to writer w
		fmt.Fprintf(w, "%s \n", strings.ToUpper(line))
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error: %s", err)
	}
	
}

func main() {
	const name = "tcpupperecho"
	log.SetPrefix(name + "\t")

	port := flag.Int("p", 8080, "port to listen on")
	flag.Parse()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Printf("listening at localHost: %s", listener.Addr())
	for {
		// this will loop to accept for connections

		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go echoUpper(conn, conn)
	}
}
