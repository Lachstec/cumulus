package types

type Class string

const (
	Admin Class = "admin" //nolint:all
	Customer Class = "customer" //nolint:all
)

type User struct {
	ID   int64
	Sub  string
	Name string `json:"name"`
	Class string `json:"class"`
}
