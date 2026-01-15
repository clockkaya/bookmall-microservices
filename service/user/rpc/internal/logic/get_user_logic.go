// 我们要写逻辑的地方
package logic

import (
	"context"
	"errors" // 用于返回错误

	"bookNew/service/user/rpc/internal/svc"
	"bookNew/service/user/rpc/model" // 引入 model
	"bookNew/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*func (l *GetUserLogic) GetUser(in *user.IdReq) (*user.UserInfoReply, error) {
	// 模拟：不管查谁，都返回这个“张三”
	return &user.UserInfoReply{
		Id:     in.Id,
		Name:   "张三 (来自 RPC)",
		Number: "2023001",
		Gender: "男",
	}, nil
}*/

func (l *GetUserLogic) GetUser(in *user.IdReq) (*user.UserInfoReply, error) {
	// 1. 调用 Model 查询数据库 (自带缓存!)
	one, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)

	if err != nil {
		// 处理“未找到”的情况
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 2. 将数据库实体转为 RPC 响应格式
	return &user.UserInfoReply{
		Id:     one.Id,
		Name:   one.Name,
		Number: one.Number,
		Gender: one.Gender,
	}, nil

	//return nil, xerr.NewErrMsg("这本书被借完了！")
}
