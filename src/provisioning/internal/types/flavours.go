package types

// Flavours represents the various Virtual Machine flavours
// that are available in our open stack cluster. The Flavours are currently:
// +----+-----------+-------+------+-----------+-------+-----------+
// | ID | Name      |   RAM | Disk | Ephemeral | VCPUs | Is Public |
// +----+-----------+-------+------+-----------+-------+-----------+
// | 1  | m1.tiny   |   512 |    1 |         0 |     1 | True      |
// | 2  | m1.small  |  2048 |   20 |         0 |     1 | True      |
// | 3  | m1.medium |  4096 |   40 |         0 |     2 | True      |
// | 4  | m1.large  |  8192 |   80 |         0 |     4 | True      |
// | 5  | m1.xlarge | 16384 |  160 |         0 |     8 | True      |
// +----+-----------+-------+------+-----------+-------+-----------+
type Flavours string

const (
	Tiny   Flavours = "1"
	Small  Flavours = "2"
	Medium Flavours = "3"
	Large  Flavours = "4"
	XLarge Flavours = "5"
)

func (f Flavours) Value() string {
	return string(f)
}
