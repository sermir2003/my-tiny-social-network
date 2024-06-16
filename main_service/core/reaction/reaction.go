package reaction

import (
	"context"
	"fmt"
	"log"
	"main_service_core/utils"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const KafkaTopicViews = "views"
const KafkaTopicLikes = "likes"

var kafkaWriterViews *kafka.Writer
var kafkaWriterLikes *kafka.Writer

func Connect() error {
	kafka_connection := fmt.Sprintf(
		"%s:%s",
		utils.GetenvSafe("KAFKA_HOST"),
		utils.GetenvSafe("KAFKA_PORT"),
	)

	kafkaWriterViews = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafka_connection},
		Topic:   KafkaTopicViews,
	})

	kafkaWriterLikes = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafka_connection},
		Topic:   KafkaTopicLikes,
	})

	return nil
}

func ReportView(post_id uint64, appraiser_id uint64) error {
	view := &View{
		PostId:      post_id,
		AppraiserId: appraiser_id,
	}

	msg, err := proto.Marshal(view)
	if err != nil {
		log.Fatalf("failed to marshal: %v\n", err)
	}

	err = kafkaWriterViews.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
	if err != nil {
		log.Fatalf("failed to write messages: %v\n", err)
	}

	return nil
}

func ReportLike(post_id uint64, appraiser_id uint64) error {
	like := &Like{
		PostId:      post_id,
		AppraiserId: appraiser_id,
	}

	msg, err := proto.Marshal(like)
	if err != nil {
		log.Fatalf("failed to marshal: %v\n", err)
	}

	err = kafkaWriterLikes.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	})
	if err != nil {
		log.Fatalf("failed to write messages: %v\n", err)
	}

	return nil
}
