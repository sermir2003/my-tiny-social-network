package db_utils

import (
	"database/sql"
	"fmt"
	"log"
	"stats_service_core/utils"

	_ "github.com/ClickHouse/clickhouse-go"
)

var DB *sql.DB

func StartUpDB() error {
	connection_string := fmt.Sprintf(
		"tcp://%s:%s?username=%s&password=%s",
		utils.GetenvSafe("STATS_DB_HOST"),
		utils.GetenvSafe("STATS_DB_PORT"),
		utils.GetenvSafe("STATS_DB_USER"),
		utils.GetenvSafe("STATS_DB_PASSWORD"),
	)
	log.Printf("trying to connect to stats_db at %s\n", connection_string)

	db, err := sql.Open("clickhouse", connection_string)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func AddView(post_id uint64, appraiser_id uint64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM stats_db.views
		WHERE post_id = ? AND appraiser_id = ?`,
		post_id,
		appraiser_id,
	).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = tx.Exec(`
		INSERT INTO stats_db.views (post_id, appraiser_id)
		VALUES (?, ?)`,
		post_id,
		appraiser_id,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func AddLike(post_id uint64, appraiser_id uint64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM stats_db.likes
		WHERE post_id = ? AND appraiser_id = ?`,
		post_id,
		appraiser_id,
	).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = tx.Exec(`
		INSERT INTO stats_db.likes (post_id, appraiser_id)
		VALUES (?, ?)`,
		post_id,
		appraiser_id,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
