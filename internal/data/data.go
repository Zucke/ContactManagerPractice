package data

import "errors"

//DBNAME is the name of the database we are working
const DBNAME = "contactmanager"

//all names of collections
const (
	UserColletion     = "user"
	ContactCollection = "contacts"
)

//Some error status
var (
	ErrorUserExist = errors.New("the user exist")
	ErrorNotFount  = errors.New("Not found")
)
