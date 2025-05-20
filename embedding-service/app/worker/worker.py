import asyncio

from app.worker.config import CONCURRENCY, DURABLE, PULL_SUBJ
from app.worker.nats_client import get_jetstream
from app.worker.processor import process_message


async def start_worker():
    js = await get_jetstream()
    consumer = await js.pull_subscribe(
        subject=PULL_SUBJ,
        durable=DURABLE,
        stream=None,  # auto-resolved to our STREAM
    )

    sem = asyncio.Semaphore(CONCURRENCY)
    print("Worker started")

    async def worker_loop():
        while True:
            try:
                msgs = await consumer.fetch(1, timeout=1.0)
            except asyncio.TimeoutError:
                continue
            async with sem:
                await process_message(js, msgs[0])

    await asyncio.gather(*(worker_loop() for _ in range(CONCURRENCY)))
