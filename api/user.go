package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	commonutils "github.com/primarybank/common/utils"
	db "github.com/primarybank/db/sqlc"
)

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := commonutils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
		Email:    req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (server *Server) GetUser(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hashedPassword *string
	if req.Password != "" {
		hash, err := commonutils.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		hashedPassword = &hash
	}

	arg := db.UpdateUserParams{
		Username: username,
		FullName: req.FullName,
		Email:    req.Email,
		// Temporary hack, need to find a way to generate this field a optional
		Password: user.Password,
	}

	if hashedPassword != nil {
		arg.Password = *hashedPassword
	}

	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := commonutils.CheckPassword(req.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := s.Config.AccessTokenDuration
	token, err := s.TokenMaker.CreateToken(user.Username, expirationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
	}

	ctx.JSON(http.StatusOK, LoginUserResponse{
		AccessToken: token,
		ExpiresIn:   time.Now().Add(expirationTime).Unix(),
	})
}
