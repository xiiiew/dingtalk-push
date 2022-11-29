package bootstrap

import (
	"dingtalk-push/models"
	"dingtalk-push/web/middleware"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

type Configurator func(*Bootstrapper)

// Bootstrapper继承和共享 iris.Application
type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time

	Sessions *sessions.Sessions
}

// Return a new Bootstrapper
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		Application:  iris.New(),
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
	}

	b.Configure(cfgs...)

	return b
}

// 配置
func (b *Bootstrapper) Bootstrapper(cs ...Configurator) {
	// CORS
	b.UseGlobal(middleware.MiddleCors)

	b.setUpErrorHandlers()

	b.UseGlobal(recover.New())

	// 配置到所有url
	mvc.New(b.Party("/"))

	// 为所有handle配置Done
	b.SetExecutionRules(iris.ExecutionRules{
		Done: iris.ExecutionOptions{Force: true},
	})
	b.Done(middleware.MiddleDone)

	b.Configure(cs...)
}

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

// 监听端口
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}

// 错误页面
func (b *Bootstrapper) setUpErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		resp, _ := json.Marshal(models.ErrorPageNotFound)
		ctx.Write(resp)
		middleware.MiddleDone(ctx)
		return
	})
}
