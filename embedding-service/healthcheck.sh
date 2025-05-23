#!/bin/sh

PORT=$(echo "$SERVER_METRICS_ADDRESS" | cut -d ':' -f2)

PORT=${PORT:-50051}

URL="localhost:$PORT"

/usr/local/bin/grpcurl -plaintext -d '{"text": "ping"}' ${URL} embeddings.v1.EmbeddingsService/GenerateTextEmbeddings

