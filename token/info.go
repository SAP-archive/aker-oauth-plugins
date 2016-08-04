package token

type Info struct {
	UserID   string   `json:"user_id"`
	UserName string   `json:"user_name"`
	Scope    []string `json:"scope"`
}
