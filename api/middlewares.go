package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/primarybank/token"
)

const (
	AuthHeaderKey   = "authorization"
	AuthType        = "bearer"
	AuthzPayloadKey = "authz_payload"
)

func AuthMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authzHeader := ctx.GetHeader(AuthHeaderKey)
		if len(authzHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResp(err))
			return
		}

		parts := strings.Fields(authzHeader)
		if len(parts) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResp(err))
			return
		}

		if strings.ToLower(parts[0]) != AuthType {
			err := errors.New("authorization type is not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResp(err))
			return
		}

		accessToken := parts[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			err := errors.New("invalid token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResp(err))
			return
		}

		ctx.Set(AuthzPayloadKey, payload)
		ctx.Next()
	}
}
