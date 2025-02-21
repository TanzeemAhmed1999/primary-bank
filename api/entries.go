package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/primarybank/db/sqlc"
)

func (s *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest
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

func (s *Server) getEntry(ctx *gin.Context) {
	var req getEntryRequest
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

func (s *Server) deleteEntry(ctx *gin.Context) {
	var req deleteEntryRequest
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

func (s *Server) listEntries(ctx *gin.Context) {
	var req listEntriesRequest
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
