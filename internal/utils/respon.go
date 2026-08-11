package utils


// respon success
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// respon error
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}