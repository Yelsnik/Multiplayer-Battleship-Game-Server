package main

import "game-client/client"

func main() {
	client.ConnectToServer("ws://localhost:8000/ws")
}
