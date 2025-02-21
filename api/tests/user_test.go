package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func createRandomUser(t *testing.T) db.User {
	password := commonutils.RandomString(10)
	hashedPassword, err := commonutils.HashPassword(password)
	require.NoError(t, err)

	return db.User{
		Username:  commonutils.RandomString(8),
		Password:  hashedPassword,
		FullName:  commonutils.RandomString(10),
		Email:     commonutils.RandomEmail(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)

	user := createRandomUser(t)

	testCases := []struct {
		name         string
		requestBody  api.CreateUserRequest
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name: "Valid Request",
			requestBody: api.CreateUserRequest{
				Username: user.Username,
				Password: user.Password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(user, nil).
					Times(1)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid Request - Empty Username",
			requestBody: api.CreateUserRequest{
				Username: "",
				Password: user.Password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid Request - Invalid Email",
			requestBody: api.CreateUserRequest{
				Username: user.Username,
				Password: user.Password,
				FullName: user.FullName,
				Email:    "invalid-email",
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
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
			c.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")

			server.CreateUser(c)

			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	server := api.NewServer(store)

	user := createRandomUser(t)

	testCases := []struct {
		name         string
		username     string
		requestBody  api.UpdateUserRequest
		buildStubs   func(store *mocks.MockStore)
		expectedCode int
	}{
		{
			name:     "Valid Update Without Password",
			username: user.Username,
			requestBody: api.UpdateUserRequest{
				FullName: "Updated Name",
				Email:    "updated.email@example.com",
				Password: "",
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Return(user, nil).
					Times(1)

				updatedUser := user
				updatedUser.FullName = "Updated Name"
				updatedUser.Email = "updated.email@example.com"

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(updatedUser, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:     "Valid Update With New Password",
			username: user.Username,
			requestBody: api.UpdateUserRequest{
				FullName: "Updated Name",
				Email:    "updated.email@example.com",
				Password: "NewSecurePassword123",
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Return(user, nil).
					Times(1)

				hashedPassword, _ := commonutils.HashPassword("NewSecurePassword123")

				updatedUser := user
				updatedUser.FullName = "Updated Name"
				updatedUser.Email = "updated.email@example.com"
				updatedUser.Password = hashedPassword

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(updatedUser, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:     "User Not Found",
			username: "non_existing_user",
			requestBody: api.UpdateUserRequest{
				FullName: "Updated Name",
				Email:    "updated.email@example.com",
				Password: "",
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), "non_existing_user").
					Return(db.User{}, sql.ErrNoRows).
					Times(1)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:     "Internal Server Error on GetUser",
			username: user.Username,
			requestBody: api.UpdateUserRequest{
				FullName: "Updated Name",
				Email:    "updated.email@example.com",
				Password: "",
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Return(db.User{}, errors.New("database error")).
					Times(1)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0) // Should not be called
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:     "Internal Server Error on UpdateUser",
			username: user.Username,
			requestBody: api.UpdateUserRequest{
				FullName: "Updated Name",
				Email:    "updated.email@example.com",
				Password: "",
			},
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Return(user, nil).
					Times(1)

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(db.User{}, errors.New("update failed")).
					Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/users/"+tc.username, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			c.Params = gin.Params{{Key: "username", Value: tc.username}}

			server.UpdateUser(c)

			require.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}
