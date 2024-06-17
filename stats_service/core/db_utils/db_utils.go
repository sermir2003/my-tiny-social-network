package db_utils

import (
	"database/sql"
	"fmt"
	"log"
	"stats_service_core/utils"

	stats_pb "stats_service_core/stats"

	_ "github.com/ClickHouse/clickhouse-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func AddView(post_id uint64, author_id uint64, appraiser_id uint64) error {
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
		INSERT INTO stats_db.views (post_id, author_id, appraiser_id)
		VALUES (?, ?, ?)`,
		post_id,
		author_id,
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

func AddLike(post_id uint64, author_id uint64, appraiser_id uint64) error {
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
		INSERT INTO stats_db.likes (post_id, author_id, appraiser_id)
		VALUES (?, ?, ?)`,
		post_id,
		author_id,
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

func GetPostStats(post_id uint64) (*stats_pb.GetPostStatsResponse, error) {
	var cnt_views uint64
	err := DB.QueryRow(`
		SELECT COUNT(*)
		FROM stats_db.views
		WHERE post_id = ?`,
		post_id,
	).Scan(&cnt_views)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select views: %v", err)
	}

	var cnt_likes uint64
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM stats_db.likes
		WHERE post_id = ?`,
		post_id,
	).Scan(&cnt_likes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select likes: %v", err)
	}

	return &stats_pb.GetPostStatsResponse{
		CntViews: cnt_views,
		CntLikes: cnt_likes,
	}, nil
}

func GetTopPosts(req *stats_pb.GetTopPostsRequest) (*stats_pb.GetTopPostsResponse, error) {
	var rows *sql.Rows
	var err error
	switch req.Type {
	case stats_pb.ReactionType_VIEW:
		rows, err = DB.Query(`
			SELECT post_id, author_id, COUNT(*) AS count
			FROM stats_db.views
			GROUP BY post_id, author_id
			ORDER BY count DESC
			LIMIT ?`,
			req.TopSize,
		)
	case stats_pb.ReactionType_LIKE:
		rows, err = DB.Query(`
			SELECT post_id, author_id, COUNT(*) AS count
			FROM stats_db.likes
			GROUP BY post_id, author_id
			ORDER BY count DESC
			LIMIT ?`,
			req.TopSize,
		)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select top: %v", err)
	}

	var top []*stats_pb.TopPostItem
	for rows.Next() {
		var item stats_pb.TopPostItem
		err := rows.Scan(&item.PostId, &item.AuthorId, &item.StatsNumber)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan: %v", err)
		}
		top = append(top, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error iterating over rows: %v", err)
	}

	return &stats_pb.GetTopPostsResponse{
		Top: top,
	}, nil
}

func GetTopUsers(req *stats_pb.GetTopUsersRequest) (*stats_pb.GetTopUsersResponse, error) {
	rows, err := DB.Query(`
		SELECT author_id, COUNT(*) AS sum_likes
		FROM stats_db.likes
		GROUP BY author_id
		ORDER BY sum_likes DESC
		LIMIT ?`,
		req.TopSize,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select top: %v", err)
	}

	var top []*stats_pb.TopUserItem
	for rows.Next() {
		var item stats_pb.TopUserItem
		err := rows.Scan(&item.UserId, &item.SumLikes)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan: %v", err)
		}
		top = append(top, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error iterating over rows: %v", err)
	}

	return &stats_pb.GetTopUsersResponse{
		Top: top,
	}, nil
}
