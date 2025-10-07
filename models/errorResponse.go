package models

type errorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err string) *errorResponse {
	return &errorResponse{
		Error: err,
	}
}
