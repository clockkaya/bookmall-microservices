package main

import (
	"bookNew/common/xerr" // 引入你的错误包
	"flag"
	"fmt"
	"net/http"

	"bookNew/service/search/api/internal/config"
	"bookNew/service/search/api/internal/handler"
	"bookNew/service/search/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/search-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// --- 👇 核心修改看这里 👇 ---
	server := rest.MustNewServer(c.RestConf,
		// 1. 接管 401 (Token 错误/缺失)
		rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
			httpx.OkJson(w, &xerr.CodeErrorResponse{
				Code: xerr.TOKEN_EXPIRE_ERROR,
				Msg:  "用户未登录或Token已失效",
			})
		}),
		// 2. 接管 404 (接口不存在)
		rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpx.OkJson(w, &xerr.CodeErrorResponse{
				Code: xerr.SERVER_COMMON_ERROR,
				Msg:  "请求的接口不存在",
			})
		})),
	)
	// --- 👆 核心修改结束 👆 ---

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 接管业务层错误
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *xerr.CodeError:
			return http.StatusOK, e.Data()
		default:
			fmt.Printf("❌ [全局错误拦截] 发生系统错误: %v\n", err)
			return http.StatusOK, &xerr.CodeErrorResponse{
				Code: xerr.SERVER_COMMON_ERROR,
				Msg:  "系统繁忙或参数错误",
			}
		}
	})

	//// --- 🔥 新增代码 Start ---
	//// 启动定时任务
	//// 传入 ctx，这样 Cron 任务里就能用数据库和 Redis 了
	//scheduler := cron.Setup(ctx)
	//// 优雅关闭：当 main 函数退出（服务停止）时，停止定时任务
	//defer scheduler.Stop()
	//// --- 🔥 新增代码 End ---

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
