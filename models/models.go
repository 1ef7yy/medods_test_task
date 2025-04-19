package models

type GenerateTokenRequest struct {
	IP   string
	Guid string
}

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
