package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"math"

	"github.com/nats-io/nats.go"
	"github.com/tech-inspire/api-contracts/api/gen/go/embeddings/v1"
)

func isNormalized(vec []float32) bool {
	var sum float64
	for _, v := range vec {
		sum += float64(v * v)
	}
	norm := math.Sqrt(sum)
	fmt.Println("norm:", norm)
	return math.Abs(norm-1.0) < 1e-4 // tolerance for float error
}

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
		PostId:   "0196fd33-810e-7b48-8633-bab816ad39b6",
		ImageUrl: "https://inspire-test.nyc3.cdn.digitaloceanspaces.com/images/post_0196f90c-7543-7767-be91-e76f83d8fb8f"}

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

	//Plainsubscribe (no JetStream) because the worker already publishes to the stream.
	//_, err = nc.Subscribe(subSubj, func(m *nats.Msg) {
	//	upd := new(embeddingsv1.PostEmbeddingsUpdatedEvent)
	//	fmt.Println("len", len(m.Data))
	//	if err = proto.Unmarshal(m.Data, upd); err != nil {
	//		log.Println("decode error:", err)
	//		return
	//	}
	//	log.Printf("post %s img[0:5]=%v  updated_at=%v",
	//		upd.PostId,
	//		upd.EmbeddingVector[:5],
	//		upd.UpdatedAt)
	//
	//	isNormalized := isNormalized(upd.EmbeddingVector)
	//	fmt.Println("isNormalized:", isNormalized)
	//
	//	if err = m.Ack(); err != nil {
	//		log.Println("ack error:", err)
	//		return
	//	}
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Printf("listening on %s …", subSubj)
	<-context.TODO().Done()
}
