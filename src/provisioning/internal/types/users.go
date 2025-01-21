package types

type Class string

const (
	Admin Class = "admin" //nolint:all
	Customer Class = "user" //nolint:all
)

func (c Class) Value() string {
	return string(c)
}

type User struct {
	ID    int64
	Sub   string
	Name  string `json:"name"`
	Class string `json:"class"`
}
