package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/primarybank/api"
	"github.com/primarybank/db/mocks"
	db "github.com/primarybank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer() *db.Transfer {
	return &db.Transfer{
		ID:            int64(1),
		FromAccountID: int64(1),
		ToAccountID:   int64(2),
		Amount:        100,
		CreatedAt:     time.Now(),
	}
}

func TestCreateTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	transfer := CreateRandomTransfer()
	testCases := []struct {
		name         string
		requestBody  api.CreateTransferRequest
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name: "Valid Request",
			requestBody: api.CreateTransferRequest{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().CreateTransfer(gomock.Any(), gomock.Any()).Return(*transfer, nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid Request - Negative Amount",
			requestBody: api.CreateTransferRequest{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        -100,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().CreateTransfer(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")
			server.CreateTransfer(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestGetTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	transfer := CreateRandomTransfer()
	testCases := []struct {
		name         string
		transferID   string
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:       "Valid Request",
			transferID: "1",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Any()).Return(*transfer, nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: tc.transferID})
			c.Request = httptest.NewRequest(http.MethodGet, "/transfers/"+tc.transferID, nil)
			server.GetTransfer(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestDeleteTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	testCases := []struct {
		name         string
		transferID   string
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:       "Valid Request",
			transferID: "1",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().DeleteTransfer(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: tc.transferID})
			c.Request = httptest.NewRequest(http.MethodDelete, "/transfers/"+tc.transferID, nil)
			server.DeleteTransfer(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestListTransfers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	testCases := []struct {
		name         string
		queryParams  string
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:        "Valid Request",
			queryParams: "page_id=1&page_size=5",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().ListTransfers(gomock.Any(), gomock.Any()).Return([]db.Transfer{}, nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, "/transfers?"+tc.queryParams, nil)
			server.ListTransfers(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}
