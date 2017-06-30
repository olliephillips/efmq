package efmq_test

import (
	"fmt"
	"log"
	"time"

	"github.com/olliephillips/efmq"
)

// This example starts publishing to a topic every second via the
// network interface "wlan0".
func Example_publisher() {
	mq, err := efmq.NewEFMQ("wlan0")
	if err != nil {
		log.Fatal(err)
	}
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		if err := mq.Publish("fermenter", "20.5"); err != nil {
			log.Fatalln(err)
		}
	}
}

// This example sets up a subscription to a topic and starts listening
// for messages from a device on the same network.
func Example_subscriber() {
	mq, err := efmq.NewEFMQ("wlan0")
	if err != nil {
		log.Fatal(err)
	}
	mq.Subscribe("fermenter")
	mq.Listen()
	for msg := range mq.Message {
		fmt.Println("topic:", msg.Topic)
		fmt.Println("message:", msg.Payload)
	}
	// Output: fermenter
	// 20.5
	// fermenter
	// 20.5
	// fermenter
	// 20.5
}
