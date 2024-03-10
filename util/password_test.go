package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	pw := RandomString(6)

	hashedPw1, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPw1)

	err = CheckPassword(pw, hashedPw1)
	require.NoError(t, err)

	wrongPw := RandomString(6)
	err = CheckPassword(wrongPw, hashedPw1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPw2, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPw1)
	require.NotEqual(t, hashedPw1, hashedPw2)
}
