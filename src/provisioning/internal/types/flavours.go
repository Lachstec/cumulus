package types

type Flavour struct {
	ID string
	Name string
	RAM	int
	Disk int
	Ephemeral int
	VCPUs int
	Is_Public bool
}

var Flavours = []Flavour {
	{
		ID: "1",
		Name: "m1.tiny",
		RAM: 512,
		Disk: 1,
		Ephemeral: 0,
		VCPUs: 1,
		Is_Public: true,
	},
	{
		ID: "2",
		Name: "m1.small",
		RAM: 2048,
		Disk: 20,
		Ephemeral: 0,
		VCPUs: 1,
		Is_Public: true,
	},
	{
		ID: "3",
		Name: "m1.medium",
		RAM: 4096,
		Disk: 40,
		Ephemeral: 0,
		VCPUs: 2,
		Is_Public: true,
	},
	{
		ID: "4",
		Name: "m1.large",
		RAM: 8192,
		Disk: 80,
		Ephemeral: 0,
		VCPUs: 4,
		Is_Public: true,
	},
	{
		ID: "5",
		Name: "m1.xlarge",
		RAM: 16384,
		Disk: 160,
		Ephemeral: 0,
		VCPUs: 8,
		Is_Public: true,
	},
}