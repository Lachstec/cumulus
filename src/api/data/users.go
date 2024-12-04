package data

type User struct {
	ID int
	Sub string
	Name string `json:"name"`
	Role string `json:"role"`
}

var Users []User

func init() {
	Users = []User {
		{
			ID: 1,
			Sub: "github|42710663",
			Name: "Janosch",
			Role: "Admin",
		},
		{
			ID: 2,
			Sub: "",
			Name: "Morids",
			Role: "Admin",
		},
		{
			ID: 3,
			Sub: "",
			Name: "Loen",
			Role: "User",
		},
		{
			ID: 4,
			Sub: "auth0|6739c35a61e7cc0097c5edfd",
			Name: "Moritzmann",
			Role: "User",
		},
	}
}