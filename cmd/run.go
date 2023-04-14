package cmd

import (
	"os"

	"github.com/philippheuer/go-kafka-consumer/consumers/embeddinggenerator"
	"github.com/philippheuer/go-kafka-consumer/pkg/api"
	"github.com/philippheuer/go-kafka-consumer/pkg/consumer"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var consumers []api.MessageConsumer

func init() {
	rootCmd.AddCommand(runCmd)
	consumers = append(consumers, embeddinggenerator.New())
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "runs a consumer",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Usage()
			os.Exit(1)
		}
		id := args[0]

		// start consumer
		for _, c := range consumer.Consumers {
			if c.ID() == id {
				consumer.RunConsumer(c)
				os.Exit(0)
			}
		}
		log.Fatal().Msgf("consumer not found: %v", id)
	},
}
