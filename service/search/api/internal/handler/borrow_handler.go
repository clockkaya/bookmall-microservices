// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"bookNew/service/search/api/internal/logic"
	"bookNew/service/search/api/internal/svc"
	"bookNew/service/search/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func borrowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BorrowReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewBorrowLogic(r.Context(), svcCtx)
		resp, err := l.Borrow(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
