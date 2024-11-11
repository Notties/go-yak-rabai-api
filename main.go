// main.go
package main

import (
	"fmt"
	"os"

	"yak.rabai/config"
	router "yak.rabai/routes"
)

func main() {
	config.ConnectDatabase()
	config.ConnectRedis()

	// Access environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	fmt.Println("DB Host:", dbHost)
	fmt.Println("DB Port:", dbPort)

	r := router.SetupRouter()

	r.Run(":8080")
}
