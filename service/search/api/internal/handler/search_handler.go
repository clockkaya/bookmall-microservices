// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"bookNew/common/response"
	"net/http"

	"bookNew/service/search/api/internal/logic"
	"bookNew/service/search/api/internal/svc"
	"bookNew/service/search/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

/*func searchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSearchLogic(r.Context(), svcCtx)
		resp, err := l.Search(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}*/

func searchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchReq
		if err := httpx.Parse(r, &req); err != nil {
			// 参数解析错误，也走统一响应
			httpx.ErrorCtx(r.Context(), w, err)
			// 注意：httpx.Parse 错误是框架级的，如果你想连这个也统一，
			// 需要更高级的全局拦截，但这里我们先不改，保持简单。
			return
		}

		l := logic.NewSearchLogic(r.Context(), svcCtx)
		resp, err := l.Search(&req)

		// 2. 核心修改：不管成功失败，全交给 Response 处理
		response.Response(w, resp, err)
	}
}
