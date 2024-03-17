package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/stevenysy/simplebank/db/mock"
	db "github.com/stevenysy/simplebank/db/sqlc"
	"github.com/stevenysy/simplebank/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTransferApi(t *testing.T) {
	amount := int64(10)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)
	acc1 := randomAccount(user1.Username)
	acc2 := randomAccount(user2.Username)
	acc3 := randomAccount(user3.Username)

	acc1.Currency = util.USD
	acc2.Currency = util.USD
	acc3.Currency = util.CNY

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Expect GetAccount to be called to validate currency
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).
					Times(1).
					Return(acc1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).
					Times(1).
					Return(acc2, nil)

				arg := db.TransferTxParams{
					FromAccountID: acc1.ID,
					ToAccountID:   acc2.ID,
					Amount:        amount,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": acc3.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Expect GetAccount to be called to validate currency
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc3.ID)).
					Times(1).
					Return(acc3, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc3.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Expect GetAccount to be called to validate currency
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).
					Times(1).
					Return(acc1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc3.ID)).
					Times(1).
					Return(acc3, nil)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).
					Times(1).
					Return(acc1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			body: gin.H{
				"from_account_id": -1,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "AccountInternalServerError",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TxInternalServerError",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Expect GetAccount to be called to validate currency
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).
					Times(1).
					Return(acc1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).
					Times(1).
					Return(acc2, nil)

				arg := db.TransferTxParams{
					FromAccountID: acc1.ID,
					ToAccountID:   acc2.ID,
					Amount:        amount,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// Start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
