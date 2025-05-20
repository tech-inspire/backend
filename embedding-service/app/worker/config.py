import os

NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")
STREAM = "POSTS"
PULL_SUBJ = "posts.*.generate_embeddings"

PUSH_SUBJ_TEMPLATE = "posts.{post_id}.embeddings_updated"
PUSH_SUBJ = PUSH_SUBJ_TEMPLATE.format(post_id="*")

DURABLE = os.getenv("WORKER_DURABLE", "embedding_worker")
CONCURRENCY = int(os.getenv("WORKER_CONCURRENCY", "4"))
