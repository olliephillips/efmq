// Package efmq provides basic MQTT like functionality for message
// publishing and subscriptions within a local area network
package efmq

import (
	"encoding/json"
	"errors"
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

// EFQM represents a connection
type EFMQ struct {
	netInterface *net.Interface
	connection   *net.PacketConn
	subscription []string
	listening    bool
	Message      chan Message
}

type Message struct {
	Topic   string `json:"tpc"`
	Payload string `json:"pyld`
}

const etherType = 0xcccc

// NewEFMQ is a factory functionm to create a value of EFMQ type
func NewEFMQ(networkInterface string) (*EFMQ, error) {
	mq := new(EFMQ)
	mq.Message = make(chan Message)
	// set network interface
	ni, err := net.InterfaceByName(networkInterface)
	if err != nil {
		return mq, errors.New("NewEFMQ: could not detect interface " + networkInterface)
	}
	// create connection/listener
	conn, err := connect(ni)
	if err != nil {
		return mq, err
	}
	// store in struct
	mq.netInterface = ni
	mq.connection = &conn
	return mq, nil
}

// connect opens network inteface to create connection for listening
func connect(ni *net.Interface) (net.PacketConn, error) {
	var conn net.PacketConn
	conn, err := raw.ListenPacket(ni, etherType)
	if err != nil {
		return conn, err
	}
	return conn, nil
}

// Subscribe takes a new subscription and stores it to slice
func (mq *EFMQ) Subscribe(topic string) {
	// add topic to subscriptions and start listener
	mq.subscription = append(mq.subscription, topic)
}

// Unsubscribe removes subscription from slice store
func (mq *EFMQ) Unsubscribe(topic string) error {
	// remove topic from subscriptions
	for i, v := range mq.subscription {
		if v == topic {
			mq.subscription = append(mq.subscription[:i], mq.subscription[i+1:]...)
		}
	}
	return nil
}

// Publish broadcasts a message on the network which comprises topic
// and payload
func (mq *EFMQ) Publish(topic string, payload string) error {
	// build a JSON object
	message := Message{
		Topic:   topic,
		Payload: payload,
	}
	// marshal to byte slice of JSON
	content, err := json.Marshal(&message)
	if err != nil {
		return errors.New("Publish: failed to marshal JSON")
	}
	// pass to despatcher
	if err := mq.despatcher(content); err != nil {
		return err
	}
	return nil
}

// despatcher handles the tranmission of message over ethernet frames
func (mq *EFMQ) despatcher(content []byte) error {
	// configure frame
	f := &ethernet.Frame{
		Destination: ethernet.Broadcast,
		Source:      mq.netInterface.HardwareAddr,
		EtherType:   etherType,
		Payload:     content,
	}
	// required for linux as mdlayher ethecho
	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}
	// prepare
	binary, err := f.MarshalBinary()
	if err != nil {
		return errors.New("despatcher: failed to marshal ethernet frame")
	}
	// send
	conn := *mq.connection
	if _, err := conn.WriteTo(binary, addr); err != nil {
		return errors.New("despatcher: failed to send message")
	}
	return nil
}

func (mq *EFMQ) Subscriptions() []string {
	return mq.subscription
}

// Listen announces the subscriptions to which we are subscribed
// and then starts listener func in go routine
func (mq *EFMQ) Listen() {
	var subs string
	subsLen := len(mq.subscription)
	for i, v := range mq.subscription {
		subs += v
		if i < subsLen-1 {
			subs += ", "
		} else {
			subs += "."
		}
	}
	// listen & log
	log.Println("Subscribed to topic(s):", subs, "Now listening...")
	go mq.listener()
}

// listener filters messages before presenting to client using topic
func (mq *EFMQ) listener() {
	var f ethernet.Frame
	var conn net.PacketConn
	var subs []string
	conn = *mq.connection
	subs = mq.subscription
	b := make([]byte, mq.netInterface.MTU)
	// handle messages indefinitely
	for {
		n, _, err := conn.ReadFrom(b)
		if err != nil {
			log.Printf("listener: failed to receive message: %v", err)
		}
		if err := (&f).UnmarshalBinary(b[:n]); err != nil {
			log.Printf("listener: failed to unmarshal ethernet frame: %v", err)
		}
		// f.Payload could be padded with zeros, need to deal before unmarshal
		var payload []byte
		for _, v := range f.Payload {
			if v != 0 {
				payload = append(payload, v)
			}
		}
		// unmarshal JSON
		message := new(Message)
		err = json.Unmarshal(payload, message)
		if err != nil {
			log.Println(err)
		}

		for _, v := range subs {
			if message.Topic == v {
				// put message on channel if matches a subscription
				mq.Message <- *message
			}
		}
	}
}
