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
	"google.golang.org/api/iterator"
)

func roomGetAll(ctx context.Context) (*model.RoomGetAllResponse, error) {
	log.Printf("room get all")

	response := &model.RoomGetAllResponse{}

	err := gcp.Firestore.RunTransaction(ctx, func(c context.Context, tx *firestore.Transaction) error {
		iter := db.Rooms.Documents(ctx)

		for {
			doc, err := iter.Next()

			if err == iterator.Done {
				break
			}

			if err != nil {
				return err
			}

			var room *model.Room

			err = doc.DataTo(&room)
			if err != nil {
				log.Printf("DataTo error")
				return err
			}

			userCount := len(room.Clients)
			if room.Host != nil {
				userCount += 1
			}

			info := &model.RoomInfo{
				Id:       room.Id,
				Name:     room.Name,
				Users:    int32(userCount),
				Location: room.Location,
			}

			response.Rooms = append(response.Rooms, info)
		}

		return nil
	}, firestore.MaxAttempts(1))

	return response, err
}

func RoomGetAll(writer http.ResponseWriter, request *http.Request) {
	response, err := roomGetAll(request.Context())
	if err != nil {
		log.Printf("roomGetAll error")
		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Json Encode RomGetAllResponse error")
		writer.WriteHeader(500)
		return
	}
}
