package zei

import (
	"fmt"

	"github.com/paypal/gatt"
)

func (z *ZEI) connect() error {
	d, err := gatt.NewDevice(
		gatt.LnxMaxConnections(1),
		gatt.LnxDeviceID(-1, true),
	)
	if err != nil {
		return fmt.Errorf("Failed to open device, err: %s\n", err)
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(z.onPeriphConnected),
		gatt.PeripheralDisconnected(z.onPeriphDisconnected),
	)

	d.Init(onStateChanged)
	return nil
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if a.LocalName != "Timeular ZEI" {
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("")

	p.Device().Connect(p)
}

func (z *ZEI) onPeriphConnected(p gatt.Peripheral, err error) {
	z.status = ZEIStatusConnected
	fmt.Println("Connected")
	z.p = p

	// Discovery services
	ss, err := p.DiscoverServices([]gatt.UUID{orientationService})
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return
	}

	for _, s := range ss {
		if !s.UUID().Equal(orientationService) && !s.UUID().Equal(ledButtonService) {
			continue
		}

		// Discovery characteristics
		cs, err := p.DiscoverCharacteristics(nil, s)
		if err != nil {
			fmt.Printf("Failed to discover characteristics, err: %s\n", err)
			continue
		}

		for _, c := range cs {
			if c.UUID().Equal(positionCharacteristic) {

				// Read the characteristic, if possible.
				if (c.Properties() & gatt.CharRead) != 0 {
					b, err := p.ReadCharacteristic(c)
					if err != nil {
						fmt.Printf("Failed to read characteristic, err: %s\n", err)
						continue
					}

					z.ChangePosition(Position(b[0]))

				}
				//Needed for cccd
				_, err := p.DiscoverDescriptors(nil, c)
				if err != nil {
					fmt.Printf("Failed to discover descriptors, err: %s\n", err)
					continue
				}

				f := func(c *gatt.Characteristic, b []byte, err error) {
					z.ChangePosition(Position(b[0]))
				}
				if err := p.SetIndicateValue(c, f); err != nil {
					fmt.Printf("Failed to subscribe characteristic, err: %s\n", err)
				}
			} else if c.UUID().Equal(ledCharacteristic) {
				p.WriteCharacteristic(c, []byte{1}, false)
			}
		}
	}
}

func (z *ZEI) onPeriphDisconnected(p gatt.Peripheral, err error) {
	z.status = ZEIStatusDisconnected
	if z.reconnect {
		fmt.Println("Disconnected, try to reconnect...")
		p.Device().Scan([]gatt.UUID{}, false)
	}
}
