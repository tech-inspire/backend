import threading
from concurrent import futures

import grpc
from embeddings.v1 import embeddings_pb2, embeddings_pb2_grpc

from app.services.clip_embedder import embed_text


class EmbeddingsServiceServicer(embeddings_pb2_grpc.EmbeddingsServiceServicer):
    def GenerateTextEmbeddings(self, request, context):
        text = request.text
        vec = embed_text(text)
        return embeddings_pb2.GenerateTextEmbeddingsResponse(
            embedding_vector=vec.astype(float)
        )


def start_grpc_server():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    embeddings_pb2_grpc.add_EmbeddingsServiceServicer_to_server(
        EmbeddingsServiceServicer(), server
    )

    server.add_insecure_port("[::]:50051")
    thread = threading.Thread(target=server.start)
    thread.start()
    print("gRPC server running on port 50051")
    return server
