package main

import (
	"eyeident/internal/server"
	"log"
)

func main() {
	srv := server.NewServer()
	log.Println("Server running at http://192.168.1.105:8080")
	log.Println("Server running at http://172.21.55.18:8080")

	//_, err := db.ConnectPostgres()
	//if err != nil {
	//	return
	//}

	log.Fatal(srv.Run(":8080"))
}
