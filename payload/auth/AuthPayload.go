package auth

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
