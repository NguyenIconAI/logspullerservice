package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/NguyenIconAI/logspullerservice/api"
	_ "github.com/NguyenIconAI/logspullerservice/docs" // Swagger docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Log Puller Service API
// @version 1.0
// @description This is a sample server for pulling log files.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func main() {
	port := flag.String("port", ":3000", "the listening port")
	flag.Parse()

	server := api.NewServer(*port)

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Printf("Server running on port %s\n", *port)
	log.Fatal(server.Start())
}
