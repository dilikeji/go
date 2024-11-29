package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"go-shop/App"
)

func main() {
	s := g.Server()
	s.Use(GlobalMiddleware)
	SetApiRoute(s)
	s.EnableAdmin()
	s.Run()
}

// GlobalMiddleware 全局中间件
func GlobalMiddleware(gReq *ghttp.Request) {
	var (
		err        error
		redisValue *gvar.Var
		header     *gvar.Var
	)
	// 执行业务逻辑前
	g.Log().Info(gReq.Context(), "Before Request")
	// 获取请求头 鉴权 配置中未设置authHeader 或没有对应请求头则跳过 鉴权移到业务代码中获取上下文时
	if header, err = g.Cfg().Get(gReq.Context(), "server.authHeader"); err == nil {
		auth := gReq.Header.Get(header.String())
		// 无视类型输出变量到控制台
		g.Dump(auth)
		// TODO 处理请求头并添加上下文变量 例如redis
		if !g.IsEmpty(auth) {
			// redis获取变量
			if redisValue, err = g.Redis().Get(gReq.Context(), auth); err != nil {
				// redis连接失败报错
				gReq.Response.WriteJsonExit(g.Map{
					"code": App.CodeToken,
					"msg":  err.Error(),
					"data": make([]interface{}, 0),
				})
				return
			}
			if !redisValue.IsEmpty() {
				gReq.SetCtxVar("Auth", redisValue)
				redisValue = g.NewVar(g.Map{
					"timeNow": gtime.Now().String(),
				})
				// redis设置变量 有效期秒
				err = g.Redis().SetEX(gReq.Context(), auth, redisValue, 1000)
			}
		}
	}
	// TODO 验签
	gReq.Middleware.Next()
	// 执行业务逻辑后
	g.Log().Info(gReq.Context(), "After Request")
}

func SetApiRoute(s *ghttp.Server) {
	s.Group("v1/api", func(group *ghttp.RouterGroup) {
		group.GET("login", App.Login)
	})
}
