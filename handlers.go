package main

import (
	"bytes"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ghuvrons/gen-scoot-com/decoder"
	"github.com/ghuvrons/gen-scoot-com/packet"
)

var reportHanlder func(vin []byte, data *packet.ReportPacket)

func onReportMessage(client mqtt.Client, msg mqtt.Message) {
	// hexString := strings.ToUpper(hex.EncodeToString())
	// binary.LittleEndian.Uint64()

	buf := bytes.NewReader(msg.Payload())
	p := &packet.ReportPacket{}
	decoder.Decode(buf, p)
	fmt.Println(p, p.Header)
	fmt.Println(p.VCU, p.VCU.BatVoltage, p.GPS)

	if reportHanlder != nil {
		reportHanlder([]byte("AADC"), &packet.ReportPacket{})
	}
}

func OnReport(f func(vin []byte, data *packet.ReportPacket)) {
	reportHanlder = f
}
