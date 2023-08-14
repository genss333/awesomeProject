package json

type AuthJson struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	TokenExp int64  `json:"token_exp"`
}

type AuthTokenJson struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	TokenExp int64  `json:"token_exp"`
}
