package helper

import (
	"context"
	"server/internal/delivery/ws"
	package_psql "server/pkg/db"
)

func GetAllRooms(ctx context.Context, client package_psql.Client, hub *ws.Hub) error {
	q := `SELECT uuid FROM rooms`

	rows, err := client.Query(ctx, q)
	if err != nil {
		return err
	}

	for rows.Next() {
		var uuid string

		err = rows.Scan(&uuid)
		if err != nil {
			return err
		}

		hub.Rooms[uuid] = &ws.Room{
			ID:      uuid,
			Name:    "test",
			Clients: make(map[string]*ws.Client),
		}
	}

	return nil
}
