package user

import (
	"encoding/json"

	log "github.com/cihub/seelog"
)

type loginData struct {
	Clientid     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

func buildLoginData(email, pwd, clientid, clientsecret string) []byte {
	l := &loginData{
		Clientid:     clientid,
		ClientSecret: clientsecret,
		Email:        email,
		Password:     pwd,
	}
	data, err := json.Marshal(l)
	if err != nil {
		log.Errorf("build login data fail: %s", err)
	}
	return data
}
