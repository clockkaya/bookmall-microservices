package svc

import (
	"bookNew/service/stock/rpc/internal/config"
	"bookNew/service/stock/rpc/model" // 引入生成的 model 包

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 🔥 核心修复：定义 StockModel 字段，Logic 层才能调用 l.svcCtx.StockModel
	StockModel model.StockModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		// 🔥 核心修复：初始化 StockModel
		StockModel: model.NewStockModel(conn, c.CacheRedis),
	}
}
