// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"bookNew/service/search/api/internal/config"
	"bookNew/service/search/api/internal/middleware"
	"bookNew/service/user/rpc/user_client" // 引入生成的客户端代码

	"github.com/zeromicro/go-zero/rest" // 引入 rest
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	UserRpc      user_client.User // 定义接口
	TimeListener rest.Middleware  // 2. 定义中间件字段 (注意名字要和 .api 里的一致)
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// 初始化 zRPC 客户端
		UserRpc: user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
		// 3. 初始化中间件
		// 注意：这里需要调用 .Handle 方法，因为框架需要的是 func(next) func
		TimeListener: middleware.NewTimeListenerMiddleware().Handle,
	}
}
