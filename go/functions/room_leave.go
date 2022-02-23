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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func roomLeave(ctx context.Context, request *model.RoomLeaveRequest) (*model.RoomLeaveResponse, error) {
	log.Printf("room leave")

	response := &model.RoomLeaveResponse{}

	err := gcp.Firestore.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
		roomRef := db.RoomRef(request.RoomId)
		room, err := db.RoomGetRef(tx, roomRef)

		if err != nil {
			log.Printf("RoomGetRef error")
			return status.Error(codes.InvalidArgument, "invalid room id")
		}

		userRef := db.UserRef(request.UserId)
		_, err = db.UserGetRef(tx, userRef)
		if err != nil {
			log.Printf("UserGetRef error")
			return status.Error(codes.InvalidArgument, "invalid user id")
		}

		err = db.RoomLeave(tx, roomRef, room, userRef)
		if err != nil {
			log.Printf("RoomLeave error")
			return err
		}

		return nil
	}, firestore.MaxAttempts(1))

	return response, err
}

func RoomLeave(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/json" {
		log.Printf("Bad request format")
		writer.WriteHeader(500)
		return
	}

	leaveRequest := model.RoomLeaveRequest{}

	if err := json.NewDecoder(request.Body).Decode(&leaveRequest); err != nil {
		log.Printf("Json Decode RoomLeaveRequest error")
		writer.WriteHeader(500)
		return
	}

	response, err := roomLeave(request.Context(), &leaveRequest)
	if err != nil {
		log.Printf("roomLeave error")
		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Json Encode RoomLeaveResponse error")
		writer.WriteHeader(500)
		return
	}
}
