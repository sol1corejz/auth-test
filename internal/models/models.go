package models

type AuthTokens struct {
	AccessToken        string
	RefreshToken       string
	HashedRefreshToken string
}

type UserAccessData struct {
	UserID       string `db:"user_id" json:"user_id"`
	UserIP       string `db:"user_ip" json:"user_ip"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type AccessData struct {
	ID int `db:"id" json:"id"`
	UserAccessData
}
