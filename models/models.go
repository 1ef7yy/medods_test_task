package models

type GenerateTokenRequest struct {
	IP         string
	Guid       string
	Generation int
}

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AccessToken struct {
	Guid       string
	Generation int
}

type RefreshToken struct {
	Guid string
	IP   string
}

type RefreshTokenRequest struct {
	IP     string
	Tokens Token
}
