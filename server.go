package main

import (
	"fmt"
	"log"

	"github.com/ghuvrons/gen-scoot-com/mqtt"
)

var closeChan = make(chan byte)

func Serve() {
	mq := &mqtt.Mqtt{}
	defer func() {
		mq.Disconnect()
		fmt.Println("Closed")
	}()

	if err := mq.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %s\n", err.Error())
	}

	subscribers := mqtt.Subscribers{
		"VCU/+/RPT": onReportMessage,
	}
	if err := mq.Subscribe(subscribers); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %s\n", err.Error())
	}
	for {
		isClosed := false
		select {
		case <-closeChan:
			isClosed = true
		}
		if isClosed {
			break
		}
	}
}

func Close() {
	closeChan <- 1
}
