import { AckPolicy, connect, NatsConnection } from "nats";
import {
  PostDeletedEventSchema,
  type PostDeletedEvent,
} from "inspire-api-contracts/api/gen/ts/posts/v1/events_pb";
import { fromBinary } from "@bufbuild/protobuf";
import { deletePostLikesData } from "../db/likesRepository";

let natsConnection: NatsConnection | null = null;

export async function startPostsDeletedSubscriber(streamName: string) {
  const nc = await connect({ servers: "nats://nats:4222" });
  natsConnection = nc;

  const js = nc.jetstream();
  const manager = await js.jetstreamManager();

  const subject = "posts.*.deleted";

  const consumerInfo = await manager.consumers.add(streamName, {
    durable_name: "likes-posts-deleted-subscriber",
    ack_policy: AckPolicy.Explicit,
    filter_subject: subject,
  });

  const consumer = js.consumers.getPullConsumerFor(consumerInfo);

  const messages = await consumer.consume({});

  console.log(`Subscribed to '${subject}'`);

  const processMessage = async (event: PostDeletedEvent) => {
    if (!event.post) {
      throw new Error("event does not have 'post' field");
    }
    await deletePostLikesData(event.post?.postId);
  };

  for await (const msg of messages) {
    try {
      console.log("Processing event", msg.subject);

      const event = fromBinary(PostDeletedEventSchema, msg.data);
      await processMessage(event);

      msg.ack();
    } catch (err) {
      console.error(
        `[${msg.seq} ${msg.subject}] Error decoding or processing:`,
        err,
      );
      // No ack = message will be redelivered depending on config
    }
  }
}

export async function stopPostsDeletedSubscriber() {
  if (natsConnection) {
    try {
      console.log("Draining NATS connection...");
      await natsConnection.drain();
      console.log("NATS connection drained.");
    } catch (err) {
      console.error("Error draining NATS:", err);
    }
  }
}
