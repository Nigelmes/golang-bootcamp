#!/bin/bash

if ! [[ $(docker image ls -q elasticsearch:7.9.2) ]]; then
  docker pull elasticsearch:7.9.2
fi

docker run -d --name elasticsearch-school21 \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  -v elasticsearch-golang-bootcamp:/usr/share/elasticsearch/data \
  elasticsearch:7.9.2