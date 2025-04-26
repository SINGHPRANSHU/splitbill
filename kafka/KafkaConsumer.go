package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeMessages(broker, groupID, topic string) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
	})
	log.Println("Starting...", consumer)

	if err != nil {
		return fmt.Errorf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Subscribe to the topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	log.Println("Error Starting...", err)

	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %v", err)
	}

	log.Println("Kafka consumer is listening...")
	

	// Consume messages in an infinite loop
	for {
		msg, err := consumer.ReadMessage(-1) // Block until a message is received
		if err != nil {
			// Log consumer errors (e.g., timeout, broker issues)
			log.Printf("Consumer error: %v\n", err)
			continue
		}

		//iif messgae on the topic 
		// handler to call for  create user 
		// Successfully received a message
		log.Printf("Received message: %s from topic: %s\n", string(msg.Value), *msg.TopicPartition.Topic)
	}
}
