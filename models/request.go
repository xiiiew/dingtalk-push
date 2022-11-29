package models

type DingtalkSendRequest struct {
	Secret      string      `json:"secret"`
	AccessToken string      `json:"access_token"`
	Message     interface{} `json:"message"`
}
