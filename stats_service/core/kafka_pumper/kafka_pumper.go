package kafka_pumper

import (
	"context"
	"fmt"
	"log"
	"stats_service_core/db_utils"
	reaction_pb "stats_service_core/reaction"
	"stats_service_core/utils"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const groupID = "pumper"
const KafkaTopicViews = "views"
const KafkaTopicLikes = "likes"

var kafkaReaderViews *kafka.Reader
var kafkaReaderLikes *kafka.Reader

func StartPumping() {
	kafka_connection := fmt.Sprintf(
		"%s:%s",
		utils.GetenvSafe("KAFKA_HOST"),
		utils.GetenvSafe("KAFKA_PORT"),
	)

	kafkaReaderViews = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafka_connection},
		Topic:   KafkaTopicViews,
		GroupID: groupID,
	})

	kafkaReaderLikes = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafka_connection},
		Topic:   KafkaTopicLikes,
		GroupID: groupID,
	})

	go pumpViews()
	go pumpLikes()
}

func pumpViews() {
	for {
		msg, err := kafkaReaderViews.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read message: %v", err)
		}

		var view reaction_pb.View
		err = proto.Unmarshal(msg.Value, &view)
		if err != nil {
			log.Fatalf("failed to unmarshal: %v", err)
		}

		err = db_utils.AddView(view.PostId, view.AppraiserId)
		if err != nil {
			log.Fatalf("failed to add view: %v", err)
		}
	}
}

func pumpLikes() {
	for {
		msg, err := kafkaReaderLikes.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read message: %v", err)
		}

		var like reaction_pb.Like
		err = proto.Unmarshal(msg.Value, &like)
		if err != nil {
			log.Fatalf("failed to unmarshal: %v", err)
		}

		err = db_utils.AddLike(like.PostId, like.AppraiserId)
		if err != nil {
			log.Fatalf("failed to add like: %v", err)
		}
	}
}
