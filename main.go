package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/NguyenIconAI/logspullerservice/api"
)

func main() {
	port := flag.String("port", ":3000", "the listening port")
	flag.Parse()

	server := api.NewServer(*port)

	fmt.Printf("Server running on port %s\n", *port)
	log.Fatal(server.Start())
}
