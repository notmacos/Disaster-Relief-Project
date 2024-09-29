# Disaster Management Dashboard Project

This project consists of a Node.js server, a SQLite database interaction module, and a Go-based WebSocket chat server.

## Components

1. **Node.js Server** (`server.js`): Serves static files from the 'public' directory.
2. **SQLite Database Interaction** (`sql/sql-interact.js`): Handles database operations for events.
3. **Go WebSocket Chat Server** (`go-chat-test/main.go`): Implements a real-time chat functionality.

## Prerequisites

- Node.js
- Go
- SQLite

## Getting Started

1. Clone the repository:
   ```
   git clone <repository-url>
   cd <project-directory>
   ```

2. Install Node.js dependencies:
   ```
   npm install express sqlite3
   ```

3. Install Go dependencies:
   ```
   go get github.com/gorilla/websocket
   ```

4. Set up the SQLite database:
   ```
   sqlite3 shellhacks.db
   CREATE TABLE events (ID INTEGER PRIMARY KEY AUTOINCREMENT, eventType TEXT, location TEXT, time TEXT, poster TEXT, description TEXT);
   .exit
   ```

5. Start the servers:
   ```
   chmod +x start.sh
   ./start.sh
   ```

   This will start both the Node.js server, Sqlite server, and the Go chat server.

## Usage

### Node.js Server

- The server runs on `http://localhost:3000`
- Static files are served from the `public` directory

### SQLite Database Interaction

- The server runs on `http://localhost:1400`

You can interact with the database using the following commands:

- Add an event:
  ```
  node sql/sql-interact.js add <eventType> <location> <time> <poster> <description>
  ```

- Remove an event:
  ```
  node sql/sql-interact.js remove <eventId>
  ```

- Read all events:
  ```
  node sql/sql-interact.js read
  ```

### Go WebSocket Chat Server

- The chat server runs on `http://localhost:8080`
- Connect to the WebSocket at `ws://localhost:8080/ws`
- The `index.html` file (not provided in the current files) should implement the client-side chat interface

## API Endpoints

- GET `/api/events`: Retrieves all events from the database

## Notes

- The chat server adds special tags to usernames:
  - "mike" and "mark" are tagged as [VOLUNTEER]
  - "rich" is tagged as [LOCAL PD]

## Contributing

[Add contribution guidelines here]

## License

[Add license information here]
