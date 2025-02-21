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
	router     *gin.Engine
	tokenMaker token.Maker
	config     config.Config
}

func NewServer(cfg config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     cfg,
	}
	server.setUpRouter()

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	// add routes to the routes
	// Account routes
	router.GET("/account/:id", server.GetAccount)
	router.GET("/accounts", server.ListAccounts)
	router.POST("/account", server.CreateAccount)
	router.PATCH("/account", server.UpdateAccount)
	router.DELETE("/account/:id", server.DeleteAccount)

	// Transfer routes
	router.GET("/transfer/:id", server.GetTransfer)
	router.GET("/transfers", server.ListTransfers)
	router.POST("/transfer", server.CreateTransfer)
	router.DELETE("/transfer/:id", server.DeleteTransfer)

	// Entry routes
	router.GET("/entry/:id", server.GetEntry)
	router.GET("/entries", server.ListEntries)
	router.POST("/entry", server.CreateEntry)
	router.DELETE("/entry/:id", server.DeleteEntry)

	// User routes
	router.POST("/user", server.CreateUser)
	router.GET("/user/:username", server.GetUser)
	router.PUT("/user/:username", server.UpdateUser)
	router.POST("/user/login", server.LoginUser)

	server.router = router
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}
