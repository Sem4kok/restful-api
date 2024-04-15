package db

import (
	"context"
	"github.com/Sem4kok/restful-api/util"
	"github.com/jackc/pgx/v5"
)

func GetAlbumsFromDB(ctx context.Context, conn *pgx.Conn) ([]util.Album, error) {
	// set a Query connection to db
	albums := make([]util.Album, 0)
	rows, err := conn.Query(ctx, "SELECT id, title, artist, price FROM albums")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a util.Album
		if err = rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}

	if err = rows.Err(); err != nil {
		return albums, err
	}

	return albums, nil
}
