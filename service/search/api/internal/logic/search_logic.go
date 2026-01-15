// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"bookNew/service/user/rpc/user_client"
	"context"
	"encoding/json" // 需要引入 json 包
	"fmt"

	"bookNew/service/search/api/internal/svc"
	"bookNew/service/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

/*func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchReply, err error) {
	// 1. 调用 RPC
	// 就像调用本地方法一样简单，不用管网络连接、序列化等细节
	userResp, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user_client.IdReq{
		Id: req.UserId,
	})

	if err != nil {
		return nil, err
	}

	// 2. 组装结果
	return &types.SearchReply{
		Name:  req.Name,
		Count: 100,
		// 这里的数据来自远程服务
		Owner: userResp.Name,
	}, nil
}*/

func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchReply, err error) {
	// 1. 从 Context 中获取 userId
	// 注意：go-zero 解析 JWT 后，数字类型默认会变成 json.Number 类型
	userIdNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	userId, err := userIdNumber.Int64()
	if err != nil {
		return nil, err
	}

	// 2. 调用 RPC (使用解析出来的 userId)
	userResp, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user_client.IdReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.SearchReply{
		Name:  req.Name,
		Count: 100,
		Owner: userResp.Name,
	}, nil
}
