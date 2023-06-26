package responses

type SuccessResponse struct {
	StatusCode  int
	SuccessCode string
	Data        interface{}
}

type ErrorResponse struct {
	StatusCode int
	Error      error
	ErrorCode  string
	Data       interface{}
}

type TokenDecoded struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
