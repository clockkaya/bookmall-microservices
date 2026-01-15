package logic

import (
	"bookNew/service/search/api/internal/svc"
	"bookNew/service/search/api/internal/types"
	"bookNew/service/stock/rpc/stock"
	"bookNew/service/user/rpc/user"
	"context"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BorrowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBorrowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowLogic {
	return &BorrowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BorrowLogic) Borrow(req *types.BorrowReq) (resp *types.BorrowReply, err error) {
	// 1. 定义 DTM 服务器地址 (从配置读)
	dtmServer := l.svcCtx.Config.DtmServer

	// 2. 生成一个全局唯一的事务 ID (GID)
	// 在真实场景中，可以使用 snowflake 算法或者 orderId
	gid := dtmgrpc.MustGenGid(dtmServer)
	l.Infof("🚀 开启 Saga 分布式事务，GID: %s", gid)

	// 3. 定义子服务的直连地址
	// DTM 是个独立进程，它需要回调我们的 RPC 服务，所以这里要给它“怎么访问我”的地址
	// 注意：UserRpc 端口是 8080，StockRpc 端口是 8081
	userRpcTarget := "127.0.0.1:8080"
	stockRpcTarget := "127.0.0.1:8081"

	// 4. 构造 User 服务的参数
	userReq := &user.AdjustPointsReq{
		Id:     req.UserId,
		Points: 10, // 假设借一本书扣 10 分
	}
	// 5. 构造 Stock 服务的参数
	stockReq := &stock.DeductReq{
		BookId: req.BookId,
		Count:  1, // 借 1 本
	}

	// 6. 核心：编排 Saga 剧本 📜
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		// 🎬 第一幕：User 服务
		Add(
			userRpcTarget+"/user.user/decrPoints",         // 正向操作 (Action)
			userRpcTarget+"/user.user/decrPointsRollback", // 补偿操作 (Compensate)
			userReq, // 参数
		).
		// 🎬 第二幕：Stock 服务
		Add(
			stockRpcTarget+"/stock.stock/deduct",         // 正向操作
			stockRpcTarget+"/stock.stock/deductRollback", // 补偿操作
			stockReq, // 参数
		)

	// 🔥 关键：告诉 DTM 我要同步等待结果！
	// 如果不加这行，Submit 只是把任务丢进队列就返回成功，你看不出后续的失败
	saga.WaitResult = true

	// 7. 提交事务！(Action!)
	// WaitResult=true 表示等待所有子事务执行完才返回，方便我们看结果
	err = saga.Submit()

	if err != nil {
		l.Errorf("❌ 借书失败，Saga 提交错误: %v", err)
		return nil, err
	}

	l.Infof("✅ 借书成功！Saga 事务 %s 已完成", gid)
	return &types.BorrowReply{}, nil
}
