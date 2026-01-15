package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// 新增：数据库配置
	Mysql struct {
		DataSource string
	}

	// 新增：缓存配置
	CacheRedis cache.CacheConf
}
