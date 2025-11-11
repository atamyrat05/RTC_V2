package helper

import (
	"context"
	"server/internal/delivery/ws"
	package_psql "server/pkg/db"
)

func GetAllRooms(ctx context.Context, client package_psql.Client, hub *ws.Hub) error {
	q := `SELECT id FROM rooms`

	rows, err := client.Query(ctx, q)
	if err != nil {
		return err
	}

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			return err
		}

		hub.Rooms[id] = &ws.Room{
			ID:      id,
			Name:    "test",
			Clients: make(map[string]*ws.Client),
		}
	}

	return nil
}
