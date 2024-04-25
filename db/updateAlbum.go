package db

import "log"

func UpdateData(cfg *AlbumConfig) {
	if _, err := cfg.Cfg.Conn.Exec(
		cfg.Cfg.Ctx,
		"UPDATE albums SET price = $1, title = $2, artist = $3 WHERE id = $4",
		cfg.NewAlbum.Price, cfg.NewAlbum.Title, cfg.NewAlbum.Artist, cfg.NewAlbum.ID,
	); err != nil {
		log.Printf("error with data update: %v", err.Error())
	}
}
