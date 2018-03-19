package requests

type AuthLoginOrSignupEmail struct {
	Email string `json:"email"`
}

type AuthSignupEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginOrSignupSSO struct {
	Token    string `json:"token"`
	Provider string `json:"provider"`
}
