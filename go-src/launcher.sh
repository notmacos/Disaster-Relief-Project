#!/bin/bash

export API_KEY="AIzaSyB9ZQZ4K9qd_r6Xjm2FZ69RYjbusagiwmQ"

# Pass all command-line arguments to gemini.go
go run gemini.go "$@"

# Capture the exit status
exit_status=$?

# If the command failed, print more detailed information
if [ $exit_status -ne 0 ]; then
    echo "Error: go run command failed with exit status $exit_status"
    echo "GOPATH: $GOPATH"
    echo "GO111MODULE: $GO111MODULE"
    go version
    go env
fi

# Exit with the same status as the go run command
exit $exit_status
