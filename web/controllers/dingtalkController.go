package controllers

import (
	"dingtalk-push/chat"
	"dingtalk-push/models"
	"encoding/json"
	"github.com/kataras/iris/v12"
)

type DingtalkController struct {
	Ctx iris.Context
}

// 发送钉钉消息
func (self *DingtalkController) PostSend() iris.Map {
	req := models.DingtalkSendRequest{}
	err := self.Ctx.ReadJSON(&req)
	if err != nil {
		return models.ErrorResponse("failed to unmarshal message")
	}
	if req.Secret == "" {
		return models.ErrorResponse("error secret")
	}
	if req.AccessToken == "" {
		return models.ErrorResponse("error access_token")
	}
	bytesMsg, err := json.Marshal(req.Message)
	if err != nil {
		return models.ErrorResponse(err.Error())
	}
	im, err := chat.UnmarshalBytes(bytesMsg)
	if err != nil {
		return models.ErrorResponse(err.Error())
	}
	success := chat.BotPoolInstance.Send(req.Secret, req.AccessToken, im)
	if success {
		return models.SuccessResponse("success")
	}
	return models.SuccessResponse("failed")
}
