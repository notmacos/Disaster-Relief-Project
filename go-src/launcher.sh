#!/bin/bash

export API_KEY="gemini_key"

# Pass all command-line arguments to gemini.go
go run -mod=mod go-src/gemini.go "$@"

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
