package repo

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Repository struct {
	Connection *sql.DB
}

func NewRepository() Repository {
	connStr := "user=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return Repository{
		Connection: db,
	}
}

func (r *Repository) GetSongPath(word string) string {
	var songPath string
	err := r.Connection.QueryRow("SELECT song_path FROM songs WHERE word = $1", word).Scan(&songPath)
	if err != nil {
		return ""
	}
	return songPath
}

func (r *Repository) AddSong(word, songPath string) error {
	_, err := r.Connection.Exec("INSERT into songs (word, song_path) VALUES ($1,$2)", word, songPath)
	return err
}

/*
	func ShowSongs(db sql.DB, err error) ([]string, error) {
		var songs []string
		rows, err = db.Query("SELECT word, songPath FROM songs")
		if err != nil {
			return nil, err
		}
		for rows.Next() {

		}
	}
*/
