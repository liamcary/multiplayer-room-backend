package functions

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"com.liamcary/multiplayer/api/model"
	"com.liamcary/multiplayer/db"
	"com.liamcary/multiplayer/gcp"

	"cloud.google.com/go/firestore"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func roomJoin(ctx context.Context, request *model.RoomJoinRequest) (*model.RoomJoinResponse, error) {
	log.Printf("room create")

	response := &model.RoomJoinResponse{}

	err := gcp.Firestore.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
		timestamp := timestamppb.Now()

		var err error

		roomRef := db.RoomRef(request.Id)
		room, err := db.RoomGetRef(tx, roomRef)

		if err != nil {
			log.Printf("RoomGet error")
			return err
		}

		log.Printf("Finding client user")

		var user *model.User
		var userRef *firestore.DocumentRef

		if request.User.Id == "" {
			log.Printf("Creating new user for host")

			userRef = db.Users.NewDoc()
			user = &model.User{
				Id:           request.User.Id,
				DisplayName:  request.User.DisplayName,
				Location:     request.User.Location,
				Platform:     request.User.Platform,
				TimeCreated:  timestamp,
				LastActivity: timestamp,
			}

			err := db.UserNew(tx, userRef, user)
			if err != nil {
				log.Printf("Failed to create new user")
				return err
			}
		} else {
			log.Printf("Finding existing user for host")

			userRef = db.UserRef(request.User.Id)
			user, err = db.UserGetRef(tx, userRef)
			if err != nil {
				log.Printf("Failed to find existing user")
				return err
			}
		}

		err = db.RoomJoin(tx, roomRef, room, userRef)
		if err != nil {
			log.Printf("RoomJoin error")
			return err
		}

		response.ClientId = userRef.ID
		return nil
	}, firestore.MaxAttempts(1))

	return response, err
}

func RoomJoin(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/json" {
		log.Printf("Bad request format")
		writer.WriteHeader(500)
		return
	}

	joinRequest := model.RoomJoinRequest{}

	if err := json.NewDecoder(request.Body).Decode(&joinRequest); err != nil {
		log.Printf("Json Decode RoomJoinResponse error")
		writer.WriteHeader(500)
		return
	}

	response, err := roomJoin(request.Context(), &joinRequest)
	if err != nil {
		log.Printf("roomJoin error")
		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Json Encode RoomJoinResponse error")
		writer.WriteHeader(500)
		return
	}
}
