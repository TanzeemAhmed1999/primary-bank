package tests

import (
	"testing"
	"time"

	"github.com/primarybank/api"
	commonutils "github.com/primarybank/common/utils"
	"github.com/primarybank/config"
	db "github.com/primarybank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *api.Server {
	cfg := config.Config{
		TokenSymmetricKey:   commonutils.RandomString(35),
		AccessTokenDuration: 3 * time.Minute,
	}

	server, err := api.NewServer(cfg, store)
	require.NoError(t, err)

	return server
}
