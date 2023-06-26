package responses

type BasicResponse struct {
	StatusCode  int
	Error       error
	MessageCode string
	Data        interface{}
}

type TokenDecoded struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
