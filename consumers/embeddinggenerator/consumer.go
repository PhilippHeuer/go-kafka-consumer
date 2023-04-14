package embeddinggenerator

import (
	"context"
	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/philippheuer/go-kafka-consumer/config"
	"github.com/philippheuer/go-kafka-consumer/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func New() *MessageConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &MessageConsumer{
		Ctx:         ctx,
		Cancel:      cancel,
		KafkaConfig: config.LoadKafkaConfig(),
		Topic:       "test",
	}
}

type MessageConsumer struct {
	Ctx         context.Context
	Cancel      context.CancelFunc
	KafkaConfig *config.KafkaConfig
	Topic       string
}

func (m *MessageConsumer) ID() string {
	return "embedding-generator"
}

func (m *MessageConsumer) Run() error {
	// make a new reader that consumes from topic-A
	kafkaReader := kafka.NewReader(api.ReaderConfig(m.KafkaConfig, m.Topic, "group-id"))

	// poll messages
	log.Info().Str("topic", m.Topic).Msg("registered kafka subscription")
	for {
		select {
		case <-m.Ctx.Done():
			log.Info().Str("topic", m.Topic).Msg("released kafka subscription")
			return nil
		default:
			message, err := kafkaReader.ReadMessage(m.Ctx)
			if err != nil {
				log.Err(err).Str("topic", m.Topic).Msg("failed to fetch message")
				continue
			}

			// process message
			event := cloudevents.NewEvent()
			err = json.Unmarshal(message.Value, &event)
			if err != nil {
				log.Err(err).Str("topic", m.Topic).Msg("failed to unmarshal kafka message")
				continue
			}
			log.Info().Msgf("event: %v", string(event.Data()))

			// complete
			log.Info().Str("topic", m.Topic).Str("event_id", event.ID()).Msgf("event has been processed successfully.")
		}
	}
}

func (m *MessageConsumer) Close() {
	log.Info().Str("topic", m.Topic).Msg("cancelling kafka subscription")
	m.Cancel()
}
