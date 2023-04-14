package consumer

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/philippheuer/go-kafka-consumer/consumers/embeddinggenerator"
	"github.com/philippheuer/go-kafka-consumer/pkg/api"
	"github.com/rs/zerolog/log"
)

var Consumers []api.MessageConsumer

func init() {
	Consumers = append(Consumers, embeddinggenerator.New())
}

func RunConsumer(consumer api.MessageConsumer) {
	// start the consumer in a goroutine
	sigDone := make(chan struct{})
	go func() {
		err := consumer.Run()
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to start consumer: %v", err)
		}
		close(sigDone) // signal that the goroutine has completed
	}()

	// wait for either a SIGTERM signal or the goroutine to complete
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		// SIGTERM received, shutdown
		consumer.Close()
	case <-sigDone:
		// goroutine completed, shutdown
		consumer.Close()
	}
}
