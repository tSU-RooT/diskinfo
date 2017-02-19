//+build debug

package atasmart

const (
	DebugEnabled = true
)

var (
	examples []*SkDisk
)

func init() {
	// Good Case
	sda := &SkDisk{
		internal: nil,

		Data: &SkDiskData{
			Serial:   "FOOBAR VENDOR A42",
			Firmware: "TEA POD",
			Model:    "ABCDEF",
		},
		Attributes: []Attribute{
			{0x01, "raw-read-error-rate", 100, 100, 16, 0},
			{0x02, "throughput-performance", 133, 100, 54, 0},
			{0x03, "spin-up-time", 253, 100, 24, 0},
			{0x05, "reallocated-sector-count", 100, 100, 5, 0},
			{0x07, "seek-error-rate", 100, 100, 67, 0},
			{0x08, "seek-time-performance", 128, 100, 20, 0},
			{0x09, "power-on-hour", 100, 100, 0, 0},
			{0x0A, "spin-retry-count", 100, 100, 60, 0},
			{0x16, "attribute-22", 100, 100, 25, 0},
			{0xB4, "unused-reserved-blocks ", 100, 100, 98, 0},
			{0xC2, "temperature-celsius-2", 171, 130, 0, 0},
			{0xC4, "reallocated-event-count", 100, 100, 0, 0},
		},
		Path: "/dev/sda",
	}
	// Bad Case
	sdb := &SkDisk{
		internal: nil,

		Data: &SkDiskData{
			Serial:   "FOOBAR VENDOR A43",
			Firmware: "TEA POD",
			Model:    "GHIDEF",
		},
		Attributes: []Attribute{
			{0x01, "raw-read-error-rate", 100, 100, 16, 0},
			{0x02, "throughput-performance", 133, 100, 54, 0},
			{0x03, "spin-up-time", 253, 100, 24, 0},
			// BAD
			{0x05, "reallocated-sector-count", 2, 2, 5, 0},
			{0x07, "seek-error-rate", 100, 100, 67, 0},
			{0x08, "seek-time-performance", 128, 100, 20, 0},
			{0x09, "power-on-hour", 100, 100, 0, 0},
			{0x0A, "spin-retry-count", 100, 100, 60, 0},
			{0x16, "attribute-22", 100, 100, 25, 0},
			{0xB4, "unused-reserved-blocks ", 100, 100, 98, 0},
			{0xC2, "temperature-celsius-2", 171, 130, 0, 0},
			{0xC4, "reallocated-event-count", 100, 100, 0, 0},
		},
		Path: "/dev/sdb",
	}
	examples = make([]*SkDisk, 2)
	examples[0] = sda
	examples[1] = sdb
}

func DiskOpen(path string) (*SkDisk, error) {
	for i, _ := range examples {
		if examples[i].Path == path {
			return examples[i], nil
		}
	}
	return nil, InternalError
}

func (sd *SkDisk) GetSize() (uint64, error) {
	return 100 * 1024 * 1024 * 1024, nil
}

func (sd *SkDisk) IsSleepMode() (bool, error) {
	return true, nil
}

func (sd *SkDisk) ParseIdentify() error {
	// Do nothing
	return nil
}

func (sd *SkDisk) ReadSmartData() error {
	// Do nothing
	return nil
}

func (sd *SkDisk) GetSmartPowerCycle() (uint64, error) {
	return 0, nil
}

func (sd *SkDisk) GetSmartPowerOn() (uint64, error) {
	return 0, nil
}

func (sd *SkDisk) GetSmartTemperature() (uint64, error) {
	return 300000, nil
}

func (sd *SkDisk) ParseSmartAttr() error {
	// Do nothing
	return nil
}

func (sd *SkDisk) Free() {
	// Do nothing
}

func (sd *SkDisk) DiskDump() int {
	// Do nothing
	return 0
}
