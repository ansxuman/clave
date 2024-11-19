package objects

import (
	"time"
)

type TOTPProfile struct {
	ID        string    `json:"id"`
	Issuer    string    `json:"issuer"`
	Secret    string    `json:"secret"`
	Period    int       `json:"period"`
	Digits    int       `json:"digits"`
	CreatedAt time.Time `json:"createdAt"`
}

type TOTPCode struct {
	Code      string `json:"code"`
	ExpiresIn int    `json:"expiresIn"`
}
