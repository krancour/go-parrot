package ardrone3

const (
	AccessoryTypeSequoia int32 = 0
	AccessoryTypeFLIR    int32 = 1
)

// Accessory ...
// TODO: Document this
type Accessory interface {
	// ID returns the ID of the accessory.
	ID() uint8
	// Type returns the type of accessory.
	Type() int32
	// UID returns the accessory's unique identifier.
	UID() string
	// SoftwareVersion returns the software version of the accessory.
	SoftwareVersion() string
	// BatteryPercent returns the accessory's percentage of battery life
	// remaining.
	BatteryPercent() uint8
}

type accessory struct {
	id              uint8
	tipe            int32
	uid             string
	softwareVersion string
	batteryPercent  uint8
}

func (a *accessory) ID() uint8 {
	return a.id
}

func (a *accessory) Type() int32 {
	return a.tipe
}

func (a *accessory) UID() string {
	return a.uid
}

func (a *accessory) SoftwareVersion() string {
	return a.softwareVersion
}

func (a *accessory) BatteryPercent() uint8 {
	return a.batteryPercent
}
