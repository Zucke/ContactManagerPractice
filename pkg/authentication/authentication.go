package authentication

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Zucke/ContactManager/pkg/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the user data
type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Nickname string             `json:"nickname" bson:"nickname"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

//Claim contiaint the claims that use the token
type Claim struct {
	User `json:"user"`
	jwt.StandardClaims
}

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

//ComparePassword macth with a password
func (u *User) ComparePassword(password string) bool {
	return u.Password == password
}

func init() {
	privateBytes, err := ioutil.ReadFile("./cert/private.rsa")
	if err != nil {
		log.Fatal("error reading private key")
	}

	publicBytes, err := ioutil.ReadFile("./cert/public.rsa.pub")

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
func GenerateJWT(user User) (string, error) {
	claims := Claim{
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

//ValidateMiddleware used to validate tokes
func ValidateMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token *jwt.Token
		token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &Claim{}, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		if !token.Valid {
			response.HTTPError(w, r, http.StatusUnauthorized, "Invalid Token")
			return
		}
		id := token.Claims.(*Claim).ID
		ctx := context.WithValue(r.Context(), primitive.ObjectID{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
	return fn

}
