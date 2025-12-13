package main

import (
	"eyeident/internal/db"
	"eyeident/internal/server"
	"eyeident/internal/workers"
	"log"
	"time"
)

func main() {
	srv := server.NewServer()
	log.Println("Server running at http://192.168.1.105:8080")
	log.Println("Server running at http://172.21.55.18:8080")

	dbPool, err := db.ConnectPostgres()
	if err != nil {
		return
	}
	defer dbPool.Close()

	worker := workers.NewWorker("DataMigratorTask")
	scheduler := workers.NewScheduler(worker, time.Minute*5)
	scheduler.Start()

	log.Fatal(srv.Run(":8080"))
}
