package response

type ApiResponse struct {
	Message string `json:"message" example:"Operation successful"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Error   string            `json:"error" example:"An error occurred"`
	Details map[string]string `json:"details"`
}
