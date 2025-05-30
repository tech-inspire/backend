import threading
from concurrent import futures

import grpc
from embeddings.v1 import embeddings_pb2, embeddings_pb2_grpc
from grpc_reflection.v1alpha import reflection
from app.services.clip_embedder import ClipEmbedder
from grpc_health.v1 import health
from grpc_health.v1 import health_pb2_grpc


class EmbeddingsServiceServicer(embeddings_pb2_grpc.EmbeddingsServiceServicer):
    def __init__(self, embedder: ClipEmbedder):
        self.embedder = embedder

    def GenerateTextEmbeddings(self, request, context):
        text = request.text
        vec = self.embedder.embed_text(text)
        return embeddings_pb2.GenerateTextEmbeddingsResponse(
            embedding_vector=vec.astype(float)
        )


def start_grpc_server(embedder: ClipEmbedder):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    service_names = (
        reflection.SERVICE_NAME,
        health.SERVICE_NAME,
        'embeddings.v1.EmbeddingsService',
    )
    reflection.enable_server_reflection(service_names, server)
    embeddings_pb2_grpc.add_EmbeddingsServiceServicer_to_server(
        EmbeddingsServiceServicer(embedder), server
    )
    health_pb2_grpc.add_HealthServicer_to_server(health.HealthServicer(), server)

    server.add_insecure_port("[::]:50051")
    thread = threading.Thread(target=server.start)
    thread.start()
    print("gRPC server running on port 50051")
    return server
