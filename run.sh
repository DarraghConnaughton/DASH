#!/bin/bash

make build
docker build -t docker-dashserver -f ./dashserver.dockerfile .
docker build -t docker-controlproxy -f ./controlproxy.dockerfile .
docker build -t docker-forwardproxy -f ./forwardproxy.dockerfile .
docker-compose up -d --build

sleep 3

#./releases/dashclient
