#!/bin/bash

make build
docker build -t docker-dashserver .

docker build -t docker-proxy -f ./proxy.dockerfile .
docker-compose up -d --build

sleep 3

#./releases/dashclient
