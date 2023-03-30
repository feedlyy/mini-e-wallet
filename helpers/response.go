package helpers

type Response struct {
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type ErrResp struct {
	Err interface{} `json:"error,omitempty"`
}
