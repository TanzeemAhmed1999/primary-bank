package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/primarybank/api"
	commonutils "github.com/primarybank/common/utils"
	"github.com/primarybank/db/mocks"
	db "github.com/primarybank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) *db.Account {
	return &db.Account{
		ID:        commonutils.RandomMoney(),
		Owner:     commonutils.RandomOwner(),
		Balance:   commonutils.RandomMoney(),
		Currency:  commonutils.RandomCurrency(),
		CreatedAt: time.Now(),
	}
}

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)

	account := CreateRandomAccount(t)
	testCases := []struct {
		name         string
		requestBody  api.CreateAccountRequest
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name: "Valid Request",
			requestBody: api.CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Return(*account, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid Request - Empty Owner",
			requestBody: api.CreateAccountRequest{
				Owner:    "",
				Currency: account.Currency,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
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
			c.Request = httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			server.CreateAccount(c)

			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestGetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)

	account := CreateRandomAccount(t)
	testCases := []struct {
		name         string
		accountID    string
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:      "Valid Request",
			accountID: strconv.Itoa(int(account.ID)),
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Return(*account, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:      "Not Found",
			accountID: "999",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(db.Account{}, sql.ErrNoRows).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: tc.accountID})
			c.Request = httptest.NewRequest(http.MethodGet, "/accounts/"+tc.accountID, nil)

			server.GetAccount(c)

			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)

	account := CreateRandomAccount(t)
	testCases := []struct {
		name         string
		accountID    string
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:      "Valid Request",
			accountID: strconv.Itoa(int(account.ID)),
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), account.ID).Return(nil).Times(1)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:      "Not Found",
			accountID: "999",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: tc.accountID})
			c.Request = httptest.NewRequest(http.MethodDelete, "/accounts/"+tc.accountID, nil)

			server.DeleteAccount(c)

			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestListAccounts(t *testing.T) {
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
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Return([]db.Account{}, nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Invalid Query Params",
			queryParams: "page_id=abc&page_size=-1",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, "/accounts?"+tc.queryParams, nil)

			server.ListAccounts(c)

			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)

	account := CreateRandomAccount(t)
	testCases := []struct {
		name         string
		requestBody  api.UpdateAccountRequest
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name: "Valid Request",
			requestBody: api.UpdateAccountRequest{
				ID:      account.ID,
				Balance: 500,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name: "Account Not Found",
			requestBody: api.UpdateAccountRequest{
				ID:      999,
				Balance: 500,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(tc.requestBody)
			c.Request = httptest.NewRequest(http.MethodPut, "/accounts", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			server.UpdateAccount(c)

			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}
