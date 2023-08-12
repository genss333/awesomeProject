package payload

type AuthPayload struct {
	Username  string `json:"username"`
	UserEmail string `json:"user_email"`
}
