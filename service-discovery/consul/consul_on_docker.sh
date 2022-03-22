#!/bin/zsh

# start a consul server
docker run -d --name=dev-consul -p 8500:8500 -e CONSUL_BIND_INTERFACE=eth0 consul

# start two consul agent
docker run -d --name=dev-consul1 -e CONSUL_BIND_INTERFACE=eth0 consul agent -dev -join=172.17.0.2

docker run -d --name=dev-consul2 -e CONSUL_BIND_INTERFACE=eth0 consul agent -dev -join=172.17.0.2

