package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type User struct {
	Username string
}

type Client struct {
	Conn net.Conn
	Name string
}

func ConnectToServer(address string) {
	yellow := color.New(color.FgHiYellow).SprintFunc()
	purple := color.New(color.FgHiMagenta).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()

	// Prompt user for their username
	fmt.Print(yellow("Enter your username: "))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()
	if username == "" {
		fmt.Println("Username cannot be empty. Enter a username")
	}

	// Dial the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// Send username to server upon connection
	data := User{
		Username: username,
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshalling username: %v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, res)
	if err != nil {
		log.Fatalf("Error sending username to server: %v", err)
	}

	// Handle incoming messages from the server in a separate goroutine
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading from server: %v", err)
				os.Exit(0)
			}

			SetColor(purple, cyan, green, red, string(message))
		}
	}()

	// Read input from the user and send it to the server
	// scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error sending message to server: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from input: %v", err)
	}
}

func SetColor(purple, cyan, green, red func(a ...interface{}) string, message string) {
	result := strings.Split(message, " ")

	if result[len(result)-1] == "err" {
		result := result[:len(result)-1]
		msg := strings.Join(result, " ")
		fmt.Println(red(msg))
	} else if result[len(result)-1] == "noerr" {
		result := result[:len(result)-1]
		msg := strings.Join(result, " ")
		fmt.Println(green(msg))
	} else if result[len(result)-1] == "info" {
		// check if the second last element is "HIT"
		if result[len(result)-2] == "HIT" {
			result := result[:len(result)-1]
			msg := strings.Join(result, " ")
			fmt.Println(green(msg))
		} else {
			result := result[:len(result)-1]
			msg := strings.Join(result, " ")
			fmt.Println(red(msg))
		}
	} else if result[0] == "Welcome" {
		fmt.Println(cyan(message))
	} else {
		fmt.Println(purple(message))
	}

}
