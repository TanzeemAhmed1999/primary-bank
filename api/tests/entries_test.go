package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
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

func CreateRandomEntry(t *testing.T) *db.Entry {
	return &db.Entry{
		ID:        commonutils.RandomMoney(),
		AccountID: commonutils.RandomMoney(),
		Amount:    commonutils.RandomMoney(),
		CreatedAt: time.Now(),
	}
}

func TestCreateEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	testCases := []struct {
		name         string
		requestBody  api.CreateEntryRequest
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name: "Valid Request",
			requestBody: api.CreateEntryRequest{
				AccountID: 1,
				Amount:    100,
			},
			buildStubs: func(store *mocks.MockStore) {
				entry := CreateRandomEntry(t)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Return(*entry, nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Invalid Request - Negative Amount",
			requestBody: api.CreateEntryRequest{
				AccountID: 1,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
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
			c.Request = httptest.NewRequest(http.MethodPost, "/entries", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")
			server.CreateEntry(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestGetEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	testCases := []struct {
		name         string
		entryID      string
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:    "Valid Request",
			entryID: "1",
			buildStubs: func(store *mocks.MockStore) {
				entry := CreateRandomEntry(t)
				store.EXPECT().GetEntry(gomock.Any(), gomock.Any()).Return(*entry, nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:    "Entry Not Found",
			entryID: "999",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().GetEntry(gomock.Any(), gomock.Any()).Return(db.Entry{}, sql.ErrNoRows).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: tc.entryID})
			c.Request = httptest.NewRequest(http.MethodGet, "/entries/"+tc.entryID, nil)
			server.GetEntry(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	testCases := []struct {
		name         string
		entryID      string
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:    "Valid Request",
			entryID: "1",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().DeleteEntry(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:    "Entry Not Found",
			entryID: "999",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().DeleteEntry(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows).Times(1)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: tc.entryID})
			c.Request = httptest.NewRequest(http.MethodDelete, "/entries/"+tc.entryID, nil)
			server.DeleteEntry(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestListEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)
	testCases := []struct {
		name         string
		queryParams  api.ListEntriesRequest
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name: "Valid Request",
			queryParams: api.ListEntriesRequest{
				PageSize: 10,
				PageID:   1,
			},
			buildStubs: func(store *mocks.MockStore) {
				entries := []db.Entry{*CreateRandomEntry(t), *CreateRandomEntry(t)}
				store.EXPECT().ListEntries(gomock.Any(), gomock.Any()).Return(entries, nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Internal Server Error",
			queryParams: api.ListEntriesRequest{
				PageSize: 10,
				PageID:   1,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().ListEntries(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, "/entries?page_size="+strconv.Itoa(int(tc.queryParams.PageSize))+"&page_id="+strconv.Itoa(int(tc.queryParams.PageID)), nil)
			server.ListEntries(c)
			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}
