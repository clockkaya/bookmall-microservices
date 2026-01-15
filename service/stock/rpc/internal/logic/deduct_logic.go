package logic

import (
	"context"

	"bookNew/service/stock/rpc/internal/svc"
	"bookNew/service/stock/rpc/stock"

	"github.com/zeromicro/go-zero/core/logx"

	// 🔥 1. 必须引入这两个包
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductLogic {
	return &DeductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Deduct 扣减库存 (Saga Action)
func (l *DeductLogic) Deduct(in *stock.DeductReq) (*stock.DeductReply, error) {
	l.Infof("🔥 [Stock Action] 正在扣减图书 %d 的库存: %d", in.BookId, in.Count)

	// 直接调用 Model 方法，不涉及任何 SQL
	err := l.svcCtx.StockModel.Deduct(l.ctx, in.BookId, in.Count)
	if err != nil {
		l.Errorf("❌ 扣减库存失败: %v", err)

		// 🔥 2. 关键修改：将 error 包装成 Aborted 状态码
		// DTM 只有看到 Aborted 或 Failed 才会停止重试并触发回滚
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &stock.DeductReply{}, nil
}
