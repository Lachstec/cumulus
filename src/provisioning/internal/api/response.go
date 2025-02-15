package api

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *Error      `json:"error,omitempty"`
}

type Error struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
