package token

//
//import (
//	"fmt"
//	"github.com/golang-jwt/jwt/v5"
//	"time"
//)
//
//const minSecretKeySize = 32
//
//// JWTMaker is a JSON Web Token maker
//type JWTMaker struct {
//	secretKey string
//}
//
//// NewJWTMaker creates a new JWTMaker
//func NewJWTMaker(secretKey string) (Maker, error) {
//	if len(secretKey) < minSecretKeySize {
//		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
//	}
//	return &JWTMaker{secretKey: secretKey}, nil
//}
//
//// CreateToken creates a new token for a specific user and duration
//func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
//	payload, err := NewPayload(username, duration)
//	if err != nil {
//		return "", err
//	}
//
//	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
//	return jwtToken.SignedString([]byte(maker.secretKey))
//}
//
//// VerifyToken checks if the input token is valid
//func (maker *JWTMaker) VerifyToken(token string) (*jwt.RegisteredClaims, error) {
//	keyFunc := func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, ErrInvalidToken
//		}
//		return []byte(maker.secretKey), nil
//	}
//
//	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
//	if err != nil {
//		return nil, err
//	}
//
//	payload, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
//	if !ok {
//		return nil, ErrInvalidToken
//	}
//
//	return payload, nil
//}
