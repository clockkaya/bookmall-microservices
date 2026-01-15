package logic

import (
	"bookNew/service/user/rpc/internal/svc"
	"bookNew/service/user/rpc/user"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecrPointsRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecrPointsRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecrPointsRollbackLogic {
	return &DecrPointsRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Saga Compensate: 补偿积分 (把扣掉的加回来)
func (l *DecrPointsRollbackLogic) DecrPointsRollback(in *user.AdjustPointsReq) (*user.AdjustPointsReply, error) {
	// 直接调用 Model 层的业务方法
	err := l.svcCtx.UserModel.DecrPointsRollback(l.ctx, in.Id, in.Points)
	if err != nil {
		return nil, err
	}
	return &user.AdjustPointsReply{}, nil
}
