package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/primarybank/config"
	db "github.com/primarybank/db/sqlc"
	"github.com/primarybank/token"
)

// Server serves all http request for banking service
type Server struct {
	store      db.Store
	Router     *gin.Engine
	TokenMaker token.Maker
	Config     config.Config
}

func NewServer(cfg config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		TokenMaker: tokenMaker,
		Config:     cfg,
	}
	server.setUpRouter()

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	// add routes to the routes
	// User routes
	router.POST("/user", server.CreateUser)
	router.POST("/user/login", server.LoginUser)

	authRoutes := router.Group("/").Use(AuthMiddleWare(server.TokenMaker))

	// routes require auth
	authRoutes.GET("/user/:username", server.GetUser)
	authRoutes.PUT("/user/:username", server.UpdateUser)

	// Account routes
	authRoutes.GET("/account/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccounts)
	authRoutes.POST("/account", server.CreateAccount)
	authRoutes.PATCH("/account", server.UpdateAccount)
	authRoutes.DELETE("/account/:id", server.DeleteAccount)

	// Transfer routes
	authRoutes.GET("/transfer/:id", server.GetTransfer)
	authRoutes.GET("/transfers", server.ListTransfers)
	authRoutes.POST("/transfer", server.CreateTransfer)
	authRoutes.DELETE("/transfer/:id", server.DeleteTransfer)

	// Entry routes
	authRoutes.GET("/entry/:id", server.GetEntry)
	authRoutes.GET("/entries", server.ListEntries)
	authRoutes.POST("/entry", server.CreateEntry)
	authRoutes.DELETE("/entry/:id", server.DeleteEntry)

	server.Router = router
}

func (s *Server) Start(addr string) error {
	return s.Router.Run(addr)
}

func errResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}
