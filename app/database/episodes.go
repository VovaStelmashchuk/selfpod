package database

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

var (
	DriverName     = "sqlite"
	DataSourceName = "files/selfpod.db"
)

func init() {
	db, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS episodes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        acast_id TEXT,
        title TEXT,
        audio_url TEXT,
        image_url TEXT,
        description TEXT,
        processing_state INTEGER
    );`
	_, err = db.Exec(createTableSQL)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}

type EpisodeProcessStatus string

const (
	Fail       EpisodeProcessStatus = "Fail"
	InProgress EpisodeProcessStatus = "InProgress"
	Success    EpisodeProcessStatus = "Success"
	NotStarted EpisodeProcessStatus = "NotStarted"
)

type Episode struct {
	Id              int
	AcastId         string
	Title           string
	AudioUrl        string
	ImageUrl        string
	Description     string
	ProcessingState EpisodeProcessStatus
}

func SaveEpisode(episode Episode) (int64, error) {
	client, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	insertSQL := `INSERT INTO episodes (
                      acast_id, 
                      title, 
                      audio_url, 
                      image_url, 
                      processing_state
				  ) VALUES (?, ?, ?, ?, ?)`

	statement, err := client.Prepare(insertSQL)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(
		episode.AcastId,
		episode.Title,
		episode.AudioUrl,
		episode.ImageUrl,
		episode.ProcessingState,
	)

	lastId, err := result.LastInsertId()
	return lastId, err
}

func GetEpisode(id int) (*Episode, error) {
	client, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	var episode Episode
	querySQL := `SELECT id, acast_id, title, audio_url, image_url, processing_state FROM episodes WHERE id = ?`
	row := client.QueryRow(querySQL, id)

	err = row.Scan(
		&episode.Id,
		&episode.AcastId,
		&episode.Title,
		&episode.AudioUrl,
		&episode.ImageUrl,
		&episode.ProcessingState,
	)
	if err != nil {
		return nil, err
	}

	return &episode, nil
}

func GetEpisodeIdsByState(state EpisodeProcessStatus) ([]int, error) {
	client, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	var ids []int
	querySQL := `SELECT id FROM episodes WHERE processing_state = ?`
	rows, err := client.Query(querySQL, state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func UpdateEpisodeState(id int, newState EpisodeProcessStatus) error {
	client, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	updateSQL := `UPDATE episodes SET processing_state = ? WHERE id = ?`
	statement, err := client.Prepare(updateSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(newState, id)
	return err
}

func AddDescriptionToEpisode(id int, description string) error {
	client, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	updateSQL := `UPDATE episodes SET description = ? WHERE id = ?`
	statement, err := client.Prepare(updateSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(description, id)
	return err

}
