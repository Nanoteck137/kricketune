package broker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EventData interface {
	GetEventType() string
}

type Event struct {
	Type string    `json:"type"`
	Data EventData `json:"data"`
}

// NOTE(patrik): Based on: https://gist.github.com/Ananto30/8af841f250e89c07e122e2a838698246
type Broker struct {
	notifier chan EventData

	newClients     chan chan EventData
	closingClients chan chan EventData
	clients        map[chan EventData]struct{}
}

func NewBroker() *Broker {
	return &Broker{
		notifier:       make(chan EventData, 1024),
		newClients:     make(chan chan EventData),
		closingClients: make(chan chan EventData),
		clients:        make(map[chan EventData]struct{}),
	}
}

func (broker *Broker) Listen() {
	for {
		select {
		case c := <-broker.newClients:
			broker.clients[c] = struct{}{}
		case c := <-broker.closingClients:
            if _, ok := broker.clients[c]; ok {
                delete(broker.clients, c)
                close(c)
            }
		case event := <-broker.notifier:
            for c := range broker.clients {
                select {
                case c <- event:
                default:
                    // Drop event for slow client instead of blocking
                }
            }
		}
	}
}

func (broker *Broker) Start() {
	go broker.Listen()
}

func (broker *Broker) EmitEvent(event EventData) {
	broker.notifier <- event
}

var _ (EventData) = (*ConnectedEvent)(nil)

type ConnectedEvent struct {
}

func (c ConnectedEvent) GetEventType() string {
	return "connected"
}

func (broker *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	rc := http.NewResponseController(w)

	eventChan := make(chan EventData, 16)
	broker.newClients <- eventChan
	defer func() {
		broker.closingClients <- eventChan
	}()

	sendEvent := func(eventData EventData) {
		fmt.Fprintf(w, "event: %s\n", eventData.GetEventType())
		fmt.Fprintf(w, "data: ")

		// event := Event{
		// 	Type: eventData.GetEventType(),
		// 	Data: eventData,
		// }

		encode := json.NewEncoder(w)
		encode.Encode(eventData)

		fmt.Fprintf(w, "\n\n")
		rc.Flush()
	}

	sendEvent(ConnectedEvent{})

	for {
		select {
		case <-r.Context().Done():
			return

		case event := <-eventChan:
			sendEvent(event)
		}
	}
}
