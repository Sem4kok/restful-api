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
	Albums *[]util.Album
}

// AddData uses goroutine while working with the database
// AddData won't change data in the db
func AddData(cfg *Config, newAlbums []util.Album) {

	// add new albums into db
	// firstly, asynchronous launch of validation requests
	for _, newAlbum := range newAlbums {
		go tryToAddAlbum(&AlbumConfig{
			Cfg:      cfg,
			NewAlbum: newAlbum,
		})
	}

}

type AlbumConfig struct {
	Cfg      *Config
	NewAlbum util.Album
}

// QueryResult structure for writing from the channel
type QueryResult struct {
	ArtistExist bool
	TitleExist  bool
}

func tryToAddAlbum(cfg *AlbumConfig) {
	resultChanel := make(chan QueryResult, 2)
	CheckForExisting(cfg, resultChanel)
	result := <-resultChanel

	log.Printf("result: %v", result)
	// if album already exist case
	if result.ArtistExist && result.TitleExist {
		return
	}

	var newId int
	if err := cfg.Cfg.Conn.QueryRow(
		cfg.Cfg.Ctx,
		"INSERT INTO albums (artist, title, price) VALUES ($1, $2, $3) RETURNING id",
		cfg.NewAlbum.Artist, cfg.NewAlbum.Title, cfg.NewAlbum.Price,
	).Scan(&newId); err != nil {
		log.Printf("error with data insertion: %v", err.Error())
	}

	// adding data into local Albums slice
	cfg.NewAlbum.ID = newId
	*cfg.Cfg.Albums = append(*cfg.Cfg.Albums, cfg.NewAlbum)
}
