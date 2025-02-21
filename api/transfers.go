package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/primarybank/db/sqlc"
)

func (s *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	args := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := s.store.CreateTransfer(ctx.Request.Context(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (s *Server) getTransfer(ctx *gin.Context) {
	var req getTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	transfer, err := s.store.GetTransfer(ctx.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (s *Server) deleteTransfer(ctx *gin.Context) {
	var req deleteTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	err := s.store.DeleteTransfer(ctx.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "transfer deleted"})
}

func (s *Server) listTransfers(ctx *gin.Context) {
	var req listTransfersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	args := db.ListTransfersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	transfers, err := s.store.ListTransfers(ctx.Request.Context(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}
