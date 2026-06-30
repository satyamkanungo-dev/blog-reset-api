package apiresponse

type APIResponse struct {
	Code    int    `json:"code"`   // http status
	Status  string `json:"status"` // success or error
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
