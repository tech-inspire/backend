import traceback
from datetime import datetime, timezone

from embeddings.v1 import events_pb2
from google.protobuf.timestamp_pb2 import Timestamp

from app.services.clip_embedder import embed_image
from app.worker.downloader import download_image


async def process_message(js, msg):
    try:
        # Parse incoming request
        req = events_pb2.GeneratePostEmbeddingsEvent()
        req.ParseFromString(msg.data)

        # Download image and compute embedding
        img = await download_image(req.image_url)
        vector = embed_image(img).tolist()

        # Build PostEmbeddingsUpdatedEvent
        now = datetime.now(timezone.utc)
        ts = Timestamp()
        ts.FromDatetime(now)
        out = events_pb2.PostEmbeddingsUpdatedEvent(
            post_id=req.post_id,
            updated_at=ts,
            embedding_vector=vector,
        )
        payload = out.SerializeToString()

        # Publish and ack
        subj = f"posts.{req.post_id}.embeddings_updated"
        await js.publish(subject=subj, payload=payload)
        await msg.ack()
        print(f"Processed post {req.post_id}")
    except Exception as e:
        print("Error processing message:", e)
        traceback.print_exc()
        await msg.nak()
