package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/stevenysy/simplebank/db/mock"
	db "github.com/stevenysy/simplebank/db/sqlc"
	"github.com/stevenysy/simplebank/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type eqCreateUserParamsMatcher struct {
	arg db.CreateUserParams
	pw  string
}

func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	if err := util.CheckPassword(e.pw, arg.HashedPassword); err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	// Check if types assignable and convert them to common type
	x1Val := reflect.ValueOf(e.arg)
	x2Val := reflect.ValueOf(x)

	if x1Val.Type().AssignableTo(x2Val.Type()) {
		x1ValConverted := x1Val.Convert(x2Val.Type())
		return reflect.DeepEqual(x1ValConverted.Interface(), x2Val.Interface())
	}

	return false
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and %v", e.arg, e.pw)
}

func EqCreateUserParams(arg db.CreateUserParams, pw string) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		arg: arg,
		pw:  pw,
	}
}

func TestCreateUserApi(t *testing.T) {
	user, pw := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  pw,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, pw)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				rsp := userResponse{
					Username:          user.Username,
					FullName:          user.FullName,
					Email:             user.Email,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}
				requireBodyMatchUser(t, recorder.Body, rsp)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "abc123#",
				"password":  pw,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
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
	}, pw
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, userRsp userResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var receivedUserRsp userResponse
	err = json.Unmarshal(data, &receivedUserRsp)
	require.NoError(t, err)
	require.Equal(t, userRsp, receivedUserRsp)
}
