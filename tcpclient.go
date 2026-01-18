package main

// This is a TCP Server in Go, this is the client side code to listen to the request that the server sends

// TCP Is listening on Port 6483, hosted on port 8080

// This is how importing is done in Go Nothing crazy, but yet totally its own thing
import(
	"fmt"
	"log"
	"bufio"	
	"flag"
	"net"
	"os"
)

// There is possibly a way to have any one anywhere connect to this. How will I figure that out?

func main() {
	const name = "writetcp"
	log.SetPrefix(name + "\t")

	// register the command-line flags -p specifies the port to connect to 
	// := is how variable dynamic type declarations work in Go, or you could be boring and go var x <type> = <whatever it is>
	port := flag.Int("p", 8080, "To connect to TCP Port")
	flag.Parse() 	// Parse Registered flags

	// you can define two things at once! Look at this it takes the 
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: *port})
	if err != nil {
		log.Fatalf("error connecting to localhost: %d %v", *port, err)
	}
	log.Print("connected to %v: will forward stdin", conn.RemoteAddr())

	defer conn.Close();

	go func() {
		/* Go Routines are used because they create an independent function that runs concurrently with other functions
		   allowing for parallelism, in case you didn't realise this is a GoRoutine */

		// This routine reads lines from the server and prints them to stdout
		// here it is reading the data from conn which is getting data from DialTCP
		for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {
			fmt.Printf("%s \n", connScanner.Text())

			if err := connScanner.Err(); err != nil {
				log.Fatalf("error from %s: %v", conn.RemoteAddr(), err)
			}
		}
	}()

	// Now for the writing part that is why we have this scanner Scan from os.Stdin
	for stdinScanner := bufio.NewScanner(os.Stdin); stdinScanner.Scan(); {
		log.Printf("sent: %s \n", stdinScanner.Text())
		if _, err := conn.Write(stdinScanner.Bytes()); err != nil {
			log.Fatalf("error writing to %s: %v", conn.RemoteAddr(), err)
		}

		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Fatalf("error writing to %s: %v", conn.RemoteAddr(), err)
		}

		if stdinScanner.Err() != nil {
			log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
		}		
	}

}

