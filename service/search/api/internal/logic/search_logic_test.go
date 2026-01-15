package logic

import (
	"bookNew/service/search/api/internal/svc"
	"bookNew/service/search/api/internal/types"
	"bookNew/service/user/rpc/user"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// 1. 定义 Mock 对象，继承 testify 的 mock.Mock
type MockUserRpc struct {
	mock.Mock
}

// 2. 实现 User 接口 (必须和 user_client.User 接口签名一致)
func (m *MockUserRpc) GetUser(ctx context.Context, in *user.IdReq, opts ...grpc.CallOption) (*user.UserInfoReply, error) {
	// 这一行是精髓：告诉 mock 框架，“有人调用了我”，并获取预设的返回值
	args := m.Called(ctx, in)

	// 断言返回值类型
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserInfoReply), args.Error(1)
}

func TestSearchLogic_Search(t *testing.T) {
	// --- A. 准备工作 (Arrange) ---

	// 1. 初始化 Mock 对象
	mockRpc := new(MockUserRpc)

	// 2. 构造 ServiceContext，把 Mock 对象注入进去！
	// 注意：这里我们不需要真实的 Config 和数据库，只要 mock 住 RPC 即可
	svcCtx := &svc.ServiceContext{
		UserRpc: mockRpc, // 关键点：偷梁换柱
	}

	// 3. 模拟 Context (带上 userId，模拟 JWT 解析后的结果)
	// 注意：go-zero 解析出来的是 json.Number
	ctx := context.WithValue(context.Background(), "userId", json.Number("1"))

	// 4. 实例化 Logic
	logic := NewSearchLogic(ctx, svcCtx)
	// 关掉日志输出，保持测试控制台清爽
	logic.Logger = logx.WithContext(ctx)

	// --- B. 设定 Mock 行为 (Stubbing) ---
	// 对应 Java: when(rpc.getUser(any, arg)).thenReturn(reply, nil)

	expectedUser := &user.UserInfoReply{
		Id:     1,
		Name:   "Mock张三", // 注意这里我们用 Mock 的名字
		Number: "8888",
	}

	// 设定规则：当 GetUser 被调用，且参数 Id=1 时，返回 expectedUser 和 nil 错误
	mockRpc.On("GetUser", mock.Anything, &user.IdReq{Id: 1}).Return(expectedUser, nil)

	// --- C. 执行测试 (Act) ---
	req := &types.SearchReq{
		Name: "Go并发编程",
	}
	resp, err := logic.Search(req)

	// --- D. 断言结果 (Assert) ---
	// 使用 testify/assert 库，语法非常友好
	assert.Nil(t, err)                    // 期望没有错误
	assert.NotNil(t, resp)                // 期望有返回值
	assert.Equal(t, "Mock张三", resp.Owner) // 核心验证：Logic 是否正确使用了 RPC 的返回值
	assert.Equal(t, 100, resp.Count)

	// 验证 Mock 方法是否真的被调用了 (Verify)
	mockRpc.AssertExpectations(t)
}

func TestSearchLogic_Search_RpcError(t *testing.T) {
	// Arrange
	mockRpc := new(MockUserRpc)
	svcCtx := &svc.ServiceContext{UserRpc: mockRpc}
	ctx := context.WithValue(context.Background(), "userId", json.Number("99")) // 假设这是 ID 99
	logic := NewSearchLogic(ctx, svcCtx)

	// Stubbing: 模拟 RPC 返回错误
	// Java: when(rpc.getUser(...)).thenThrow(...)
	mockRpc.On("GetUser", mock.Anything, &user.IdReq{Id: 99}).Return(nil, errors.New("RPC连接超时"))

	// Act
	resp, err := logic.Search(&types.SearchReq{Name: "Go"})

	// Assert
	assert.NotNil(t, err)                   // 期望报错
	assert.Nil(t, resp)                     // 期望响应为空
	assert.Equal(t, "RPC连接超时", err.Error()) // 验证错误信息

	mockRpc.AssertExpectations(t)
}
