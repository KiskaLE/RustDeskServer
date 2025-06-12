package user

type LoginPayload struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type RegisterPayload struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}
