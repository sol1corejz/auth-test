package models

type AuthTokens struct {
	AccessToken        string
	RefreshToken       string
	HashedRefreshToken string
}
