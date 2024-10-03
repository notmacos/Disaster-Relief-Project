package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

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

var (
	lastMessageTime = make(map[string]time.Time)
	rateLimitMutex  sync.Mutex
)

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	// Define a regular expression for allowed characters
	allowedChars := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?'"()\-:;]+$`)

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Rate limiting
		rateLimitMutex.Lock()
		lastTime, exists := lastMessageTime[clientIP]
		if exists && time.Since(lastTime) < 5*time.Second {
			rateLimitMutex.Unlock()
			log.Printf("Rate limited message from %s (IP: %s)", msg.Username, clientIP)
			continue // Skip this message and wait for the next one
		}
		lastMessageTime[clientIP] = time.Now()
		rateLimitMutex.Unlock()

		// Check if the message content exceeds the character limit
		if len(msg.Content) > 84 {
			log.Printf("Ignored message exceeding character limit from %s (IP: %s)", msg.Username, clientIP)
			continue // Skip this message and wait for the next one
		}

		// Check if the message content contains only allowed characters
		if !allowedChars.MatchString(msg.Content) {
			log.Printf("Ignored message with invalid characters from %s (IP: %s)", msg.Username, clientIP)
			continue // Skip this message and wait for the next one
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

		// Execute launcher.sh with the message content
		cmd := exec.Command("./launcher.sh", "message", msg.Content)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing launcher.sh: %v", err)
			msg.Content = "[message restricted]"
		} else {
			// Trim any whitespace from the output and update the message content
			msg.Content = strings.TrimSpace(string(output))
		}

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
