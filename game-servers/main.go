package main

import (
	"game-sever/server"
	"game-sever/utils"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {

	// Create a new server.
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	server := server.NewServer(config.GameSeverAddress, rdb)

	// Start the server.
	server.Start()

}
