package svc

import (
	"bookNew/service/user/rpc/internal/config"
	"bookNew/service/user/rpc/model" // 1. 引入 model 包

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 2. 定义 Model 接口
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 3. 初始化数据库连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		// 4. 初始化 Model (传入 conn 和 CacheRedis 配置)
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
