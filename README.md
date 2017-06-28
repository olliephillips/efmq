# Ethernet Frames Message Queue (EFMQ)

EFMQ provides an MQTT pub/sub style abstraction for Ethernet Frame messaging.  

EFMQ works a bit like MQTT, but only over your Local Area Network (LAN). Since message traffic is effectively broadcast peer-to-peer, no broker is required, either local or remotely hosted. Since messages are routed locally, within your LAN, it's more secure than using one of the few free remote MQTT servers to bounce your traffic off.

Messaging can be two-way. Each node can operate as publisher and subscriber at the same time, which opens the door for control as well as monitoring applications.

This package leans heavilly on @mdlayher's [raw](https://github.com/mdlayher/raw) and [ethernet](https://github.com/mdlayher/ethernet) packages, which do almost all the heavy lifting.

## Usage
Basic publisher and subscriber examples are provided below. Nodes can publish and subscribe to multiple topics.

The API follows MQTT's API loosely.

```
// Create connection
mq, _ := efmq.NewEFMQ(networkInterface string)

// Publish
mq.Publish(topic string, payload string)

// Subscribe
mq.Subscribe(topic string)

// Unsubscribe
mq.Unsubscribe(topic string)

// List subscriptions
mq.Subscriptions()

// Start listening
mq.Listen()

// Message channel
mq.Message

// Message 
Message struct {
	Topic string
	Payload string
}
```

### Publisher example
The code below will publish data to the `fermenter` topic every second. `en1` is the network interface on Mac (my Mac at least). On a Raspberry Pi it might be `wlan0`. Use `netstat -i` to discover.

```
mq, err := efmq.NewEFMQ("en1") 
if err != nil {
	log.Fatal(err)
}
t := time.NewTicker(1 * time.Second)
for range t.C {
	if err := mq.Publish("fermenter", "20.5"); err != nil {
		log.Fatalln(err)
	}
}
```

### Subscriber example
The code below sets up a subcription to the `fermenter` topic and then listens for messages. Messages are received on a channel.

```
mq, err := efmq.NewEFMQ("wlan0")
if err != nil {
	log.Fatal(err)
}
mq.Subscribe("fermenter")
mq.listen()
for msg := range mq.Message {
	fmt.Println("topic:", msg.Topic)
	fmt.Println("message:", msg.Payload)
}
```

## Todo
- Better test coverage
- Check message does not exceed frame byte data limit (1500 bytes?)