#!/bin/bash

# Run Node.js server
node server.js &

# Run Sqlite server
node SQL-Interaction/sql-interact.js &

# Run Go program
cd go-chat-test
go run main.go &

# Wait for both processes to finish
wait
