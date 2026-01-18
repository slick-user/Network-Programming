package main

import(
	"fmt"
	"net"
	"flag"
	"log"
	"bufio"
	"os"
	"strings"
)

var (
	host, path, method string
	port 							 int
)

func main() {
	// Initializing flags
	flag.StringVar(&method, "method", "GET", "HTTP method to use")
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.IntVar(&port, "port", 8080, "port to connect to")
	flag.StringVar(&path, "path", "/", "path to request")
	flag.Parse()

	ip, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, ip)
	if err != nil {
		panic(err)
	}

	log.Printf("connected to %s (@ %s)", host, conn.RemoteAddr())

	defer conn.Close()

	// Setting the required fileds to get the HTTP Req
	var reqfields = []string{
		fmt.Sprintf("%s %s HTPP/1.1", method, path),
		"Host: ", host,
		"User-Agent: httpget",
		"",
	}

	// This creates the request taking the fields required
	request := strings.Join(reqfields, "\r\n") + "\r\n"

	// Writes them to the connected address
	conn.Write([]byte(request))
	log.Printf("sent request:\n %s", request)

	// This should then write to the Conn
	for scanner := bufio.NewScanner(conn); scanner.Scan(); {
		line := scanner.Bytes()
		if _, err := fmt.Fprintf(os.Stdout, "%s \n", line); err != nil {
			log.Printf("error writing to connection: %s", err)
		}
		if scanner.Err() != nil {
			log.Printf("error writing to connection: %s", err)
			return
		}
	}
}
