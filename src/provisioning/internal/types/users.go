package types

type User struct {
	ID   int
	Sub  string
	Name string `json:"name"`
	Role string `json:"role"`
}
