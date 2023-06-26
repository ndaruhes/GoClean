package responses

type LoginResponse struct {
	Token string `json:"token"`
}

type CurrentUserResponse struct {
	Name string `json:"name"`
}
