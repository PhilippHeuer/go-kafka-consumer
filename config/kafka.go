package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	BootstrapServers string
	GroupID          string
	AutoOffsetReset  string
}

func LoadKafkaConfig() *KafkaConfig {
	viper.SetConfigFile(".env")
	viper.SetDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	viper.SetDefault("KAFKA_GROUP_ID", "go-kafka-consumer")
	viper.SetDefault("KAFKA_AUTO_OFFSET_RESET", "latest")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to read config file: %v", err)
	}

	return &KafkaConfig{
		BootstrapServers: viper.GetString("KAFKA_BOOTSTRAP_SERVERS"),
		GroupID:          viper.GetString("KAFKA_GROUP_ID"),
		AutoOffsetReset:  viper.GetString("KAFKA_AUTO_OFFSET_RESET"),
	}
}
