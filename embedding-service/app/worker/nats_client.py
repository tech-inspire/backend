import nats
from nats.js.api import RetentionPolicy, StorageType, StreamConfig

from app.worker.config import NATS_URL, PULL_SUBJ, PUSH_SUBJ, STREAM


async def get_jetstream():
    nc = await nats.connect(servers=[NATS_URL])
    js = nc.jetstream()
    await ensure_stream(js)
    return js


async def ensure_stream(js):
    streams = await js.streams_info()
    if STREAM not in streams:
        cfg = StreamConfig(
            name=STREAM,
            subjects=[PULL_SUBJ, PUSH_SUBJ],
            retention=RetentionPolicy.LIMITS,
            storage=StorageType.FILE,
        )
        await js.add_stream(cfg)
