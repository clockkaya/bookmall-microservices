package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StockModel = (*customStockModel)(nil)

type (
	// 1. 在接口里定义业务方法
	StockModel interface {
		stockModel
		Deduct(ctx context.Context, bookId int64, count int64) error
		DeductRollback(ctx context.Context, bookId int64, count int64) error
	}

	customStockModel struct {
		*defaultStockModel
	}
)

func NewStockModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StockModel {
	return &customStockModel{
		defaultStockModel: newStockModel(conn, c, opts...),
	}
}

// 2. 实现扣减库存 (Action)
func (m *customStockModel) Deduct(ctx context.Context, bookId int64, count int64) error {
	// 先查询获取缓存 Key (通过唯一索引 book_id 查询)
	// 注意：这里我们假设 Unique Key 是 book_id，生成的 FindOneByBookId
	stock, err := m.FindOneByBookId(ctx, bookId)
	if err != nil {
		return err
	}

	// 🔥 硬编码缓存 Key，避免 IDE 报红和依赖 _gen.go 变量
	// 这里的 Key 格式要和 _gen.go 里的逻辑保持一致 (通常是 cache:表名:索引字段:)
	// 如果不确定，可以在 redis-cli 里 keys * 看看 goctl 生成的 key 长啥样
	// 这里假设是标准的 cache:stock:id: 和 cache:stock:bookId:
	stockIdKey := fmt.Sprintf("%s%v", "cache:stock:id:", stock.Id)
	stockBookIdKey := fmt.Sprintf("%s%v", "cache:stock:bookId:", stock.BookId)

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		// 核心 SQL：扣减库存，且保证 num >= count
		query := fmt.Sprintf("update %s set num = num - ? where book_id = ? and num >= ?", m.table)
		result, err = conn.ExecCtx(ctx, query, count, bookId, count)
		if err != nil {
			return nil, err
		}

		// 检查是否有行被更新
		affected, _ := result.RowsAffected()
		if affected == 0 {
			// 如果没更新，说明库存不足 (num < count)
			return nil, fmt.Errorf("库存不足")
		}
		return result, nil
	}, stockIdKey, stockBookIdKey) // 传入 Key 自动清理缓存

	return err
}

// 3. 实现回滚库存 (Compensate)
func (m *customStockModel) DeductRollback(ctx context.Context, bookId int64, count int64) error {
	stock, err := m.FindOneByBookId(ctx, bookId)
	if err != nil {
		return err
	}

	stockIdKey := fmt.Sprintf("%s%v", "cache:stock:id:", stock.Id)
	stockBookIdKey := fmt.Sprintf("%s%v", "cache:stock:bookId:", stock.BookId)

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		// 核心 SQL：加回库存
		query := fmt.Sprintf("update %s set num = num + ? where book_id = ?", m.table)
		return conn.ExecCtx(ctx, query, count, bookId)
	}, stockIdKey, stockBookIdKey)

	return err
}
