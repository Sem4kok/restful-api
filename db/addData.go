package db

import (
	"context"
	"github.com/Sem4kok/restful-api/util"
	"github.com/jackc/pgx/v5"
	"log"
)

type Config struct {
	Ctx    context.Context
	Conn   *pgx.Conn
	Albums []util.Album
}

// AddData uses goroutine while working with the database
// AddData won't change data in the db
func AddData(cfg *Config, newAlbums []util.Album) error {

	// add new albums into db
	// firstly, asynchronous launch of validation requests
	for _, newAlbum := range newAlbums {
		go tryToAddAlbum(&AlbumConfig{
			Cfg:      cfg,
			NewAlbum: newAlbum,
		})
	}

	return nil
}

// ХОЧУ СДЕЛАТЬ МЕТОДЫ
type AlbumConfig struct {
	Cfg      *Config
	NewAlbum util.Album
}

// structure for writing from the channel
type QueryResult struct {
	ArtistExist bool
	TitleExist  bool
}

func tryToAddAlbum(cfg *AlbumConfig) {
	resultChanel := make(chan QueryResult, 1)
	CheckForExisting(cfg, resultChanel)
	result := <-resultChanel

	// if album already exist case
	if result.ArtistExist && result.TitleExist {
		// TODO implement writing message to client about existing
		// TODO we should send UPDATE request
		return
	}

	if _, err := cfg.Cfg.Conn.Exec(
		cfg.Cfg.Ctx,
		"INSERT INTO albums (artist, title, price) VALUES ($1, $2, $3)",
		cfg.NewAlbum.Artist, cfg.NewAlbum.Title, cfg.NewAlbum.Price,
	); err != nil {
		log.Printf("error with data insertion: %v", err.Error())
	}
	// adding data into local Albums slice
	cfg.Cfg.Albums = append(cfg.Cfg.Albums, cfg.NewAlbum)

	// TODO message <- error message
}
