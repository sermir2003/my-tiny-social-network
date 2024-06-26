package db_utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	post_pb "ugc_service_core/post"
	"ugc_service_core/utils"

	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var db *sql.DB

func StartUpDB() (err error) {
	connection_string := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		utils.GetenvSafe("POSTGRES_HOST"),
		utils.GetenvSafe("POSTGRES_PORT"),
		utils.GetenvSafe("POSTGRES_USER"),
		utils.GetenvSafe("POSTGRES_PASSWORD"),
		utils.GetenvSafe("POSTGRES_DB"),
		utils.GetenvSafe("TIME_ZONE"),
	)
	log.Printf("trying to connect to ugc_db at %s\n", connection_string)

	db, err = sql.Open("postgres", connection_string)
	if err != nil {
		return err
	}

	return nil
}

func Create(req *post_pb.CreateRequest) (*post_pb.CreateResponse, error) {
	var post_id uint64
	err := db.QueryRow(`
		INSERT INTO posts (author_id, content, create_timestamp)
		VALUES ($1, $2, NOW())
		RETURNING post_id`,
		req.AuthorId,
		req.Content,
	).Scan(&post_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert: %v", err)
	}
	return &post_pb.CreateResponse{
		PostId: post_id,
	}, nil
}

type AccessResult int

const (
	SUCCESS       AccessResult = 0
	ACCESS_DENIED AccessResult = 1
	NOT_FOUND     AccessResult = 2
)

func beginTxAndCheckAccess(post_id uint64, author_id uint64) (AccessResult, *sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, nil, err
	}

	var existing_author_id uint64
	err = tx.QueryRow(`
		SELECT author_id
		FROM posts
		WHERE post_id = $1`,
		post_id,
	).Scan(&existing_author_id)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return NOT_FOUND, nil, nil
	} else if err != nil {
		tx.Rollback()
		return 0, nil, err
	}
	if author_id != existing_author_id {
		tx.Rollback()
		return ACCESS_DENIED, nil, nil
	}

	return SUCCESS, tx, nil
}

func Update(req *post_pb.UpdateRequest) (*post_pb.UpdateResponse, error) {
	access, tx, err := beginTxAndCheckAccess(req.PostId, req.AuthorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start tx or check access: %v", err)
	}
	if access == ACCESS_DENIED {
		return &post_pb.UpdateResponse{
			Access: post_pb.AccessResult_ACCESS_DENIED,
		}, nil
	}
	if access == NOT_FOUND {
		return &post_pb.UpdateResponse{
			Access: post_pb.AccessResult_NOT_FOUND,
		}, nil
	}

	_, err = tx.Exec(`
		UPDATE posts
		SET (content, update_timestamp) = ($1, NOW())
		WHERE post_id = $2`,
		req.Content,
		req.PostId,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &post_pb.UpdateResponse{
		Access: post_pb.AccessResult_SUCCESS,
	}, nil
}

func Delete(req *post_pb.DeleteRequest) (*post_pb.DeleteResponse, error) {
	access, tx, err := beginTxAndCheckAccess(req.PostId, req.AuthorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start tx or check access: %v", err)
	}
	if access == ACCESS_DENIED {
		return &post_pb.DeleteResponse{
			Access: post_pb.AccessResult_ACCESS_DENIED,
		}, nil
	}
	if access == NOT_FOUND {
		return &post_pb.DeleteResponse{
			Access: post_pb.AccessResult_NOT_FOUND,
		}, nil
	}

	_, err = tx.Exec(`
		DELETE FROM posts
		WHERE post_id = $1`,
		req.PostId,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &post_pb.DeleteResponse{
		Access: post_pb.AccessResult_SUCCESS,
	}, nil
}

func scanPost(row interface {
	Scan(dest ...interface{}) error
}) (*post_pb.Post, error) {
	var post post_pb.Post
	var create_timestamp time.Time
	var update_timestamp sql.NullTime
	err := row.Scan(&post.PostId, &post.AuthorId, &post.Content, &create_timestamp, &update_timestamp)
	if err != nil {
		return nil, err
	}
	post.CreateTimestamp = timestamppb.New(create_timestamp)
	if update_timestamp.Valid {
		post.UpdateTimestamp = timestamppb.New(update_timestamp.Time)
	}
	return &post, nil
}

func GetById(req *post_pb.GetByIdRequest) (*post_pb.GetByIdResponse, error) {
	row := db.QueryRow(`
		SELECT post_id, author_id, content, create_timestamp, update_timestamp
		FROM posts
		WHERE post_id = $1`,
		req.PostId,
	)
	post, err := scanPost(row)
	if err == sql.ErrNoRows {
		return &post_pb.GetByIdResponse{
			Access: post_pb.AccessResult_NOT_FOUND,
		}, nil
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select: %v", err)
	}

	return &post_pb.GetByIdResponse{
		Access: post_pb.AccessResult_SUCCESS,
		Post:   post,
	}, nil
}

func GetPagination(req *post_pb.GetPaginationRequest) (*post_pb.GetPaginationResponse, error) {
	rows, err := db.Query(`
		SELECT post_id, author_id, content, create_timestamp, update_timestamp
		FROM posts
		WHERE author_id = $1
		ORDER BY create_timestamp DESC
		LIMIT $2 OFFSET $3`,
		req.AuthorId,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select: %v", err)
	}
	defer rows.Close()

	var posts []*post_pb.Post
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan: %v", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error iterating over rows: %v", err)
	}

	return &post_pb.GetPaginationResponse{
		Posts: posts,
	}, nil
}
