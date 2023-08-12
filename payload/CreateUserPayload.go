package payload

type CreateUserPayload struct {
	Username  string `json:"username"`
	UserEmail string `json:"user_email"`
	Address   string `json:"address"`
	Tel       string `json:"tel"`
	PId       string `json:"p_id"`
	Image     string `json:"image"`
}
