package payload

type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
