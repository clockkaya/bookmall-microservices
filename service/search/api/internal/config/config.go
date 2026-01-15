//加入 RPC 配置定义
// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc" // 引入 zrpc
)

type Config struct {
	rest.RestConf
	// 新增：UserRpc 的客户端配置
	UserRpc zrpc.RpcClientConf
	// 新增：对应 YAML 里的 Auth 块
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	// 新增 Dtm 服务器地址配置
	DtmServer string
}
