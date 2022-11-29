package chat

import (
	"crypto/hmac"
	"crypto/sha256"
	"dingtalk-push/conf"
	"dingtalk-push/httprequest"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Errors
const (
	ErrBadAccessToken = "bad access token"
	ErrBadSecret      = "bad secret"
	ErrSign           = "bad Sign"
)

// Dingtalk bot API.
type BotAPI struct {
	Secret      string `json:"secret"`
	AccessToken string `json:"access_token"`

	Client *http.Client `json:"-"`
}

// Dingtalk response
type BotResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func NewBotAPI(secret, accessToken string) (*BotAPI, error) {
	timeout := time.Duration(conf.ConfigYamlInstance.DingtalkConfig.TimeOut) * time.Second
	return NewBotAPIWithClient(secret, accessToken, &http.Client{Timeout: timeout})
}

func NewBotAPIWithClient(secret, accessToken string, client *http.Client) (*BotAPI, error) {
	bot := &BotAPI{
		Secret:      secret,
		AccessToken: accessToken,
		Client:      client,
	}
	return bot, nil
}

// 签名
func (self *BotAPI) sign() (url.Values, error) {
	v := url.Values{}
	if self.AccessToken != "" {
		v.Add("access_token", self.AccessToken)
	} else {
		return v, errors.New(ErrBadAccessToken)
	}
	if self.Secret == "" {
		return v, errors.New(ErrBadSecret)
	}

	timestamp := int(time.Now().Unix()) * 1000
	v.Add("timestamp", strconv.Itoa(timestamp))

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, self.Secret)
	hmacCtx := hmac.New(sha256.New, []byte(self.Secret))
	_, err := hmacCtx.Write([]byte(stringToSign))
	if err != nil {
		return v, errors.New(ErrSign)
	}
	cipherStr := hmacCtx.Sum(nil)
	sign := url.QueryEscape(base64.StdEncoding.EncodeToString(cipherStr))
	v.Add("sign", sign)
	return v, nil
}

// 发送消息
func (self *BotAPI) Send(m IMessage) (_ bool, err error) {
	// set message type
	m.SetMessageType()
	v, err := self.sign()
	if err != nil {
		return
	}
	u := fmt.Sprintf("%s?%s", APIEndpoint, v.Encode())
	bytes, err := httprequest.HTTPPostJsonWithClient(u, m, self.Client)
	if err != nil {
		return
	}
	botResponse := BotResponse{}
	err = json.Unmarshal(bytes, &botResponse)
	if err != nil {
		return
	}
	if botResponse.Errcode != 0 {
		err = errors.New(botResponse.Errmsg)
		return
	}
	return true, nil
}
