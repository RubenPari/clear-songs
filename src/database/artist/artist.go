package artist

import (
	"log"

	"github.com/RubenPari/clear-songs/src/database"
	"github.com/RubenPari/clear-songs/src/models"
)

func Add(artist models.Artist) bool {
	db := database.GetDB()

	exist := db.QueryRow("SELECT * FROM artists WHERE id = ?", artist.Id.String())

	var artistFounded models.Artist

	_ = exist.Scan(&artistFounded.Id, &artistFounded.Name, &artistFounded.Uri)

	if artistFounded.Id != ("") {
		log.Default().Println("Artist already exists")
		log.Default().Println(artistFounded)
		return false
	}

	result, errInsert := db.Exec("INSERT INTO artists (id, name, uri) VALUES (?, ?, ?)", artist.Id, artist.Name, artist.Uri)

	if errInsert != nil {
		log.Default().Println("Error inserting artist")
		log.Default().Println(errInsert)
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected == int64(1)
}

func Get(id string) models.Artist {
	db := database.GetDB()

	exist := db.QueryRow("SELECT * FROM artists WHERE id = ?", id)

	var artist models.Artist

	_ = exist.Scan(&artist.Id, &artist.Uri, &artist.Name)

	return artist
}

func Delete(id string) bool {
	db := database.GetDB()

	result, err := db.Exec("DELETE FROM artists WHERE id = ?", id)

	if err != nil {
		log.Default().Println("Error deleting artist")
		log.Default().Println(err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected == int64(1)
}
