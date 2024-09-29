package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by skipping origin check
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan Message)           // Broadcast channel

// Message defines the structure of a chat message
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	IP       string `json:"ip"` // Adding IP to the message structure
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleHome serves the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// handleConnections upgrades initial GET requests to WebSockets
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// Get the client's IP address
	clientIP := r.RemoteAddr

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Set the IP address of the sender
		msg.IP = clientIP

		// Append special tags for specific usernames
		switch msg.Username {
		case "Mike", "Mark":
			msg.Username += " [VOLUNTEER]"
		case "Rich":
			msg.Username += " [MDPD]"
		case "Jake":
			msg.Username += " [MDFR]"
		case "Jamie":
			msg.Username += " [USCG]"
		}

		// Log the message and IP to the terminal
		fmt.Printf("Received message from %s (IP: %s): %s\n", msg.Username, msg.IP, msg.Content)

		// Send the message to the broadcast channel
		broadcast <- msg
	}
}

// handleMessages broadcasts incoming messages to all clients
func handleMessages() {
	for {
		msg := <-broadcast

		// Execute geminiLauncher.sh with the message content
		cmd := exec.Command("./geminiLauncher.sh", "message", msg.Content)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing geminiLauncher.sh: %v", err)
			// You might want to handle this error, perhaps by skipping this message
			continue
		}

		// Trim any whitespace from the output and update the message content
		msg.Content = strings.TrimSpace(string(output))

		// Broadcast the updated message to all clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
