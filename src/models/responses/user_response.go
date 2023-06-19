package responses

type LoginResponse struct {
	TokenID string `json:"tokenId"`
}

type CurrentUserResponse struct {
	Name string `json:"name"`
}
