package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/primarybank/db/sqlc"
)

// Server serves all http request for banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
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

	server.router = router
	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}
