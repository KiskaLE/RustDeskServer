package account

type LoginPayload struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type LogoutPayload struct {
	RefreshToken string `json:"RefreshToken"`
}

type RegisterPayload struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"RefreshToken"`
}
