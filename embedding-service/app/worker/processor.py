import traceback
from datetime import datetime, timezone

from embeddings.v1 import events_pb2
from google.protobuf.timestamp_pb2 import Timestamp

from app.services.clip_embedder import embed_image
from app.worker.downloader import ImageDownloader


class MessageProcessor:
    def __init__(self, js_client):
        self.js = js_client
        self.downloader = ImageDownloader()

    async def process(self, msg) -> None:
        try:
            req = events_pb2.GeneratePostEmbeddingsEvent()
            req.ParseFromString(msg.data)

            img = await self.downloader.download(req.image_url)
            vector = embed_image(img).tolist()

            now = datetime.now(timezone.utc)
            ts = Timestamp()
            ts.FromDatetime(now)

            out = events_pb2.PostEmbeddingsUpdatedEvent(
                post_id=req.post_id,
                updated_at=ts,
                embedding_vector=vector,
            )
            payload = out.SerializeToString()

            subj = f"posts.{req.post_id}.embeddings_updated"
            await self.js.publish(subject=subj, payload=payload)
            await msg.ack()
            print(f"Processed post {req.post_id}")

        except Exception as e:
            print("Error processing message:", e)
            traceback.print_exc()
            await msg.nak()

    async def close(self):
        await self.downloader.close()
