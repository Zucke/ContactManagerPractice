package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//User is the user data
type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Nickname string             `json:"nickname" bson:"nickname"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

//UserData is the db data
type UserData struct {
	data *Data
	coll *mongo.Collection
}

//ComparePassword macth with a password
func (u *User) ComparePassword(password string) bool {
	return u.Password == password
}

//VarifyUserByNickname encharge of query if a user exist
func (ud *UserData) VarifyUserByNickname(ctx context.Context, nickname string) (User, error) {
	result := User{}
	err := ud.coll.FindOne(ctx, bson.M{"nickname": nickname}).Decode(&result)
	if err != nil {
		return result, ErrorNotFount
	}
	return result, nil

}

//NewUser add a new User
func (ud *UserData) NewUser(ctx context.Context, info *User) error {
	if _, err := ud.VarifyUserByNickname(ctx, info.Nickname); err == nil {
		return ErrorUserExist
	}

	_, err := ud.coll.InsertOne(ctx, &info)
	return err

}

//DeleteUser delete a user by id
func (ud *UserData) DeleteUser(ctx context.Context) error {
	userContact := NewUserContact()
	err := userContact.DeleteAll(ctx)
	if err != nil {
		return err
	}
	_, err = ud.coll.DeleteOne(ctx, bson.M{"_id": ctx.Value(primitive.ObjectID{})})
	return err
}

//NewUserData reference the db data and colletion
func NewUserData() *UserData {
	return &UserData{
		data: New(),
		coll: data.DBCollection(UserColletion),
	}

}
