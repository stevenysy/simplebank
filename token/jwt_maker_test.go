package token

//
//import (
//	"fmt"
//	"github.com/golang-jwt/jwt/v5"
//	"github.com/stevenysy/simplebank/util"
//	"github.com/stretchr/testify/require"
//	"testing"
//	"time"
//)
//
//func TestJWTMaker(t *testing.T) {
//	maker, err := NewJWTMaker(util.RandomString(32))
//	require.NoError(t, err)
//
//	username := util.RandomOwner()
//	duration := time.Minute
//
//	issuedAt := time.Now()
//	expiredAt := issuedAt.Add(duration)
//
//	token, err := maker.CreateToken(username, duration)
//	require.NoError(t, err)
//	require.NotEmpty(t, token)
//
//	payload, err := maker.VerifyToken(token)
//	require.NoError(t, err)
//	require.NotEmpty(t, payload)
//
//	require.NotEmpty(t, payload.ID)
//	require.Equal(t, username, payload.Subject)
//	require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
//	require.WithinDuration(t, issuedAt, payload.NotBefore.Time, time.Second)
//	require.WithinDuration(t, expiredAt, payload.ExpiresAt.Time, time.Second)
//}
//
//func TestExpiredJWTToken(t *testing.T) {
//	maker, err := NewJWTMaker(util.RandomString(32))
//	require.NoError(t, err)
//
//	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
//	require.NoError(t, err)
//	require.NotEmpty(t, token)
//
//	payload, err := maker.VerifyToken(token)
//	fmt.Println(err)
//	require.Error(t, err)
//	require.Nil(t, payload)
//}
//
//func TestInvalidJWTTokenAlgNone(t *testing.T) {
//	payload, err := NewPayload(util.RandomOwner(), time.Minute)
//	require.NoError(t, err)
//	require.NotEmpty(t, payload)
//
//	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
//	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
//	require.NoError(t, err)
//
//	maker, err := NewJWTMaker(util.RandomString(32))
//	require.NoError(t, err)
//
//	payload, err = maker.VerifyToken(token)
//	require.Error(t, err)
//	require.Nil(t, payload)
//}
