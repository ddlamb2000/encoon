#!/bin/sh
docker network create encoon-network
docker network connect encoon-network postgres-encoon 
docker network connect encoon-network encoon