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

func roomCreate(ctx context.Context, request *model.RoomCreateRequest) (*model.RoomCreateResponse, error) {
	log.Printf("room create")

	response := &model.RoomCreateResponse{}

	err := gcp.Firestore.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
		timestamp := timestamppb.Now()

		var err error
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

		log.Printf("Creating Room")

		roomRef := db.Rooms.NewDoc()
		room := &model.Room{
			Name:        request.Name,
			Host:        userRef,
			Location:    user.Location,
			TimeCreated: timestamp,
		}

		err = db.RoomNew(tx, roomRef, room)
		if err != nil {
			log.Printf("RoomNew error")
			return err
		}

		response.HostId = userRef.ID
		response.RoomId = roomRef.ID
		return nil
	}, firestore.MaxAttempts(1))

	return response, err
}

func RoomCreate(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/json" {
		log.Printf("Bad request format")
		writer.WriteHeader(500)
		return
	}

	createRequest := model.RoomCreateRequest{}

	if err := json.NewDecoder(request.Body).Decode(&createRequest); err != nil {
		log.Printf("Json Decode RoomCreateRequest error")
		writer.WriteHeader(500)
		return
	}

	response, err := roomCreate(request.Context(), &createRequest)
	if err != nil {
		log.Printf("roomCreate error")
		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Json Encode RoomCreateResponse error")
		writer.WriteHeader(500)
		return
	}
}
