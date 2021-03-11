package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"offersCreator/tools"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"

	controller "offersCreator/controllers"
	. "offersCreator/entities"
)

const WISH_REQUESTED_TOPIC_READER_ID = 0
const RAW_WISH_RESULT_TOPIC_READER_ID = 1
const WISH_RESULT_TOPIC_WRITER_ID = 0

var readers = make([]*kafka.Reader, 2)
var wg sync.WaitGroup

func setUpKafka() {
	setupKafkaReaders()
	setupKafkaWriters()
}

func setupKafkaWriters() {
	configWriter := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "wish-result",
		ClientId:  "offersCreator",
	}
	tools.SetUpWriter(WISH_RESULT_TOPIC_WRITER_ID, configWriter)
}

func setupKafkaReaders() {
	// SetUpKafka
	configReader := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "wish-requested",
		ClientId:  "offersCreator",
	}
	readers[WISH_REQUESTED_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)

	configReader = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "raw-wish-result",
		ClientId:  "offersCreator",
	}
	readers[RAW_WISH_RESULT_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)
}

func listenKafka(readerId int) {
	for {
		m, err := readers[readerId].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		go messageHandlers(readerId, m)
	}
	wg.Done()
}

func messageHandlers(readerId int, m kafka.Message) {
	switch readerId {
	case WISH_REQUESTED_TOPIC_READER_ID:
		{
			var parsedMessage InitialWishRequest
			err := json.Unmarshal(m.Value, &parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling book validation message:", err)
			}
			controller.WishRequestedHandler(parsedMessage)
		}
	case RAW_WISH_RESULT_TOPIC_READER_ID:
		{
			var parsedMessage WishWithPossibilities
			err := json.Unmarshal(m.Value, &parsedMessage)
			if err != nil {
				log.Panic("Error unmarshaling offer possibilities search message:", err)
			}
			controller.RawWishHandler(&parsedMessage, WISH_RESULT_TOPIC_WRITER_ID)
		}
	}
}

func main() {
	// Setup readers & writers
	setUpKafka()

	wg.Add(2)

	go listenKafka(WISH_REQUESTED_TOPIC_READER_ID)
	go listenKafka(RAW_WISH_RESULT_TOPIC_READER_ID)

	wg.Wait()
}
