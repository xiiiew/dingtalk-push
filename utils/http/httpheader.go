package utils

import (
	"github.com/kataras/iris/v12"
	"strings"
)

// 取ip地址
func GetIP(ctx iris.Context) string {
	ipXOFF := ctx.Request().Header.Get("X-Original-Forwarded-For")
	ipXFF := ctx.Request().Header.Get("X-Forwarded-For")
	ipRI := ctx.Request().Header.Get("X-Real-IP")
	if ipXOFF != "" {
		ipXOFF = strings.Replace(ipXOFF, "[", "", -1)
		ipXOFF = strings.Replace(ipXOFF, "]", "", -1)
		ipXOFFList := strings.Split(ipXOFF, ",")
		if len(ipXOFFList) > 0 {
			return ipXOFFList[0]
		}
	} else if ipXFF != "" {
		ipXFF = strings.Replace(ipXFF, "[", "", -1)
		ipXFF = strings.Replace(ipXFF, "]", "", -1)
		ipXFFList := strings.Split(ipXFF, ",")
		if len(ipXFFList) > 0 {
			return ipXFFList[0]
		}
	} else if ipRI != "" {
		return ipRI
	}
	return ctx.Request().RemoteAddr
}
