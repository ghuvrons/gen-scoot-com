package main

import (
	"fmt"
	"time"

	"github.com/ghuvrons/gen-scoot-com/packet"
)

func main() {
	/** test */
	/** run */

	go func() {
		time.Sleep(10 * time.Second)
		Close()
	}()
	OnReport(func(vin []byte, data *packet.ReportPacket) {
		fmt.Println("[REPORT] ", vin, &data)
	})
	Serve()
}
