package initialize

import (
	"os"
	"strings"

	"github.com/bignyap/go-utilities/pubsub"
)

func LoadPubSub() (pubsub.PubSubClient, error) {

	pubCfg := pubsub.Config{
		Type:      os.Getenv("PUBSUB_TYPE"),
		Enabled:   os.Getenv("PUBSUB_ENABLED") == "true",
		Namespace: os.Getenv("PUBSUB_NAMESPACE"),
	}

	switch pubCfg.Type {
	case "redis":
		pubCfg.Redis = &pubsub.RedisConfig{
			URL:      os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
		}
	case "kafka":
		pubCfg.Kafka = &pubsub.KafkaConfig{
			Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			GroupID: os.Getenv("KAFKA_GROUP_ID"),
			Topic:   os.Getenv("KAFKA_TOPIC"),
		}
	case "rabbitmq":
		pubCfg.RabbitMQ = &pubsub.RabbitMQConfig{
			URL:       os.Getenv("RABBITMQ_URL"),
			QueueName: os.Getenv("RABBITMQ_QUEUE"),
		}
	}

	return pubsub.NewPubSub(pubCfg)
}
