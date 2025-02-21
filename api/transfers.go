package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/primarybank/db/sqlc"
)

// TODO: add a validation to check both FromAccount and ToAccount currency matches
func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
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

func (s *Server) GetTransfer(ctx *gin.Context) {
	var req GetTransferRequest
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

func (s *Server) DeleteTransfer(ctx *gin.Context) {
	var req DeleteTransferRequest
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

func (s *Server) ListTransfers(ctx *gin.Context) {
	var req ListTransfersRequest
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
