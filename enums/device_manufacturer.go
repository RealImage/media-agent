package enums

type DeviceManufacturer string

const (
	DeviceManufacturerDolby DeviceManufacturer = "Dolby"
	DeviceManufacturerQube  DeviceManufacturer = "Qube Cinema"
)

func (d DeviceManufacturer) String() string {
	return string(d)
}
