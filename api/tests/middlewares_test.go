package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/primarybank/api"
	"github.com/primarybank/token"
	"github.com/stretchr/testify/require"
)

func addAuthz(
	t *testing.T,
	tokenMaker token.Maker,
	req *http.Request,
	authzType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authzHeader := fmt.Sprintf("%s %s", authzType, token)
	req.Header.Set(api.AuthHeaderKey, authzHeader)
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		setupAuth func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		checkResp func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Valid",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthz(t, tokenMaker, req, api.AuthType, "user", time.Minute)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthz",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnSupportedAuthz",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthz(t, tokenMaker, req, "Invalid", "user", time.Minute)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthz(t, tokenMaker, req, api.AuthType, "user", -time.Minute)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.Router.GET(authPath, api.AuthMiddleWare(server.TokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tt.setupAuth(t, req, server.TokenMaker)
			server.Router.ServeHTTP(recorder, req)
			tt.checkResp(t, recorder)
		})
	}
}
