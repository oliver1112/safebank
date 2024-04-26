package lib

import (
	"encoding/base64"
	"encoding/json"
)

type UserToken struct {
	UserID int64 `json:"userID"`
}

func (u *UserToken) DecodeToken(token string) {
	tokenEncrypt, _ := base64.StdEncoding.DecodeString(token)
	_ = json.Unmarshal(tokenEncrypt, &u)
}

func (u *UserToken) EncodeToken() string {
	dataJson, _ := json.Marshal(u)
	return base64.StdEncoding.EncodeToString(dataJson)
}
