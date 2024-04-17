package db

import "log"

func CheckForExisting(cfg *AlbumConfig, result chan QueryResult) {
	var isArtistExist, isTitleExist bool

	// send a SQL request for check existing
	if err := cfg.Cfg.Conn.QueryRow(
		cfg.Cfg.Ctx,
		"SELECT true FROM albums WHERE artist = $1 LIMIT 1",
		cfg.NewAlbum.Artist,
	).Scan(&isArtistExist); err != nil {
		log.Printf(err.Error())
	}

	// send a SQL request for check existing
	if err := cfg.Cfg.Conn.QueryRow(
		cfg.Cfg.Ctx,
		"SELECT true FROM albums WHERE title = $1 LIMIT 1",
		cfg.NewAlbum.Artist,
	).Scan(&isTitleExist); err != nil {
		log.Printf(err.Error())
	}

	result <- QueryResult{
		ArtistExist: isArtistExist,
		TitleExist:  isTitleExist,
	}
}
