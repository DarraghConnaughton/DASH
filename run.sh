#!/bin/bash

make build

docker build -t docker-dashserver .
docker run -p 8080:8080 -d -it --name  docker-dashserver-container docker-dashserver

# Write file indicating a successful initialisation to the container file system once present, continue with client.
sleep 3

./releases/dashclient