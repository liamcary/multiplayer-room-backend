package db

import (
	"com.liamcary/multiplayer/api/model"
	"com.liamcary/multiplayer/gcp"

	"cloud.google.com/go/firestore"
)

var (
	Users = gcp.Firestore.Collection("users")
)

func UserRef(userId string) *firestore.DocumentRef {
	return Users.Doc(userId)
}

func UserGet(tx *firestore.Transaction, userId string) (*model.User, error) {
	return UserGetRef(tx, UserRef(userId))
}

func UserGetRef(tx *firestore.Transaction, userRef *firestore.DocumentRef) (*model.User, error) {
	doc, err := tx.Get(userRef)
	if err != nil {
		return nil, err
	}

	return UserLoad(doc)
}

func UserLoad(doc *firestore.DocumentSnapshot) (*model.User, error) {
	user := &model.User{}

	err := doc.DataTo(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UserNew(tx *firestore.Transaction, userRef *firestore.DocumentRef, user *model.User) error {
	user.Id = userRef.ID

	err := tx.Create(userRef, user)
	if err != nil {
		return err
	}

	return nil
}

func UserDelete(tx *firestore.Transaction, userRef *firestore.DocumentRef) error {
	err := tx.Delete(userRef)
	if err != nil {
		return err
	}

	return nil
}
