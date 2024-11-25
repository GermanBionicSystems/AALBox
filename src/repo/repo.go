// Package repo provides functionalities to interact with a database
// for managing songs associated with tags.package repo
package repo

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Repository struct to manage database interactions
type Repository struct {
	Connection *sql.DB
}

// NewRepository initializes a new Repository with a database connection
func NewRepository() Repository {
	connStr := "user=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to open database connection:", err)
		panic(err)
	}

	// It's a good practice to test the database connection
	if err := db.Ping(); err != nil {
		fmt.Println("Failed to establish database connection:", err)
		panic(err)
	}

	// Create the table only if it does not exist
	tableCreationQuery := `CREATE TABLE IF NOT EXISTS songs (
		tag_id TEXT PRIMARY KEY,
		song_path TEXT NOT NULL
	)`
	if _, err := db.Exec(tableCreationQuery); err != nil {
		fmt.Println("Failed to create table:", err)
	}

	return Repository{
		Connection: db,
	}
}

// GetSongPath retrieves the song path associated with a given word
func (r *Repository) GetSongPath(tagId string) string {
	var songPath string
	err := r.Connection.QueryRow("SELECT song_path FROM songs WHERE tag_id = $1", tagId).Scan(&songPath)
	if err != nil {
		return ""
	}
	return songPath
}

// AddSong adds a new song to the database
func (r *Repository) AddSong(tagId, songPath string) error {
	_, err := r.Connection.Exec("INSERT into songs (tag_id, song_path) VALUES ($1,$2)", tagId, songPath)
	return err
}

// RemoveSong deletes a song form the database
func (r *Repository) RemoveSong(tagId string) error {
	_, err := r.Connection.Exec("DELETE FROM songs WHERE tag_id =$1", tagId)
	return err
}
