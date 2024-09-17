package requests

type GenerateTokenRequest struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`
}
