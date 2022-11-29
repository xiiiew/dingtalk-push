package router

import (
	"dingtalk-push/bootstrap"
	"dingtalk-push/web/controllers"
	"github.com/kataras/iris/v12/mvc"
)

// 创建路由
func RegisterRouter(b *bootstrap.Bootstrapper) {
	mvc.New(b.Party("/dingtalk")).
		Handle(new(controllers.DingtalkController))
}
