// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
)

type TimeListenerMiddleware struct {
}

func NewTimeListenerMiddleware() *TimeListenerMiddleware {
	return &TimeListenerMiddleware{}
}

// Handle 就像 Java 的 doFilter
func (m *TimeListenerMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// --- 🟢 前置逻辑 (Pre-Handle) ---
		startTime := time.Now()
		logx.Info(">>> 请求进来了: ", r.RequestURI)

		// --- 🟡 放行，执行业务 Logic (Call Next) ---
		// 这行代码非常关键！如果不调用 next，请求就断在这里了，Logic 根本不会执行。
		next(w, r)

		// --- 🔴 后置逻辑 (Post-Handle) ---
		duration := time.Since(startTime)
		logx.WithContext(r.Context()).Infof("<<< 请求结束: %s, 耗时: %v", r.RequestURI, duration)
	}
}
