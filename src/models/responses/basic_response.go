package responses

type SuccessResponse struct {
	SuccessCode string
	StatusCode  int
	Data        interface{}
}

type ErrorResponse struct {
	Error      error
	ErrorCode  string
	StatusCode int
	Data       interface{}
}

type TokenDecoded struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
