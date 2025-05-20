import asyncio

from app.worker.config import CONCURRENCY, DURABLE, PULL_SUBJ
from app.worker.nats_client import get_jetstream
from app.worker.processor import MessageProcessor


async def start_worker():
    js = await get_jetstream()
    processor = MessageProcessor(js)

    consumer = await js.pull_subscribe(
        subject=PULL_SUBJ,
        durable=DURABLE,
        stream=None,  # auto-resolved to our STREAM
    )

    print("Worker started")

    async def worker_loop():
        while True:
            try:
                msgs = await consumer.fetch(1, timeout=1.0)
            except asyncio.TimeoutError:
                continue
            await processor.process(msgs[0])

    tasks = [worker_loop() for _ in range(CONCURRENCY)]
    try:
        await asyncio.gather(*tasks)
    finally:
        await processor.close()
