#!/bin/bash

docker pull postgres:13.1
if [ ! "$(docker ps -q -f name=pg-shortener)" ]; then
  if [ "$(docker ps -aq -f status=exited -f name=pg-shortener)" ]; then
    docker rm pg-shortener
  fi
  docker run --name=pg-shortener -p 5432:5432 -v shortener-db-pg:/var/lib/postgresql/data \
   -e POSTGRES_DB='shortener-db' -e POSTGRES_PASSWORD='1110' -d postgres:13.1
  ss -tnlp | grep 5432
fi
