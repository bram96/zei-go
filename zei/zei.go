package zei

import (
	"github.com/paypal/gatt"
)

var ledButtonService = gatt.MustParseUUID("c7e70020c84711e681758c89a55d403c")
var orientationService = gatt.MustParseUUID("c7e70010c84711e681758c89a55d403c")
var positionCharacteristic = gatt.MustParseUUID("c7e70012c84711e681758c89a55d403c")
var accelerometerCharacteristic = gatt.MustParseUUID("c7e70011c84711e681758c89a55d403c")
var ledCharacteristic = gatt.MustParseUUID("c7e70022c84711e681758c89a55d403c")
var pushButtonCharacteristic = gatt.MustParseUUID("c7e70021c84711e681758c89a55d403c")

type Status string

//ZEI connection statusses
const (
	ZEIStatusDisconnected Status = "Disconnected"
	ZEIStatusConnected    Status = "Connected"
)

type Position byte

var PositionName = map[Position]string{
	0: "TIP_UP",
	1: "DOWN_NORTH",
	2: "DOWN_WEST",
	3: "DOWN_SOUTH",
	4: "DOWN_EAST",
	5: "UP_NORTH",
	6: "UP_WEST",
	7: "UP_SOUTH",
	8: "UP_EAST",
	9: "TIP_DOWN",
}

type ZEI struct {
	position  Position
	p         gatt.Peripheral
	reconnect bool
	status    Status
	hook      Hook
}

//ChangePosition changes the current active position
func (z *ZEI) ChangePosition(p Position) {
	if z.position == p {
		return
	}
	z.position = p
	if z.hook != nil {
		z.hook.OnPositionChanged(p)
	}
}

//Position returns the current active position
func (z *ZEI) Position() Position {
	return z.position
}

//ConnectionStatus returns the current connection status
func (z *ZEI) ConnectionStatus() Status {
	return z.status
}

func (z *ZEI) Disconnect() {
	z.reconnect = false
	z.p.Device().CancelConnection(z.p)
}

func NewZEIConnection(h Hook) (*ZEI, error) {
	z := &ZEI{
		reconnect: true,
		hook:      h,
	}
	return z, z.connect()
}
