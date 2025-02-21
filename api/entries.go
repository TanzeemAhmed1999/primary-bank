package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/primarybank/db/sqlc"
)

func (s *Server) CreateEntry(ctx *gin.Context) {
	var req CreateEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	args := db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	entry, err := s.store.CreateEntry(ctx.Request.Context(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

func (s *Server) GetEntry(ctx *gin.Context) {
	var req GetEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	entry, err := s.store.GetEntry(ctx.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

func (s *Server) DeleteEntry(ctx *gin.Context) {
	var req DeleteEntryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	err := s.store.DeleteEntry(ctx.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "entry deleted"})
}

func (s *Server) ListEntries(ctx *gin.Context) {
	var req ListEntriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResp(err))
		return
	}

	args := db.ListEntriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	transfers, err := s.store.ListEntries(ctx.Request.Context(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResp(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}
