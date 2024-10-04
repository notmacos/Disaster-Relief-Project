#!/bin/bash

# Run Node.js server
node server.js &

# Run Sqlite server
node sql/sql-interact.js &

# Run Go program
cd go-src && go run . &

# Wait for all processes to finish
wait
