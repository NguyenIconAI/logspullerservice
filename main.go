package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NguyenIconAI/logspullerservice/api"
	"github.com/NguyenIconAI/logspullerservice/constants"
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

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	port := flag.String("port", ":3000", "the listening port")
	apiKey := os.Getenv(constants.ApiKeyEnvVar)
	flag.Parse()

	if apiKey == "" {
		log.Fatal("API key is required")
	}
	server := api.NewServer(*port)

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Printf("Server running on port %s\n", *port)
	log.Fatal(server.Start())
}
