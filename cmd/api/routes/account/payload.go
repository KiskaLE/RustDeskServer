package account

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutPayload struct {
	RefreshToken string `json:"refresh_token"`
}

type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token"`
}
