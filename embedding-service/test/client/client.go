package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/tech-inspire/api-contracts/api/gen/go/embeddings/v1"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	msg := &embeddingsv1.GeneratePostEmbeddingsEvent{
		PostId:   uuid.NewString(),
		ImageUrl: "https://fastly.picsum.photos/id/523/512/512.jpg?hmac=InuSwViS76D-vMXugy2GJxRUDXtbRXw7OewLEHABeB4",
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	subj := fmt.Sprintf("posts.%s.generate_embeddings", msg.PostId)
	ack, err := js.Publish(subj, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("sent to stream %s seq=%d", ack.Stream, ack.Sequence)
	// nc.Drain()

	subSubj := "posts.*.embeddings_updated"

	// Plain subscribe (no JetStream) because the worker already publishes to the stream.
	_, err = nc.Subscribe(subSubj, func(m *nats.Msg) {
		upd := new(embeddingsv1.PostEmbeddingsUpdatedEvent)
		fmt.Println("len", len(m.Data))
		if err = proto.Unmarshal(m.Data, upd); err != nil {
			log.Println("decode error:", err)
			return
		}
		log.Printf("post %s img[0:5]=%v  updated_at=%v",
			upd.PostId,
			upd.EmbeddingVector[:5],
			upd.UpdatedAt)

		if err = m.Ack(); err != nil {
			log.Println("ack error:", err)
			return
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s …", subSubj)
	<-context.TODO().Done()
}
