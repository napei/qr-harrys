#!/bin/bash

docker-compose pull
docker-compose up --force-recreate --remove-orphans -d
docker image prune -f
