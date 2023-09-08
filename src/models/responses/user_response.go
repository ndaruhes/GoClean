package responses

type LoginResponse struct {
	Role  string `json:"role"`
	Token string `json:"token"`
}

type CurrentUserResponse struct {
	Name string `json:"name"`
}
