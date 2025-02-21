package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/primarybank/db/sqlc"
)

// Server serves all http request for banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	// add routes to the routes
	// Account routes
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/account", server.createAccount)
	router.PATCH("/account", server.updateAccount)
	router.DELETE("/account/:id", server.deleteAccount)

	// Transfer routes
	router.GET("/transfer/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfers)
	router.POST("/transfer", server.createTransfer)
	router.DELETE("/transfer/:id", server.deleteTransfer)

	// Entry routes
	router.GET("/entry/:id", server.getEntry)
	router.GET("/entries", server.listEntries)
	router.POST("/entry", server.createEntry)
	router.DELETE("/entry/:id", server.deleteEntry)

	server.router = router
	return server
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}
