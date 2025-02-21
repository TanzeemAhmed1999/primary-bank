package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/primarybank/db/sqlc"
)

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(ctx.Request.Context(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	account, err := s.store.GetAccount(ctx.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (s *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	args := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := s.store.ListAccounts(ctx.Request.Context(), args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (s *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	args := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	err := s.store.UpdateAccount(ctx.Request.Context(), args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// We should ideally soft delete account
func (s *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	err := s.store.DeleteAccount(ctx.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
