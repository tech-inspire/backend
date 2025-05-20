import asyncio, io, os, msgpack, aiohttp, nats
import datetime
from typing import Optional
import traceback
from google.protobuf.timestamp_pb2 import Timestamp
from datetime import datetime, timezone
from embeddings.v1 import events_pb2
from nats.js.api import StreamConfig, RetentionPolicy, StorageType
from PIL import Image
from app.services.clip_embedder import embed_image, embed_text

NATS_URL   = os.getenv("NATS_URL", "nats://localhost:4222")
STREAM     = "POSTS"
PULL_SUBJ  = "posts.*.generate_embeddings"
PUSH_SUBJ  = "posts.{post_id}.embeddings_updated"
DURABLE    = os.getenv("WORKER_DURABLE", "embedding_worker")
CONCURRENCY = int(os.getenv("WORKER_CONCURRENCY", "4"))

_http_session: Optional[aiohttp.ClientSession] = None

async def _ensure_stream(js):
    if STREAM not in (await js.streams_info()):
        cfg = StreamConfig(
            name=STREAM,
            subjects=[PULL_SUBJ, "posts.*.embeddings_updated"],
            retention=RetentionPolicy.LIMITS,
            storage=StorageType.FILE,
        )
        await js.add_stream(cfg)

async def _get_http():
    global _http_session
    if _http_session is None:
        _http_session = aiohttp.ClientSession()
    return _http_session

async def _download_image(url: str) -> Image.Image:
    session = await _get_http()
    async with session.get(url) as resp:
        resp.raise_for_status()
        data = await resp.read()
    return Image.open(io.BytesIO(data)).convert("RGB")

async def _process_msg(js, msg):
    try:
        request = events_pb2.GeneratePostEmbeddingsEvent()
        request.ParseFromString(msg.data)
        print(request)

        post_id      = request.post_id
        image_url    = request.image_url

        img = await _download_image(image_url)
        img_vec  = embed_image(img).tolist()
        # desc_vec = embed_text(description).tolist()

        now = datetime.now(timezone.utc)
        ts = Timestamp()
        ts.FromDatetime(now)

        event = events_pb2.PostEmbeddingsUpdatedEvent(
            post_id=request.post_id,
            updated_at=ts,
            embedding_vector=img_vec,
        )
        event_bytes = event.SerializeToString()
        await js.publish(
            subject=PUSH_SUBJ.format(post_id=post_id),
            payload=event_bytes,
        )
        await msg.ack()
        print(f"✔ processed post {post_id}")

    except Exception as exc:
        print("✖ error:", exc)
        traceback.print_exc()
        await msg.nak()

async def main():
    nc = await nats.connect(servers=[NATS_URL])
    js = nc.jetstream()
    await _ensure_stream(js)

    consumer = await js.pull_subscribe(
        subject=PULL_SUBJ,
        durable=DURABLE,
        stream=STREAM,
    )

    sem = asyncio.Semaphore(CONCURRENCY)

    async def worker_loop():
        while True:
            try:
                msgs = await consumer.fetch(1, timeout=1.0)
            except asyncio.TimeoutError:
                continue
            async with sem:  # limit concurrent tasks
                await _process_msg(js, msgs[0])

    # spawn concurrency workers
    await asyncio.gather(*(worker_loop() for _ in range(CONCURRENCY)))

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass
