package main

import (
	"log"
	"stats_service_core/db_utils"
	"stats_service_core/http_server"
	"stats_service_core/kafka_pumper"
	"stats_service_core/stats_server"
)

func main() {
	err := db_utils.StartUpDB()
	for err != nil {
		log.Println(err, "retrying...")
		err = db_utils.StartUpDB()
	}

	kafka_pumper.StartPumping()
	go http_server.RunServer()
	stats_server.RunServer()
}
