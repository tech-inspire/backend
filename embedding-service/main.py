import asyncio
import threading

from app.grpc import embedding
from app.services import clip_embedder
from app.workers.embedding_consumer import start_worker


def run_worker_loop():
    try:
        asyncio.run(start_worker())
    except KeyboardInterrupt:
        pass


def main():
    # print("Preloading captioner...")
    # captioner.preload()
    print("Preloading clip_embedder...")
    clip_embedder.preload()

    print("Starting gRPC server...")
    grpc_server = embedding.start_grpc_server()

    print("Starting NATS Jetstream worker...")
    threading.Thread(target=run_worker_loop, daemon=True).start()

    try:
        grpc_server.wait_for_termination()
    except KeyboardInterrupt:
        grpc_server.stop(0)


if __name__ == "__main__":
    main()
