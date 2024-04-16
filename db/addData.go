package db

import (
	"context"
	"github.com/Sem4kok/restful-api/util"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	Ctx    context.Context
	Conn   *pgx.Conn
	Albums []util.Album
}

// AddData uses goroutine while working with the database
func AddData(cfg *Config, newAlbums []util.Album) error {

	// add new albums into db
	// firstly, asynchronous launch of validation requests
	for _, newAlbum := range newAlbums {
		go tryToAddAlbum(&albumConfig{
			Cfg:      cfg,
			NewAlbum: newAlbum,
		})
	}

	return nil
}

type albumConfig struct {
	Cfg      *Config
	NewAlbum util.Album
}

func tryToAddAlbum(c *albumConfig) {
}
