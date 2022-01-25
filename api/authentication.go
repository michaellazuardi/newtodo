package api

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
}

//why do we need this struct?
type authCustomClaims struct {
	Name               string `json:"name"`
	User               bool   `json:"user"`
	jwt.StandardClaims        //you can use type Standard claims from jwt library
}

type jwtServiceImpl struct {
	secretKey string
	issuer    string
}

func getSecretKey() string {
	//Getenv retrieves the value of the environment variable named by the key.
	//It returns the value, which will be empty if the variable is not present.
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func JWTAuthService() JWTService {
	return &jwtServiceImpl{
		secretKey: getSecretKey(),
		issuer:    "Michael",
	}
}

func (s *jwtServiceImpl) GenerateToken(email string, isUser bool) string {
	claims := &authCustomClaims{
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	//encoded string
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (s *jwtServiceImpl) ValidateToken(encodedToken string) (*jwt.Token, error) {

	//I dont know how to interpret this code
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])

		}
		return []byte(s.secretKey), nil
	})
}
