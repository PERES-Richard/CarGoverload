package main

import (
	. "multipleSearchingAggregator/entities"
	"multipleSearchingAggregator/tools"
	"multipleSearchingAggregator/controllers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"sync"
)

const NEW_WISH_TOPIC_READER_ID = 0
const SEARCH_RESULT_TOPIC_READER_ID = 1

const RAW_WISH_RESULT_TOPIC_WRITER_ID = 0

var readers = make([]*kafka.Reader, 4)
var wg sync.WaitGroup

func main() {
	// Setup readers & writers
	setUpKafka()

	wg.Add(4)

	go listenKafka(SEARCH_RESULT_TOPIC_READER_ID)
	go listenKafka(NEW_WISH_TOPIC_READER_ID)

	wg.Wait()
}

func setUpKafka() {
	setupKafkaReaders()
	setupKafkaWriters()
}

func setupKafkaWriters() {
	configWriter := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "raw-wish-result",
		ClientId:  "multiple-searching-aggregator",
	}
	tools.SetUpWriter(RAW_WISH_RESULT_TOPIC_WRITER_ID,configWriter)
}

func setupKafkaReaders() {
	configReader := tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "new-wish",
		ClientId:  "multiple-searching-aggregator",
	}
	readers[NEW_WISH_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)
	configReader = tools.KafkaConfig{
		BrokerUrl: os.Getenv("KAFKA"),
		Topic:     "search-result",
		ClientId:  "multiple-searching-aggregator",
	}
	readers[SEARCH_RESULT_TOPIC_READER_ID] = tools.GetUpKafkaReader(configReader)
}

func listenKafka(readerId int) {
	for {
		m, err := readers[readerId].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln("Error reader " + string(rune(readerId)) + ": ", err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		messageHandlers(readerId, m)
	}
	wg.Done()
}

func messageHandlers(readerId int, m kafka.Message) {
	switch readerId {
		case SEARCH_RESULT_TOPIC_READER_ID:
		{
			var parsedMessage SearchResultMessage
			err := json.Unmarshal(m.Value, &parsedMessage)
			if err != nil {
				log.Panic("Error unmarshalling availability result message:", err)
			}
			controllers.SearchResultHandler(parsedMessage)
		}
		case NEW_WISH_TOPIC_READER_ID:
		{
			var parsedMessage NewWishMessageResult
			err := json.Unmarshal(m.Value, &parsedMessage)
			if err != nil {
				log.Panic("Error unmarshalling availability result message:", err)
			}
			controllers.NewWishHandler(parsedMessage)
		}
	}
}
