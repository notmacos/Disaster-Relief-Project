#!/bin/bash

export API_KEY="VERY IMPORTANT!"

# Pass all command-line arguments to eventRecommendations.go
go run eventRecommendations.go "$@"
