package responses

type SuccessResponse struct {
	SuccessCode string
	StatusCode  int
	Data        interface{}
}

type ErrorResponse struct {
	Error      error
	StatusCode int
	FormErrors map[string][]string
	Data       interface{}
}

type TokenDecoded struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
