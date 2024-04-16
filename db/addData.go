package db

import (
	"context"
	"github.com/Sem4kok/restful-api/util"
	"github.com/jackc/pgx/v5"
)

type AddConfig struct {
	Ctx       context.Context
	Conn      *pgx.Conn
	Albums    []util.Album
	NewAlbums []util.Album
}

func AddData(c AddConfig) error {

	// add new albums into db
	// firstly checking for distinct

	return nil
}
