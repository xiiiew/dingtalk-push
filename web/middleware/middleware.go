package middleware

import (
	"dingtalk-push/conf"
	utils "dingtalk-push/utils/http"
	loggerRequest "dingtalk-push/utils/log/request"
	loggerResponse "dingtalk-push/utils/log/response"
	"fmt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
)

// 访问日志
func MiddleRequest(ctx iris.Context) {
	ctx.Record()
	ip := utils.GetIP(ctx)
	token := ctx.GetHeader("token")
	uri := ctx.Request().URL.Path
	params := ctx.URLParams()
	method := ctx.Request().Method
	header := ctx.Request().Header

	log := fmt.Sprintf("token=%s, ip=%s, method=%s, uri=%s, params=%+v, header=%+v", token, ip, method, uri, params, header)
	loggerRequest.Info(log)
}

// CORS
func MiddleCors(ctx iris.Context) {
	if conf.ConfigYamlInstance.HTTPConfig.Cors.Enable {
		ctx.Header("Access-Control-Allow-Origin", conf.ConfigYamlInstance.HTTPConfig.Cors.AccessControlAllowOrigin)
		if ctx.Request().Method == "OPTIONS" {
			ctx.Header("Access-Control-Allow-Methods", conf.ConfigYamlInstance.HTTPConfig.Cors.AccessControlAllowMethods)
			ctx.Header("Access-Control-Allow-Headers", conf.ConfigYamlInstance.HTTPConfig.Cors.AccessControlAllowHeaders)
			ctx.Header("Content-type", "application/json;charset=UTF-8")
			ctx.StatusCode(204)
			return
		}
	}
	MiddleRequest(ctx)
	ctx.Next()
	//MiddleDone(ctx)
}

// 所有handle的Done
func MiddleDone(ctx iris.Context) {
	ip := utils.GetIP(ctx)
	token := ctx.GetHeader("token")
	uri := ctx.Request().URL.Path
	params := ctx.URLParams()
	method := ctx.Request().Method
	header := ctx.Request().Header

	formData := ctx.FormValues()
	body, _ := ioutil.ReadAll(ctx.Request().Body)
	ctx.Request().Body.Close()
	response := ctx.Recorder().Body()

	loggerResponse.InfoWithFields(
		uri,
		loggerResponse.Fields{
			"token":     token,
			"ip":        ip,
			"method":    method,
			"params":    params,
			"form-data": formData,
			"body":      string(body),
			"response":  string(response),
			"header":    header,
		},
	)

	ctx.Recorder().ResetBody()
	ctx.Recorder().Write(response)
}
