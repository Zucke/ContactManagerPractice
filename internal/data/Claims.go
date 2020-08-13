package data

import jwt "github.com/dgrijalva/jwt-go"

//Claim has the necesary data for jwt claims
type Claim struct {
	User `json:"user"`
	jwt.StandardClaims
}
