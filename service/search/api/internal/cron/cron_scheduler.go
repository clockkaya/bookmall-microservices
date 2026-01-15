package cron

import (
	"bookNew/service/search/api/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

// Setup 初始化并启动定时任务
// 返回 *cron.Cron 是为了让 main 函数能在服务停止时优雅关闭它
func Setup(svcCtx *svc.ServiceContext) *cron.Cron {
	// 1. 创建 cron 实例
	// WithSeconds() 是关键，否则默认最小粒度是分钟，测试起来太慢
	//c := cron.New(cron.WithSeconds())
	// 使用 Chain 包装器
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	// 2. 添加任务
	// 这里的 spec "*/5 * * * * *" 表示每 5 秒执行一次
	_, err := c.AddFunc("*/5 * * * * *", func() {
		// 这里可以调用 Logic 层的方法，或者直接写逻辑
		// 比如：svcCtx.UserModel.Find...
		logx.Info("⏰ [Cron] 定时任务触发：正在刷新热门搜索缓存...")

		// 模拟业务逻辑：打印一下 Redis 配置信息（证明拿到了 svcCtx）
		// logx.Infof("当前 Redis 地址: %s", svcCtx.Config.UserRpc.Etcd.Hosts)
	})

	if err != nil {
		logx.Errorf("添加定时任务失败: %v", err)
	}

	// 3. 启动调度器（这是异步的，不会阻塞主线程）
	c.Start()
	logx.Info("🚀 Cron 调度器已启动")

	return c
}
