package lib

import (
	"encoding/base64"
	"encoding/json"
)

type UserToken struct {
	UserID    int64  `json:"userID"`
	ExpiresAt int64  `json:"expiresAt"`
	IP        string `json:"ip"`
}

func (u *UserToken) DecodeToken(token string) {
	tokenEncrypt, _ := base64.StdEncoding.DecodeString(token)
	tokenEncrypt, _ = RsaDecrypt(tokenEncrypt)
	_ = json.Unmarshal(tokenEncrypt, &u)
}

func (u *UserToken) EncodeToken() string {
	dataJson, _ := json.Marshal(u)
	tokenEncrypt, _ := RsaEncrypt(dataJson)
	return base64.StdEncoding.EncodeToString(tokenEncrypt)
}
