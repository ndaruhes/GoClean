package requests

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOAuthRequest struct {
	Code string `json:"code"`
}

type RegisterWithEmailPasswordRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmationPassword"`
}

type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}
