package db

import (
	"log"
)

// DeleteAlbumFromDB returns bool -> did album exist
// true if yes. false if didn't.
// returns error if there were problem with deleting element
func DeleteAlbumFromDB(cfg *Config, id int) (bool, error) {
	// delete from database
	result, err := cfg.Conn.Exec(cfg.Ctx,
		"DELETE FROM albums WHERE id = $1", id)
	if err != nil {
		log.Printf("error with deleting album: %v", err.Error())
		return false, err
	}

	if result.RowsAffected() > 0 {
		return true, nil
	}

	return false, nil
}
