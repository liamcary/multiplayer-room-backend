package gcp

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

var (
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
)

func init() {
	ctx := context.Background()

	var err error

	App, err = firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Firestore, err = App.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
