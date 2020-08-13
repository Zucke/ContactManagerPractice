package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Contact contain the info of a contact and the owner
type Contact struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Number string             `json:"number" bson:"number"`
	Name   string             `json:"name" bson:"name"`
}

//UserContact db info
type UserContact struct {
	data *Data
	coll *mongo.Collection
}

//GetByName get by name the contacts
func (uc UserContact) GetByName(ctx context.Context, name string) (Contact, error) {
	result := Contact{}
	err := uc.coll.FindOne(ctx, bson.M{"name": name, "user_id": ctx.Value(primitive.ObjectID{})}).Decode(&result)
	if err != nil {
		return result, ErrorNotFount
	}

	return result, nil
}

//GetAll the contacts
func (uc UserContact) GetAll(ctx context.Context) ([]Contact, error) {
	result := []Contact{}
	cursor, err := uc.coll.Find(ctx, bson.M{"user_id": ctx.Value(primitive.ObjectID{})})
	if err != nil {
		return result, ErrorNotFount
	}
	cursor.All(ctx, &result)

	return result, nil
}

//GetByNumber get by number a contact
func (uc UserContact) GetByNumber(ctx context.Context, number string) (Contact, error) {
	result := Contact{}
	err := uc.coll.FindOne(ctx, bson.M{"number": number, "user_id": ctx.Value(primitive.ObjectID{})}).Decode(&result)
	if err != nil {
		return result, ErrorNotFount
	}
	return result, nil
}

//InsertContact insert one contact
func (uc UserContact) InsertContact(contact Contact) error {
	ctx := context.Background()
	_, err := uc.coll.InsertOne(ctx, contact)
	return err

}

//DeleteByName delete one contact asociate to a user and a name passed
func (uc UserContact) DeleteByName(ctx context.Context, name string) (Contact, error) {
	contact, err := uc.GetByName(ctx, name)
	if err != nil {
		return contact, err
	}

	_, err = uc.coll.DeleteOne(ctx, contact)
	return contact, err

}

//DeleteAll delete all contact asociate to a user
func (uc UserContact) DeleteAll(ctx context.Context) error {
	_, err := uc.coll.DeleteMany(ctx, bson.M{"user_id": ctx.Value(primitive.ObjectID{})})
	return err

}

//UpdateContactByName update a existing contact for a user
func (uc UserContact) UpdateContactByName(ctx context.Context, currendContactName string, updatedContact *Contact) (Contact, error) {
	oldContact, err := uc.GetByName(ctx, currendContactName)

	if err != nil {
		return oldContact, err
	}
	updatedContact.ID = oldContact.ID
	updatedContact.UserID = oldContact.UserID

	_, err = uc.coll.UpdateOne(ctx, oldContact, *updatedContact)

	return *updatedContact, err
}

//NewUserContact return db info
func NewUserContact() UserContact {
	return UserContact{
		data: New(),
		coll: data.DBCollection(ContactCollection),
	}
}
