#!/bin/bash

CONTAINER_NAME=docker-scylla-node1-1
SCYLLADB_USER=cassandra
SCYLLADB_PASSWORD=cassandra

docker exec -i $CONTAINER_NAME cqlsh -u $SCYLLADB_USER -p $SCYLLADB_PASSWORD -e "COPY posts.posts_by_id TO '/tmp/posts_by_id.csv' WITH HEADER = TRUE"
docker cp $CONTAINER_NAME:/tmp/posts_by_id.csv ./posts_by_id.csv
docker exec $CONTAINER_NAME rm /tmp/posts_by_id.csv
