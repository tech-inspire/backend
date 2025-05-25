#!/bin/bash

KEYSPACE="posts"
SNAPSHOT_NAME="backup_$(date +%Y%m%d_%H%M%S)"
BACKUP_DIR="./backups/$SNAPSHOT_NAME"

CONTAINER_NAME_PREFIX="docker-"
CONTAINER_NAME_SUFFIX="-1"
NODES=(scylla-node1 scylla-node2 scylla-node3 scylla-node4 scylla-node5)

echo "Creating snapshot '$SNAPSHOT_NAME' for keyspace '$KEYSPACE' on all nodes..."

mkdir -p "$BACKUP_DIR"

for NODE in "${NODES[@]}"; do
  echo ">>> Taking snapshot on $NODE"
  CONTAINER=${CONTAINER_NAME_PREFIX}${NODE}${CONTAINER_NAME_SUFFIX}
  docker exec "$CONTAINER" nodetool snapshot -t "$SNAPSHOT_NAME" "$KEYSPACE"

  echo ">>> Copying snapshot data from $NODE"

  # Get list of snapshot directories for the keyspace
  TABLE_DIRS=$(docker exec "$CONTAINER" find /var/lib/scylla/data/$KEYSPACE -type d -path "*/snapshots/$SNAPSHOT_NAME")

  echo $TABLE_DIRS
  for SNAP_PATH in $TABLE_DIRS; do
    HOST_NODE_DIR="$BACKUP_DIR/$CONTAINER${SNAP_PATH#/var/lib/scylla}"
    mkdir -p "$HOST_NODE_DIR"

    echo ">> Copying from $CONTAINER:$SNAP_PATH to $HOST_NODE_DIR"
    docker cp "$CONTAINER:$SNAP_PATH" "$HOST_NODE_DIR"
  done
done

echo "✅ Backup complete. Files saved to: $BACKUP_DIR"
