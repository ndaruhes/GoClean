package responses

type BasicResponse struct {
	StatusCode int
	Error      error
	Data       interface{}
	Pagination int
	Message    string
}
