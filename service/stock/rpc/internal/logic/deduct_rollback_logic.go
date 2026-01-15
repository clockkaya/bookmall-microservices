package logic

import (
	"context"

	"bookNew/service/stock/rpc/internal/svc"
	"bookNew/service/stock/rpc/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeductRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductRollbackLogic {
	return &DeductRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 补偿库存 (Compensate) - 如果后续步骤失败，把库存加回来
func (l *DeductRollbackLogic) DeductRollback(in *stock.DeductReq) (*stock.DeductReply, error) {
	// todo: add your logic here and delete this line

	return &stock.DeductReply{}, nil
}
