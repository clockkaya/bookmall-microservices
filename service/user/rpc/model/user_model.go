package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// 1. 在接口里增加这两个方法，Logic 层才能调用
	UserModel interface {
		userModel
		DecrPoints(ctx context.Context, id int64, points int64) error
		DecrPointsRollback(ctx context.Context, id int64, points int64) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

// 2. 实现扣减积分 (Action)
func (m *customUserModel) DecrPoints(ctx context.Context, id int64, points int64) error {
	user, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", "cache:user:id:", user.Id)
	userNumberKey := fmt.Sprintf("%s%v", "cache:user:number:", user.Number)

	// 执行事务/更新
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set points = points - ? where id = ? and points >= ?", m.table)
		return conn.ExecCtx(ctx, query, points, id, points)
	}, userIdKey, userNumberKey)

	return err
}

// 3. 实现回滚积分 (Compensate)
func (m *customUserModel) DecrPointsRollback(ctx context.Context, id int64, points int64) error {
	user, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", "cache:user:id:", user.Id)
	userNumberKey := fmt.Sprintf("%s%v", "cache:user:number:", user.Number)

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set points = points + ? where id = ?", m.table)
		return conn.ExecCtx(ctx, query, points, id)
	}, userIdKey, userNumberKey)

	return err
}
