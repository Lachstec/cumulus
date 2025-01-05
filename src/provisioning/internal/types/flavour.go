package types

// Flavour represents the various Virtual Machine flavours
// that are available in our open stack cluster. The Flavour are currently:
// +----+-----------+-------+------+-----------+-------+-----------+
// | ID | Name      |   RAM | Disk | Ephemeral | VCPUs | Is Public |
// +----+-----------+-------+------+-----------+-------+-----------+
// | 1  | m1.tiny   |   512 |    1 |         0 |     1 | True      |
// | 2  | m1.small  |  2048 |   20 |         0 |     1 | True      |
// | 3  | m1.medium |  4096 |   40 |         0 |     2 | True      |
// | 4  | m1.large  |  8192 |   80 |         0 |     4 | True      |
// | 5  | m1.xlarge | 16384 |  160 |         0 |     8 | True      |
// +----+-----------+-------+------+-----------+-------+-----------+
type Flavour string

const (
	Tiny   Flavour = "1" //nolint:all
	Small  Flavour = "2" //nolint:all
	Medium Flavour = "3" //nolint:all
	Large  Flavour = "4" //nolint:all
	XLarge Flavour = "5" //nolint:all
)

func (f Flavour) Value() string {
	return string(f)
}

func (f Flavour) AvailableRam() int {
	var ram int

	switch f {
	case Tiny:
		ram = 512
	case Small:
		ram = 2048
	case Medium:
		ram = 4096
	case Large:
		ram = 8192
	case XLarge:
		ram = 16384
	}

	return ram
}
