package db

import (
	"com.liamcary/multiplayer/api/model"
	"com.liamcary/multiplayer/gcp"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Rooms = gcp.Firestore.Collection("rooms")
)

type RoomCreateParams struct {
	Id     string
	HostId string
}

func RoomRef(roomId string) *firestore.DocumentRef {
	return Rooms.Doc(roomId)
}

func RoomGet(tx *firestore.Transaction, roomId string) (*model.Room, error) {
	return RoomGetRef(tx, RoomRef(roomId))
}

func RoomGetRef(tx *firestore.Transaction, roomRef *firestore.DocumentRef) (*model.Room, error) {
	doc, err := tx.Get(roomRef)

	if err != nil {
		return nil, err
	}

	return RoomLoad(doc)
}

func RoomLoad(doc *firestore.DocumentSnapshot) (*model.Room, error) {
	room := &model.Room{}

	err := doc.DataTo(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func RoomNew(tx *firestore.Transaction, roomRef *firestore.DocumentRef, room *model.Room) error {
	room.Id = roomRef.ID

	err := tx.Create(roomRef, room)
	if err != nil {
		return err
	}

	return nil
}

func RoomJoin(tx *firestore.Transaction, roomRef *firestore.DocumentRef, room *model.Room, userRef *firestore.DocumentRef) error {
	err := tx.Update(roomRef, []firestore.Update{
		{
			Path:  "Clients",
			Value: append(room.Clients, userRef),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func RoomLeave(tx *firestore.Transaction, roomRef *firestore.DocumentRef, room *model.Room, userRef *firestore.DocumentRef) error {
	if room.Host != nil && room.Host.ID == userRef.ID {
		room.Host = nil

		return tx.Update(roomRef, []firestore.Update{
			{
				Path:  "Host",
				Value: nil,
			},
		})
	}

	for i := 0; i < len(room.Clients); i++ {
		if room.Clients[i].ID != userRef.ID {
			continue
		}

		copy(room.Clients[i:], room.Clients[i+1:])
		room.Clients = room.Clients[:len(room.Clients)-1]

		return tx.Update(roomRef, []firestore.Update{
			{
				Path:  "Clients",
				Value: room.Clients,
			},
		})
	}

	return status.Error(codes.InvalidArgument, "Couldnt find user in room")
}
