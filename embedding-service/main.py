import asyncio
import threading

from app.grpc import embedding
from app.services import clip_embedder
from app.services.clip_embedder import ClipEmbedder
from app.worker.worker import start_worker


def run_worker_loop(embedder: ClipEmbedder):
    try:
        asyncio.run(start_worker(embedder))
    except KeyboardInterrupt:
        pass


def main():
    embedder = clip_embedder.ClipEmbedder()

    print("Starting gRPC server...")
    grpc_server = embedding.start_grpc_server(embedder)

    print("Starting NATS Jetstream worker...")
    threading.Thread(
        target=run_worker_loop, args=(embedder,), daemon=True
    ).start()

    try:
        grpc_server.wait_for_termination()
    except KeyboardInterrupt:
        grpc_server.stop(0)


if __name__ == "__main__":
    main()
