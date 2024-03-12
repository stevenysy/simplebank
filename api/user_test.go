package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/stevenysy/simplebank/db/mock"
	db "github.com/stevenysy/simplebank/db/sqlc"
	"github.com/stevenysy/simplebank/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserApi(t *testing.T) {
	user, pw := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username":  user.Username,
				"password":  pw,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				rsp := createUserResponse{
					Username:          user.Username,
					FullName:          user.FullName,
					Email:             user.Email,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}
				requireBodyMatchUser(t, recorder.Body, rsp)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(&tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (db.User, string) {
	pw := util.RandomString(6)
	hashedPw, err := util.HashPassword(pw)
	require.NoError(t, err)

	return db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPw,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}, hashedPw
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, userRsp createUserResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var receivedUserRsp createUserResponse
	err = json.Unmarshal(data, &receivedUserRsp)
	require.NoError(t, err)
	require.Equal(t, userRsp, receivedUserRsp)
}
