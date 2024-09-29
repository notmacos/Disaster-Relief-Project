#!/bin/bash

export API_KEY="AIzaSyCQwhCrrOQKZflnqOBmm6iQQ_F80Mnw25k"

# Pass all command-line arguments to eventRecommendations.go
go run eventRecommendations.go "$@"
