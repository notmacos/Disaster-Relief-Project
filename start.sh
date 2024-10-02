#!/bin/bash

# Run Node.js server
node server.js &

# Run Sqlite server
node sql/sql-interact.js &

# Run Go program
go run go-src/chat.go &

# Wait for both processes to finish
wait
