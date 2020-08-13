package authentication

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Zucke/ContactManager/internal/data"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateBytes, err := ioutil.ReadFile("./private.rsa")
	if err != nil {
		log.Fatal("error reading private key")
	}

	publicBytes, err := ioutil.ReadFile("./public.rsa.pub")

	if err != nil {
		log.Fatal("error reading public key")
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("error parsing private key")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("error parsing private key")
	}
}

//GenerateJWT generate a JWT token to a user
func GenerateJWT(user data.User) (string, error) {
	claims := data.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "log a user",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return result, nil

}

//ValidateToken velidate a token from a logged user
func ValidateToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	var token *jwt.Token
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &data.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	return token, err

}
