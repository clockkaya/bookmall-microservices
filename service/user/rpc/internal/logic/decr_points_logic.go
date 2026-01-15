package logic

import (
	"bookNew/service/user/rpc/internal/svc"
	"bookNew/service/user/rpc/user"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecrPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecrPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecrPointsLogic {
	return &DecrPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Saga Action: 扣积分
func (l *DecrPointsLogic) DecrPoints(in *user.AdjustPointsReq) (*user.AdjustPointsReply, error) {
	// 直接调用 Model 层的业务方法
	err := l.svcCtx.UserModel.DecrPoints(l.ctx, in.Id, in.Points)
	if err != nil {
		return nil, err
	}
	return &user.AdjustPointsReply{}, nil
}
