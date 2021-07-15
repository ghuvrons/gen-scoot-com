package packet

type ReportPacket struct {
	Header *HeaderPacket
	VCU    *VcuPacket
	GPS    *GpsPacket
}

type HeaderPacket struct {
	Prefix       string `len:"2"`
	Size         uint8  `unit:"Bytes"`
	Vin          uint32 `type:"uint32"`
	SendDatetime int64  `len:"7"`
}

type VcuPacket struct {
	FrameID     uint8   ``
	LogDatetime int64   `len:"7"`
	State       int8    `type:"int8"`
	Events      uint16  ``
	LogBuffered uint8   ``
	BatVoltage  float32 `len:"1" unit:"mVolt" factor:"18.0"`
	Uptime      float32 `type:"uint32" unit:"hour" factor:"0.000277"`
}

type GpsPacket struct {
	Active    bool    ``
	SatInUse  uint8   `unit:"Sat"`
	HDOP      float32 `len:"1" factor:"0.1"`
	VDOP      float32 `len:"1" factor:"0.1"`
	Speed     uint8   `unit:"Kph"`
	Heading   float32 `len:"1" unit:"Deg" factor:"2.0"`
	Longitude float32 `factor:"0.0000001"`
	Latitude  float32 `factor:"0.0000001"`
	Altitude  float32 `len:"2" unit:"m" factor:"0.1"`
}
