package requests

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOAuthRequest struct {
	Code string `json:"code"`
}

type RegisterWithEmailPasswordRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required"`
	PhoneNumber          string `json:"phoneNumber" validate:"required"`
	Password             string `json:"password" validate:"required"`
	ConfirmationPassword string `json:"confirmationPassword" validate:"required,eqfield=Password"`
}

type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}
